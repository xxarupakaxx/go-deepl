package deepl

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

func (c *Client) convertLang(lang lang) string {
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
