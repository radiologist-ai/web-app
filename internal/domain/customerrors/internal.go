package customerrors

import (
	"errors"
	"fmt"
)

var (
	InternalError    = errors.New("internal error. ")
	InternalErrorSQL = fmt.Errorf("%wfailed to execute sql query. ", InternalError)
)
