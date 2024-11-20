package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type GroupPostgres struct {
	db *sqlx.DB
}

func NewGroupPostgres(db *sqlx.DB) *GroupPostgres {
	return &GroupPostgres{db: db}
}

func (r *GroupPostgres) Create(userID, roleID int, group models.UserCreateGroup) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var ID int
	createGroup := fmt.Sprintf("INSERT INTO %s (name, description) values ($1, $2) RETURNING ID", groupsTable)
	row1 := r.db.QueryRow(createGroup, group.GroupName, group.GroupDescription)
	if err := row1.Scan(&ID); err != nil {
		tx.Rollback()
		return 0, err
	}
	var userGroupID int
	userGroup := fmt.Sprintf("INSERT INTO %s (userId, groupID, roleID, playListID) values ($1, $2, $3, $4) RETURNING ID", userGroupTable)
	row2 := r.db.QueryRow(userGroup, userID, ID, roleID, group.PlayListID)
	if err := row2.Scan(&userGroupID); err != nil {
		tx.Rollback()
		return 0, err
	}
	var playListID int
	playList := fmt.Sprintf("INSERT INTO %s (name, description, groupID) values ($1, $2, $3) RETURNING ID", playlistsTable)
	row3 := r.db.QueryRow(playList, group.PlayListName, group.PlayListDescription, ID)
	if err := row3.Scan(&playListID); err != nil {
		tx.Rollback()
		return 0, err
	}
	return ID, tx.Commit()
}

func (r *GroupPostgres) GetAll() ([]models.Group, error) {
	var groups []models.Group
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC ", groupsTable)
	err := r.db.Select(&groups, query)
	return groups, err
}

func (r *GroupPostgres) GetById(id int) (*models.Group, error) {
	var group models.Group
	query := fmt.Sprintf("SELECT * FROM %s WHERE id= $1", groupsTable)
	err := r.db.Get(&group, query, id)
	if group == (models.Group{}) {
		return nil, nil
	}
	return &group, err
}

func (r *GroupPostgres) Update(id int, group models.UpdateGroupInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if group.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *group.Name)
		argId++
	}

	if group.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *group.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", groupsTable, setQuery, argId)
	args = append(args, id)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *GroupPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", groupsTable)
	_, err := r.db.Exec(query, id)
	return err
}
