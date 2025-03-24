package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-api-project/internal/repository"

	"github.com/gorilla/mux"
)

type PollRequest struct {
	Question string `json:"question"`
}

type VoteRequest struct {
	Vote string `json:"vote"`
}

func CreatePoll(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pollReq PollRequest
		err := json.NewDecoder(r.Body).Decode(&pollReq)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Extract the username from the JWT token
		username, ok := r.Context().Value("username").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		pollID, err := repo.CreatePoll(pollReq.Question, username)
		if err != nil {
			http.Error(w, "Error creating poll", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Poll created successfully",
			"poll_id": pollID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func ListPolls(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		polls, err := repo.ListPolls()
		if err != nil {
			http.Error(w, "Error listing polls", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(polls)
	}
}

func SubmitVote(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pollID, err := strconv.Atoi(vars["pollID"])
		if err != nil {
			http.Error(w, "Invalid poll ID", http.StatusBadRequest)
			return
		}

		var voteReq VoteRequest
		err = json.NewDecoder(r.Body).Decode(&voteReq)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Extract the username from the JWT token
		username, ok := r.Context().Value("username").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user has already voted
		hasVoted, err := repo.HasUserVoted(pollID, username)
		if err != nil {
			http.Error(w, "Error checking vote status", http.StatusInternalServerError)
			return
		}
		if hasVoted {
			http.Error(w, "User has already voted", http.StatusConflict)
			return
		}

		var vote bool
		if voteReq.Vote == "yes" {
			vote = true
		} else if voteReq.Vote == "no" {
			vote = false
		} else {
			http.Error(w, "Invalid vote value", http.StatusBadRequest)
			return
		}

		err = repo.SubmitVote(pollID, username, vote)
		if err != nil {
			http.Error(w, "Error submitting vote", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Vote submitted successfully"))
	}
}
