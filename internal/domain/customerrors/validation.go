package customerrors

import (
	"errors"
	"fmt"
)

var (
	ValidationError                               = errors.New("validation error. ")
	ValidationErrorPassword                       = fmt.Errorf("%winvalid password. ", ValidationError)
	ValidationErrorEmail                          = fmt.Errorf("%winvalid email. ", ValidationError)
	ValidationErrorPasswordTooShort               = fmt.Errorf("%wpassword must contain at least 8 characters. ", ValidationErrorPassword)
	ValidationErrorPasswordTooLong                = fmt.Errorf("%wpassword can't contain more than 64 characters. ", ValidationErrorPassword)
	ValidationErrorPasswordUnacceptableCharacters = fmt.Errorf("%wpassword can contain only latin letters and digits. ", ValidationErrorPassword)
	ValidationErrorPasswordNoLetters              = fmt.Errorf("%wpassword must containat least 1 letter. ", ValidationErrorPassword)
	ValidationErrorPasswordNoDigits               = fmt.Errorf("%wpassword must contain at least 1 digit. ", ValidationErrorPassword)
	ValidationErrorFirstName                      = fmt.Errorf("%winvalid first name. ", ValidationError)
	ValidationErrorLastName                       = fmt.Errorf("%winvalid last name. ", ValidationError)
	ValidationErrorFirstNameEmpty                 = fmt.Errorf("%wempty first name. ", ValidationErrorFirstName)
	ValidationErrorLastNameEmpty                  = fmt.Errorf("%wempty last name. ", ValidationErrorLastName)
)
