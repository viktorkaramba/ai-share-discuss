package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, email, platform) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Username, user.Email, user.Platform)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CreateUserSpotify(spotifyUri string, token models.ApiToken, user models.User) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var ID int
	createUser := fmt.Sprintf("INSERT INTO %s (username, email, platform) values ($1, $2, $3) RETURNING id", usersTable)
	row1 := r.db.QueryRow(createUser, user.Username, user.Email, user.Platform)
	if err := row1.Scan(&ID); err != nil {
		tx.Rollback()
		return 0, err
	}
	var userSpotifyID int
	userGroup := fmt.Sprintf("INSERT INTO %s (userId, spotifyUri, accessToken, refreshToken) values ($1, $2, $3, $4) RETURNING ID", userSpotifyTable)
	row2 := r.db.QueryRow(userGroup, ID, spotifyUri, token.AccessToken, token.RefreshToken)
	if err := row2.Scan(&userSpotifyID); err != nil {
		tx.Rollback()
		return 0, err
	}
	return ID, nil
}

func (r *AuthPostgres) CreateUserYouTubeMusic(token string, user models.User) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var ID int
	createUser := fmt.Sprintf("INSERT INTO %s (username, email, platform) values ($1, $2, $3) RETURNING id", usersTable)
	row1 := r.db.QueryRow(createUser, user.Username, user.Email, user.Platform)
	if err := row1.Scan(&ID); err != nil {
		tx.Rollback()
		return 0, err
	}
	var userYouTubeID int
	userGroup := fmt.Sprintf("INSERT INTO %s (userId, accessToken) values ($1, $2) RETURNING ID", userYouTubeMusicTable)
	row2 := r.db.QueryRow(userGroup, ID, token)
	if err := row2.Scan(&userYouTubeID); err != nil {
		tx.Rollback()
		return 0, err
	}
	return ID, nil
}

func (r *AuthPostgres) GetUser(email string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1", usersTable)
	err := r.db.Get(&user, query, email)
	if user == (models.User{}) {
		return nil, nil
	}
	return &user, err
}

func (r *AuthPostgres) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, userID)
	if user == (models.User{}) {
		return nil, nil
	}
	return &user, err
}

func (r *AuthPostgres) GetUserSpotifyByID(userID int) (*models.UserSpotify, error) {
	var user models.UserSpotify
	query := fmt.Sprintf("SELECT * FROM %s WHERE userID=$1", userSpotifyTable)
	err := r.db.Get(&user, query, userID)
	if user == (models.UserSpotify{}) {
		return nil, nil
	}
	return &user, err
}
