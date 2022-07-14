package deepl

import "errors"

var (
	ErrNilText       = errors.New("text should not nil")
	ErrNilTargetLang = errors.New("text should not nil")
)

type ErrMessage struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}