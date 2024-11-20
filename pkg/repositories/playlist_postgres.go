package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type PlayListPostgres struct {
	db *sqlx.DB
}

func NewPlayListPostgres(db *sqlx.DB) *PlayListPostgres {
	return &PlayListPostgres{db: db}
}

func (r *PlayListPostgres) Create(playlist models.PlayList) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, description, groupID) values ($1, $2, $3) RETURNING id", playlistsTable)
	row := r.db.QueryRow(query, playlist.Name, playlist.Description, playlist.GroupID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PlayListPostgres) GetAll() ([]models.PlayList, error) {
	var playLists []models.PlayList
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC ", playlistsTable)
	err := r.db.Select(&playLists, query)
	return playLists, err
}

func (r *PlayListPostgres) GetById(id int) (*models.PlayList, error) {
	var playList models.PlayList
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", playlistsTable)
	err := r.db.Get(&playList, query, id)
	if playList == (models.PlayList{}) {
		return nil, nil
	}
	return &playList, err
}

func (r *PlayListPostgres) GetByGroupId(id int) (*models.PlayList, error) {
	var playList models.PlayList
	query := fmt.Sprintf("SELECT * FROM %s WHERE groupID=$1", playlistsTable)
	err := r.db.Get(&playList, query, id)
	if playList == (models.PlayList{}) {
		return nil, nil
	}
	return &playList, err
}

func (r *PlayListPostgres) Update(id int, playList models.UpdatePlayListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if playList.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *playList.Name)
		argId++
	}

	if playList.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *playList.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", playlistsTable, setQuery, argId)
	args = append(args, id)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *PlayListPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", playlistsTable)
	_, err := r.db.Exec(query, id)
	return err
}
