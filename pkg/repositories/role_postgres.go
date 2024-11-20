package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type RolePostgres struct {
	db *sqlx.DB
}

func NewRolePostgres(db *sqlx.DB) *RolePostgres {
	return &RolePostgres{db: db}
}

func (r *RolePostgres) Create(role models.Role) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) values ($1) RETURNING id", rolesTable)
	row := r.db.QueryRow(query, role.Name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RolePostgres) GetAll() ([]models.Role, error) {
	var roles []models.Role
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC ", rolesTable)
	err := r.db.Select(&roles, query)
	return roles, err
}

func (r *RolePostgres) GetById(id int) (*models.Role, error) {
	var role models.Role
	query := fmt.Sprintf("SELECT * FROM %s WHERE id= $1", rolesTable)
	err := r.db.Get(&role, query, id)
	if role == (models.Role{}) {
		return nil, nil
	}
	return &role, err
}

func (r *RolePostgres) GetByName(name string) (*models.Role, error) {
	var role models.Role
	query := fmt.Sprintf("SELECT * FROM %s WHERE name= $1", rolesTable)
	err := r.db.Get(&role, query, name)
	if role == (models.Role{}) {
		return nil, nil
	}
	return &role, err
}

func (r *RolePostgres) Update(id int, role models.UpdateRoleInput) error {
	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2", rolesTable)
	_, err := r.db.Exec(query, role.Name, id)
	return err
}

func (r *RolePostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", rolesTable)
	_, err := r.db.Exec(query, id)
	return err
}
