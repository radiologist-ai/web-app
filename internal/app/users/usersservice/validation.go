package usersservice

import (
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"net/mail"
	"unicode"
)

func (s *Service) ValidatePassword(pwd string) error {
	runes := []rune(pwd)
	if len(runes) < 8 {
		return customerrors.ValidationErrorPasswordTooShort
	}
	if len(runes) > 64 {
		return customerrors.ValidationErrorPasswordTooLong
	}
	var (
		containsLetter bool
		containsDigit  bool
	)
	for _, r := range runes {
		if unicode.IsDigit(r) {
			containsDigit = true
		} else if unicode.IsLetter(r) {
			containsLetter = true
		} else {
			return customerrors.ValidationErrorPasswordUnacceptableCharacters
		}
	}
	if !containsLetter {
		return customerrors.ValidationErrorPasswordNoLetters
	}
	if !containsDigit {
		return customerrors.ValidationErrorPasswordNoDigits
	}
	return nil
}

// TODO check if email already in use
func (s *Service) validateRegisterForm(form domain.UserForm) error {
	if form.LastName == "" {
		return customerrors.ValidationErrorLastNameEmpty
	}
	if form.FirstName == "" {
		return customerrors.ValidationErrorFirstNameEmpty
	}
	if _, err := mail.ParseAddress(form.Email); err != nil {
		return customerrors.ValidationErrorEmail
	}
	if err := s.ValidatePassword(form.Password); err != nil {
		return err
	}

	return nil
}
