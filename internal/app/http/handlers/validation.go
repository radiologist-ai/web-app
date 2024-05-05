package handlers

import (
	"errors"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
)

func ValidationErrorToResponseText(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, customerrors.ValidationErrorPasswordTooShort):
		return "Password too short, use at least 8 characters."
	case errors.Is(err, customerrors.ValidationErrorPasswordTooLong):
		return "Password too long, use at most 64 characters."
	case errors.Is(err, customerrors.ValidationErrorPasswordNoLetters):
		return "Password must contain at least one letter."
	case errors.Is(err, customerrors.ValidationErrorPasswordNoDigits):
		return "Password must contain at least one digits."
	case errors.Is(err, customerrors.ValidationErrorPasswordUnacceptableCharacters):
		return "Password can only contain letters and digits."
	default:
		return err.Error()
	}
}
