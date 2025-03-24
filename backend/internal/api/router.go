package api

import (
	"go-api-project/internal/repository"

	"github.com/gorilla/mux"
)

func NewRouter(repo repository.Repository) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", Login(repo)).Methods("POST")
	router.HandleFunc("/signup", Signup(repo)).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(JWTAuthMiddleware)
	protected.HandleFunc("/polls", CreatePoll(repo)).Methods("POST")
	protected.HandleFunc("/polls", ListPolls(repo)).Methods("GET")
	protected.HandleFunc("/polls/{pollID}/vote", SubmitVote(repo)).Methods("POST")

	return router
}
