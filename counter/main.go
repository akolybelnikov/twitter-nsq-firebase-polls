package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"flag"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

const updateDuration = 1 * time.Second

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

func main() {

	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	log.Println("Connecting to database...")
	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CONFIG"))
	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	defer client.Close()

	pollData := client.Collection("polls")

	var counts map[string]int
	var countsLock sync.Mutex

	log.Println("Connecting to nsq...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}

	log.Println("Waiting for votes on nsq...")
	ticker := time.NewTicker(updateDuration)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-ticker.C:
			doCount(&countsLock, &counts, pollData, client)
		case <-termChan:
			ticker.Stop()
			q.Stop()
		case <-q.StopChan:
			// finished
			return
		}
	}

}

func doCount(countsLock *sync.Mutex, counts *map[string]int, pollData *firestore.CollectionRef, client *firestore.Client) {
	countsLock.Lock()
	defer countsLock.Unlock()
	if len(*counts) == 0 {
		log.Println("No new votes, skipping database update")
		return
	}
	log.Println("Updating database...")
	log.Println(*counts)
	ok := true
	ctx := context.Background()
	for option, count := range *counts {
		sel := pollData.Where("options", "in", []string{option}).Documents(ctx)
		batch := client.Batch()
		for {
			doc, err := sel.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			batch.Set(doc.Ref, map[string]interface{}{
				"results." + option: count,
			}, firestore.MergeAll)
		}
		_, err := batch.Commit(ctx)
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Println("failed to update:", err)
			ok = false
		}
	}
	if ok {
		log.Println("Finished updating database...")
		*counts = nil // reset counts
	}
}
