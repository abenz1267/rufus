package rufus

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

// Translation handles all the functionality and data for translating
type Translation struct {
	Amount    int                 `json:"-"`
	Languages map[string]int      `json:"languages,omitempty"`
	Locales   map[string]string   `json:"locales,omitempty"`
	Phrases   map[string][]string `json:"phrases,omitempty"`
	URL       map[string][]string `json:"url,omitempty"`
}

func (t *Translation) loadData() error {
	translationFile, err := ioutil.ReadFile("translation.json")
	if err != nil {
		log.Println("No translation file. Setting LanguageAmount to 1.")
		t.Amount = 1
	} else {
		if err := json.Unmarshal(translationFile, t); err != nil {
			return err
		}

		t.Amount = len(t.Languages)
	}

	if err := t.checkIfDataIsValid(); err != nil {
		return err
	}

	return nil
}

func (t Translation) checkIfDataIsValid() error {
	if len(t.Locales) != t.Amount {
		return errors.New("not enough locales for translation defined")
	}

	if err := t.checkTranslationMaps(t.Phrases, "phrases"); err != nil {
		return err
	}

	return t.checkTranslationMaps(t.URL, "url")
}

func (t Translation) checkTranslationMaps(mapToCheck map[string][]string, translationType string) error {
	var keysWithError []string
	var errorMessage strings.Builder

	for k, v := range mapToCheck {
		if len(v) != t.Amount {
			keysWithError = append(keysWithError, k)
		}
	}

	if keysWithError != nil {
		errorMessage.WriteString("missing translation for ")
		errorMessage.WriteString(translationType)
		errorMessage.WriteString(": ")
		errorMessage.WriteString(strings.Join(keysWithError, ", "))

		return errors.New(errorMessage.String())
	}

	return nil
}

// Translate phrase to according language
func (t Translation) Translate(s, lang string) string {
	if val, ok := t.Phrases[s]; ok {
		return val[t.Languages[lang]]
	}

	return "Phrase can't be translated"
}

// TranslateURL phrase to URL friendly translation
func (t Translation) TranslateURL(s, lang string) string {
	if val, ok := t.URL[s]; ok {
		return val[t.Languages[lang]]
	}

	return "URL can't be translated"
}
