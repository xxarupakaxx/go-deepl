package deepl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
	TargetLang         lang
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

func (c *Client) Translate(params TranslateParams) (string, error) {
	u, err := url.Parse(c.baseURL.String() + "translate")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header = c.header

	params.Text = strings.Replace(params.Text, "。", ".", -1)

	q := req.URL.Query()
	if params.Text == "" {
		return "", ErrNilText
	}
	if params.TargetLang == 0 {
		return "", ErrNilTargetLang
	}

	q.Add("target_lang", convertLang(params.TargetLang))
	q.Add("text", params.Text)

	if params.SourceLang != 0 {
		q.Add("source_lang", convertLang(params.SourceLang))
	}
	if params.PreserveFormatting != "" {
		q.Add("preserve_formatting", string(params.PreserveFormatting))
	}
	if params.SplitSentences != "" {
		q.Add("split_sentences", string(params.SplitSentences))
	}

	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusRequestEntityTooLarge {
		splitText := strings.Split(params.Text, ".")
		for i, s := range splitText {
			if i == 0 {
				q.Set("text", s)
			} else {
				q.Add("text", s)
			}
		}

		req.URL.RawQuery = q.Encode()

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		var errMessage ErrMessage
		if err = json.Unmarshal(body, &errMessage); err != nil {
			return "", err
		}

		return "", fmt.Errorf("%s", errMessage.DisplayMessage())
	}

	var data response
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	if len(data.Translations) >= 2 {
		var text string
		for _, t := range data.Translations {
			text += t.Text
		}

		return text, nil
	}

	return data.Translations[0].Text, nil
}

func convertLang(lang lang) string {
	switch lang {
	case Bulgarian:
		return "BG"
	case Chinese:
		return "ZH"
	case Czech:
		return "CS"
	case Danish:
		return "DA"
	case Dutch:
		return "NL"
	case English:
		return "EN"
	case Estonian:
		return "ET"
	case Finnish:
		return "FI"
	case French:
		return "FR"
	case German:
		return "DE"
	case Greek:
		return "EL"
	case Hungarian:
		return "HU"
	case Indonesian:
		return "ID"
	case Italian:
		return "IT"
	case Japanese:
		return "JA"
	case Latvian:
		return "LV"
	case Lithuanian:
		return "LT"
	case Polish:
		return "PL"
	case Portuguese:
		return "PT"
	case Romanian:
		return "RO"
	case Russian:
		return "RU"
	case Slovak:
		return "SK"
	case Slovenian:
		return "SL"
	case Spanish:
		return "ES"
	case Swedish:
		return "SV"
	case Turkish:
		return "TR"
	default:
		return ""
	}
}
