package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api-project/config"
	"go-api-project/internal/models"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PostgresRepository struct {
	*sql.DB
}

func ConnectPostgresDB() (Repository, error) {
	cfg := config.LoadConfig()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.DB.Close()
}

func (repo *PostgresRepository) CreatePoll(question, createdBy string) (int, error) {
	var pollID int
	err := repo.QueryRow("INSERT INTO polls (question, created_by) VALUES ($1, $2) RETURNING id", question, createdBy).Scan(&pollID)
	return pollID, err
}

func (repo *PostgresRepository) ListPolls() ([]models.Poll, error) {
	rows, err := repo.Query("SELECT id, question, created_by, created_at, yes_count, no_count FROM polls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var polls []models.Poll
	for rows.Next() {
		var poll models.Poll
		if err := rows.Scan(&poll.ID, &poll.Question, &poll.CreatedBy, &poll.CreatedAt, &poll.YesCount, &poll.NoCount); err != nil {
			return nil, err
		}
		polls = append(polls, poll)
	}
	return polls, nil
}

func (repo *PostgresRepository) HasUserVoted(pollID int, username string) (bool, error) {
	var exists bool
	err := repo.QueryRow("SELECT EXISTS(SELECT 1 FROM votes WHERE poll_id=$1 AND username=$2)", pollID, username).Scan(&exists)
	return exists, err
}

func (repo *PostgresRepository) SubmitVote(pollID int, username string, vote bool) error {
	_, err := repo.Exec("INSERT INTO votes (poll_id, username, vote) VALUES ($1, $2, $3)", pollID, username, vote)
	if err != nil {
		return err
	}

	if vote {
		_, err = repo.Exec("UPDATE polls SET yes_count = yes_count + 1 WHERE id = $1", pollID)
	} else {
		_, err = repo.Exec("UPDATE polls SET no_count = no_count + 1 WHERE id = $1", pollID)
	}
	return err
}

func (repo *PostgresRepository) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = repo.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	return err
}

func (repo *PostgresRepository) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	err := repo.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

func (repo *PostgresRepository) UsernameExists(username string) (bool, error) {
	var exists bool
	err := repo.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	return exists, err
}
