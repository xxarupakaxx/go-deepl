package deepl

import (
	"errors"
	"fmt"
)

var (
	ErrNilText       = errors.New("text should not nil")
	ErrNilTargetLang = errors.New("text should not nil")
)

type ErrMessage struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (e *ErrMessage) Error() error {
	return fmt.Errorf("error Message: %s, Details: %s", e.Message, e.Detail)
}
