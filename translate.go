package deepl

import (
	"context"
	"net/http"
	"net/url"
)

type lang int32

const (
	Bulgarian lang = iota + 1
	Czech
	Danish
	German
	Greek
	English
	Spanish
	Estonian
	Finnish
	French
	Hungarian
	Indonesian
	Italian
	Japanese
	Lithuanian
	Latvian
	Dutch
	Polish
	Portuguese
	Romanian
	Russian
	Slovak
	Slovenian
	Swedish
	Turkish
	Chinese
)

type splitSentence string

const (
	NoSplit    splitSentence = "0"
	Default                  = "1"
	Nonewlines               = "nonewlines"
)

type preserveFormatting string

const (
	NoPreserveFormat preserveFormatting = "0"
	PreserveFormat                      = "1"
)

type TranslateParams struct {
	Text               string
	SourceLang         lang
	TargetLang         string
	SplitSentences     splitSentence
	PreserveFormatting preserveFormatting
}
type response struct {
	Translations []translation `json:"translations"`
}

type translation struct {
	Language string `json:"detected_source_language"`
	Text     string `json:"text"`
}

