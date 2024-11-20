package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserGroupPostgres struct {
	db *sqlx.DB
}

func NewUserGroupPostgres(db *sqlx.DB) *UserGroupPostgres {
	return &UserGroupPostgres{db: db}
}

func (r *UserGroupPostgres) Create(userGroup models.UserGroup) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (userid, groupid, roleid, playlistid) values ($1, $2, $3, $4) RETURNING id", userGroupTable)
	row := r.db.QueryRow(query, userGroup.UserID, userGroup.GroupID, userGroup.RoleID, userGroup.PlayListID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserGroupPostgres) GetAll() ([]models.AllUserGroupData, error) {
	var groups []models.AllUserGroupData
	query := fmt.Sprintf("SELECT ug.id, u.id as userid, u.id, u.email, r.name as rolename, ug.playlistid "+
		"FROM %s ug JOIN %s u ON u.id=ug.userid JOIN %s r ON ug.roleid=r.id ORDER BY ug.id ASC ",
		userGroupTable, usersTable, rolesTable)
	err := r.db.Select(&groups, query)
	return groups, err
}

func (r *UserGroupPostgres) GetById(id int) (*models.AllUserGroupData, error) {
	var userGroup models.AllUserGroupData
	query := fmt.Sprintf("SELECT ug.id, u.id as userid, u.email, r.name as rolename, ug.playlistid "+
		"FROM %s ug JOIN %s u ON u.id=ug.userid JOIN %s r ON ug.roleid=r.id WHERE ug.groupID=$1",
		userGroupTable, usersTable, rolesTable)
	err := r.db.Get(&userGroup, query, id)
	if userGroup == (models.AllUserGroupData{}) {
		return nil, nil
	}
	return &userGroup, err
}

func (r *UserGroupPostgres) GetUsersByGroupId(id int) ([]models.UserGroupData, error) {
	var userGroupData []models.UserGroupData
	query := fmt.Sprintf("SELECT u.id, u.username, u.email, u.platform, r.name as rolename, p.id as playlistid "+
		"FROM %s ug JOIN %s u ON u.id=ug.userid JOIN %s r ON ug.roleid=r.id JOIN %s p ON p.groupid=ug.groupid "+
		"WHERE ug.groupID=$1 ORDER BY u.id ASC ", userGroupTable, usersTable, rolesTable, playlistsTable)
	err := r.db.Select(&userGroupData, query, id)
	return userGroupData, err
}

func (r *UserGroupPostgres) GetByGroupIdAndUserIDAllData(groupID, userID int) (*models.UserGroupData, error) {
	var userGroupData models.UserGroupData
	query := fmt.Sprintf("SELECT u.id, u.username, u.email, u.platform, r.name as rolename, p.id as playlistid "+
		"FROM %s ug JOIN %s u ON u.id=ug.userid JOIN %s r ON ug.roleid=r.id JOIN %s p ON p.groupid=ug.groupid "+
		"WHERE ug.groupID=$1 AND ug.userID=$2", userGroupTable, usersTable, rolesTable, playlistsTable)
	err := r.db.Get(&userGroupData, query, groupID, userID)
	if userGroupData == (models.UserGroupData{}) {
		return nil, nil
	}
	return &userGroupData, err
}

func (r *UserGroupPostgres) GetByGroupIdAndUserIDRole(groupID, userID int, role string) (*models.UserGroup, error) {
	var userGroup models.UserGroup
	query := fmt.Sprintf("SELECT ug.id, ug.userid, ug.groupid, ug.roleid, ug.playlistid "+
		"FROM %s ug JOIN %s r ON r.id=ug.roleid AND r.name=$1 WHERE ug.groupID=$2 AND ug.userid=$3", userGroupTable, rolesTable)
	err := r.db.Get(&userGroup, query, role, groupID, userID)
	if userGroup == (models.UserGroup{}) {
		return nil, nil
	}
	return &userGroup, err
}

func (r *UserGroupPostgres) GetByGroupIdAndUserID(groupID, userID int) (*models.UserGroup, error) {
	var userGroup models.UserGroup
	query := fmt.Sprintf("SELECT * FROM %s WHERE groupID=$1 AND userid=$2", userGroupTable)
	err := r.db.Get(&userGroup, query, groupID, userID)
	if userGroup == (models.UserGroup{}) {
		return nil, nil
	}
	return &userGroup, err
}

func (r *UserGroupPostgres) GetByGroupIdSpotifyUser(id int) ([]models.UserGroupToken, error) {
	var users []models.UserGroupToken
	query := fmt.Sprintf("SELECT ug.id, ug.userid, ug.groupid, ug.roleid, us.accessToken, us.spotifyUri, "+
		"ug.playlistid FROM %s ug JOIN %s u ON u.id = ug.userid AND  u.platform = 'Spotify' "+
		"JOIN %s us ON u.id = us.userid WHERE ug.groupID=$1 ORDER BY ug.id ASC ", userGroupTable, usersTable, userSpotifyTable)
	err := r.db.Select(&users, query, id)
	return users, err
}

func (r *UserGroupPostgres) GetByGroupIdYouTubeMusicUser(id int) ([]models.UserGroupToken, error) {
	var users []models.UserGroupToken
	query := fmt.Sprintf("SELECT ug.id, ug.userid, ug.groupid, ug.roleid, ut.accessToken, "+
		"ug.playlistid FROM %s ug JOIN %s u ON u.id = ug.userid AND  u.platform='YouTubeMusic' "+
		"JOIN %s ut ON u.id = ut.userid WHERE ug.groupID=$1 ORDER BY ug.id ASC", userGroupTable, usersTable, userYouTubeMusicTable)
	err := r.db.Select(&users, query, id)
	return users, err
}

func (r *UserGroupPostgres) Update(id, updateUserID, roleID int) error {
	query := fmt.Sprintf("UPDATE %s SET roleid=$1 WHERE groupID=$2 AND userID=$3", userGroupTable)
	_, err := r.db.Exec(query, roleID, id, updateUserID)
	return err
}

func (r *UserGroupPostgres) Delete(userID, groupID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE userID=$1 AND groupID=$2", userGroupTable)
	_, err := r.db.Exec(query, userID, groupID)
	return err
}
