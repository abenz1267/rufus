package rufus

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	return nil
}

// Translate phrase to according language
func (t Translation) Translate(s, lang string) string {
	languageIndex := t.Languages[lang]
	return t.Phrases[s][languageIndex]
}

// TranslateURL phrase to URL friendly translation
func (t Translation) TranslateURL(s, lang string) string {
	languageIndex := t.Languages[lang]
	return t.URL[s][languageIndex]
}
