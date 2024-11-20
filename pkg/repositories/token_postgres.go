package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type TokenPostgres struct {
	db *sqlx.DB
}

func NewTokenPostgres(db *sqlx.DB) *TokenPostgres {
	return &TokenPostgres{db: db}
}

func (t *TokenPostgres) GetSpotifyToken(spotifyUri string) (models.ApiToken, error) {
	var token models.ApiToken
	query := fmt.Sprintf("SELECT us.accessToken, us.refreshToken FROM %s us WHERE us.spotifyUri=$1", userSpotifyTable)
	err := t.db.Get(&token, query, spotifyUri)
	return token, err
}

func (t *TokenPostgres) GetYouTubeMusicToken(userID int) (models.ApiToken, error) {
	var token models.ApiToken
	query := fmt.Sprintf("SELECT ut.accesstoken FROM %s ut WHERE ut.userID=$1", userYouTubeMusicTable)
	err := t.db.Get(&token, query, userID)
	return token, err
}

func (t *TokenPostgres) Create(token models.Token) (int, error) {
	var ID int
	query := fmt.Sprintf("INSERT INTO %s (userID, tokenValue, revoked) values ($1, $2, $3) RETURNING id", tokensTable)
	row := t.db.QueryRow(query, token.UserID, token.TokenValue, strconv.FormatBool(token.Revoked))
	if err := row.Scan(&ID); err != nil {
		return 0, err
	}
	return ID, nil
}

func (t *TokenPostgres) GetByToken(token string) (*models.Token, error) {
	var newToken models.Token
	query := fmt.Sprintf("SELECT * FROM %s WHERE tokenValue=$1 ORDER BY id ASC ", tokensTable)
	err := t.db.Get(&newToken, query, token)
	if newToken == (models.Token{}) {
		return &models.Token{}, nil
	}
	return &newToken, err
}

func (t *TokenPostgres) Update(token string) error {
	query := fmt.Sprintf("UPDATE %s SET revoked='true' WHERE tokenValue=$1", tokensTable)
	_, err := t.db.Exec(query, token)
	return err
}

func (t *TokenPostgres) UpdateYouTubeAccessToken(token string, userID int) error {
	query := fmt.Sprintf("UPDATE %s SET accessToken=$1 WHERE userID=$2", userYouTubeMusicTable)
	_, err := t.db.Exec(query, token, userID)
	return err
}

func (t *TokenPostgres) RevokeAllUserTokens(userID int) error {
	query := fmt.Sprintf("UPDATE %s SET revoked='true' WHERE userID=$1", tokensTable)
	_, err := t.db.Exec(query, userID)
	return err
}

func (t *TokenPostgres) UpdateSpotifyTokenBySpotifyUri(accessToken, spotifyUri string) error {
	query := fmt.Sprintf("UPDATE %s SET accessToken=$1 WHERE spotifyUri=$2", userSpotifyTable)
	_, err := t.db.Exec(query, accessToken, spotifyUri)
	return err
}
