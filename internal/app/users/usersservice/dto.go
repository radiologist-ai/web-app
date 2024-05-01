package usersservice

import (
	"crypto/sha256"
	"github.com/radiologist-ai/web-app/internal/domain"
)

func (s *Service) userFormToUserRepoModel(in domain.UserForm) domain.UserRepoModel {
	var out domain.UserRepoModel
	out.FirstName = in.FirstName
	out.LastName = in.LastName
	out.Email = in.Email
	out.IsDoctor = in.IsDoctor
	hasher := sha256.New()
	hasher.Write([]byte(in.Password))
	out.PasswordHash = hasher.Sum(nil)
	return out
}
