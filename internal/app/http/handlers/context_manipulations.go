package handlers

import (
	"context"
	"github.com/radiologist-ai/web-app/internal/domain"
)

func SetCurrentUser(ctx context.Context, user domain.UserRepoModel) error {
	return nil
}

func GetCurrentUser(ctx context.Context) (*domain.UserRepoModel, bool) {
	currentUser, ok := ctx.Value(domain.CurrentUserCtxKey).(domain.UserRepoModel)
	if !ok {
		return nil, false
	}
	return &currentUser, true
}
