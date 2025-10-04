package service

import (
	"fmt"
	"strings"
)

// LANGUAGES maps language codes to language names
var LANGUAGES = map[string]string{
	"en":  "english",
	"zh":  "chinese",
	"de":  "german",
	"es":  "spanish",
	"ru":  "russian",
	"ko":  "korean",
	"fr":  "french",
	"ja":  "japanese",
	"pt":  "portuguese",
	"tr":  "turkish",
	"pl":  "polish",
	"ca":  "catalan",
	"nl":  "dutch",
	"ar":  "arabic",
	"sv":  "swedish",
	"it":  "italian",
	"id":  "indonesian",
	"hi":  "hindi",
	"fi":  "finnish",
	"vi":  "vietnamese",
	"he":  "hebrew",
	"uk":  "ukrainian",
	"el":  "greek",
	"ms":  "malay",
	"cs":  "czech",
	"ro":  "romanian",
	"da":  "danish",
	"hu":  "hungarian",
	"ta":  "tamil",
	"no":  "norwegian",
	"th":  "thai",
	"ur":  "urdu",
	"hr":  "croatian",
	"bg":  "bulgarian",
	"lt":  "lithuanian",
	"la":  "latin",
	"mi":  "maori",
	"ml":  "malayalam",
	"cy":  "welsh",
	"sk":  "slovak",
	"te":  "telugu",
	"fa":  "persian",
	"lv":  "latvian",
	"bn":  "bengali",
	"sr":  "serbian",
	"az":  "azerbaijani",
	"sl":  "slovenian",
	"kn":  "kannada",
	"et":  "estonian",
	"mk":  "macedonian",
	"br":  "breton",
	"eu":  "basque",
	"is":  "icelandic",
	"hy":  "armenian",
	"ne":  "nepali",
	"mn":  "mongolian",
	"bs":  "bosnian",
	"kk":  "kazakh",
	"sq":  "albanian",
	"sw":  "swahili",
	"gl":  "galician",
	"mr":  "marathi",
	"pa":  "punjabi",
	"si":  "sinhala",
	"km":  "khmer",
	"sn":  "shona",
	"yo":  "yoruba",
	"so":  "somali",
	"af":  "afrikaans",
	"oc":  "occitan",
	"ka":  "georgian",
	"be":  "belarusian",
	"tg":  "tajik",
	"sd":  "sindhi",
	"gu":  "gujarati",
	"am":  "amharic",
	"yi":  "yiddish",
	"lo":  "lao",
	"uz":  "uzbek",
	"fo":  "faroese",
	"ht":  "haitian creole",
	"ps":  "pashto",
	"tk":  "turkmen",
	"nn":  "nynorsk",
	"mt":  "maltese",
	"sa":  "sanskrit",
	"lb":  "luxembourgish",
	"my":  "myanmar",
	"bo":  "tibetan",
	"tl":  "tagalog",
	"mg":  "malagasy",
	"as":  "assamese",
	"tt":  "tatar",
	"haw": "hawaiian",
	"ln":  "lingala",
	"ha":  "hausa",
	"ba":  "bashkir",
	"jw":  "javanese",
	"su":  "sundanese",
	"yue": "cantonese",
}

// TO_LANGUAGE_CODE maps language names to their ISO codes
var TO_LANGUAGE_CODE = map[string]string{
	// Common language aliases
	"burmese":       "my",
	"valencian":     "ca",
	"flemish":       "nl",
	"haitian":       "ht",
	"letzeburgesch": "lb",
	"pushto":        "ps",
	"panjabi":       "pa",
	"moldavian":     "ro",
	"moldovan":      "ro",
	"sinhalese":     "si",
	"castilian":     "es",
	"mandarin":      "zh",
}

// Initialize the reverse mapping
func init() {
	// Populate TO_LANGUAGE_CODE with reversed entries from LANGUAGES
	for code, language := range LANGUAGES {
		TO_LANGUAGE_CODE[language] = code
	}
}

// GetLanguageCode converts a language name to its ISO code
// Example: "indonesian" -> "id"
func GetLanguageCode(languageName string) (string, bool) {
	// Convert to lowercase for case-insensitive matching
	languageName = strings.ToLower(languageName)

	// Look up the language code
	code, ok := TO_LANGUAGE_CODE[languageName]
	if !ok {
		code = languageName
	}
	fmt.Println("code", code, languageName)

	if code == "" {
		return code, false
	}
	return code, true
}
