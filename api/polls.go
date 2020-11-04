package main

import (
	"context"
	"net/http"

	"google.golang.org/api/iterator"
)

type poll struct {
	// unique ID of the poll
	ID string `json:"id"`
	Title string `json:"title"`
	Options []string `json:"options"`
	Results map[string]int `json:"results,omitempty"`
	APIKey string `json:"apikey"`
}

func (s *Server) handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handlePollsGet(w, r)
		return
	case "POST":
		s.handlePollsPost(w, r)
		return
	case "DELETE":
		s.handlePollsDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
}

func (s *Server)handlePollsGet(w http.ResponseWriter, r *http.Request)  {
	ctx := context.Background()
	p := NewPath(r.URL.Path)
	var result []*poll
	if p.HasID() {
		dsnap, err := s.db.Collection("polls").Doc(p.ID).Get(ctx)
		if err != nil {
			respondErr(w, r, http.StatusNotFound, err)
			return
		}
		var entry poll
		dsnap.DataTo(&entry)
		result = append(result, &entry)
	} else {
		iter := s.db.Collection("polls").Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				respondErr(w, r, http.StatusInternalServerError, "failed to read polls", err)
				return
			}
			var entry poll
			doc.DataTo(&entry)
			result = append(result, &entry)
		}
	}
	respond(w, r, http.StatusOK, &result)
}

func (s *Server)handlePollsPost(w http.ResponseWriter, r *http.Request)  {
	ctx := context.Background()
	ref := s.db.Collection("polls").NewDoc().ID
	var p poll
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read poll from request", err)
		return
	}
	apikey, ok := APIKey(r.Context())
	if ok {
		p.APIKey = apikey
	}
	p.ID = ref
	_, err := s.db.Collection("polls").Doc(ref).Set(ctx, p)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert poll", err)
		return
	}
	w.Header().Set("Location", "polls/"+p.ID)
	respond(w, r, http.StatusCreated, p.ID)
}

func (s *Server)handlePollsDelete(w http.ResponseWriter, r *http.Request)  {
	ctx := context.Background()
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all polls")
		return
	}
	_, err := s.db.Collection("polls").Doc(p.ID).Delete(ctx)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete poll", err)
		return
	}
	respond(w, r, http.StatusOK, nil)
}