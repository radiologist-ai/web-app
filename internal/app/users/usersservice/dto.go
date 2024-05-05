package usersservice

import (
	"github.com/radiologist-ai/web-app/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) userFormToUserRepoModel(in domain.UserForm) (domain.UserRepoModel, error) {
	var out domain.UserRepoModel
	out.FirstName = in.FirstName
	out.LastName = in.LastName
	out.Email = in.Email
	out.IsDoctor = in.IsDoctor
	var err error
	out.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		return domain.UserRepoModel{}, err
	}
	return out, nil
}
