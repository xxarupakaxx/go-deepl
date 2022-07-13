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


