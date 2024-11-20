package services

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
)

type RoleService struct {
	repo repositories.Role
}

func NewRoleService(repo repositories.Role) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (s *RoleService) Create(role models.Role) (int, error) {
	return s.repo.Create(role)
}

func (s *RoleService) GetAll() ([]models.Role, error) {
	return s.repo.GetAll()
}

func (s *RoleService) GetById(id int) (*models.Role, error) {
	return s.repo.GetById(id)
}

func (s *RoleService) Update(id int, role models.UpdateRoleInput) error {
	if err := role.Validate(); err != nil {
		return err
	}
	return s.repo.Update(id, role)
}

func (s *RoleService) Delete(id int) error {
	return s.repo.Delete(id)
}
