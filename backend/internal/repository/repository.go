package repository

import (
	"go-api-project/internal/models"
)

type Repository interface {
	Close() error
	CreatePoll(question, createdBy string) (int, error)
	ListPolls() ([]models.Poll, error)
	HasUserVoted(pollID int, username string) (bool, error)
	SubmitVote(pollID int, username string, vote bool) error
	CreateUser(username, password string) error
	AuthenticateUser(username, password string) (*models.User, error)
	UsernameExists(username string) (bool, error)
}
