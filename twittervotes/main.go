package main

import (
	"bufio"
	"context"
	"encoding/json"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/nsqio/go-nsq"
	"google.golang.org/api/option"
)

var (
	conn       net.Conn
	authClient *oauth.Client
	creds      *oauth.Credentials
)

// Poll data struct
type Poll struct {
	options []string
}

// Tweet data structure
type Tweet struct {
	Text string
}

var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

func main() {
	// https://stackoverflow.com/questions/35419263/using-a-configuration-file-with-a-compiled-go-program?noredirect=1&lq=1
	// configFile := os.Getenv("FIREBASE_CONFIG")
	// data, err := ioutil.ReadFile(configFile)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// var config string
	// err = json.Unmarshal(data, &config)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
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

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if conn != nil {
					conn.Close()
					conn = nil
				}
				netc, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				conn = netc
				return netc, nil
			},
		},
	}
	creds = &oauth.Credentials{
		Token:  os.Getenv("SP_TWITTER_ACCESSTOKEN"),
		Secret: os.Getenv("SP_TWITTER_ACCESSSECRET"),
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  os.Getenv("SP_TWITTER_KEY"),
			Secret: os.Getenv("SP_TWITTER_SECRET"),
		},
	}

	twitterStopChan := make(chan struct{}, 1)
	publisherStopChan := make(chan struct{}, 1)
	stop := false
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stop = true
		log.Println("stopping...")
		closeConn()
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	votes := make(chan string)

	go func() {
		pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher: Stopping...")
		pub.Stop()
		log.Println("Publisher: Stopped")
		publisherStopChan <- struct{}{}
	}()

	go func() {
		defer func() {
			twitterStopChan <- struct{}{}
		}()
		for {
			if stop {
				log.Println("Twitter: Stopped")
				return
			}
			time.Sleep(2 * time.Second) // calm

			var options []string
			iter := client.Collection("polls").Documents(ctx)
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalf("Failed to iterate: %v", err)
				}
				var poll Poll
				doc.DataTo(&poll)
				options = append(options, poll.options ...)
			}

			hashtags := make([]string, len(options))
			for i := range options {
				hashtags[i] = "#" + strings.ToLower(options[i])
			}

			form := url.Values{"track": {strings.Join(hashtags, ",")}}

			formEnc := form.Encode()

			u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
			if err != nil {
				log.Println("creating filter request failed:", err)
				return
			}

			req, err := http.NewRequest("POST", u.String(), strings.NewReader(formEnc))
			if err != nil {
				log.Println("creating filter request failed:", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
			req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", u, form))

			resp, err := httpClient.Do(req)
			if err != nil {
				log.Println("Error getting response:", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				// this is a nice way to see what the error actually is:
				s := bufio.NewScanner(resp.Body)
				s.Scan()
				log.Println(s.Text())
				log.Println(hashtags)
				log.Println("StatusCode =", resp.StatusCode)
				continue
			}

			reader = resp.Body
			decoder := json.NewDecoder(reader)
			for {
				var t Tweet
				if err := decoder.Decode(&t); err == nil {
					for _, option := range options {
						if strings.Contains(
							strings.ToLower(t.Text),
							strings.ToLower(option),
						) {
							log.Println("vote:", option)
							votes <- option
						}
					}
				} else {
					break
				}
			}
		}
	}()

	// update by forcing the connection to close
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			if stop {
				break
			}
		}
	}()

	<-twitterStopChan
	close(votes)
	<-publisherStopChan

}
