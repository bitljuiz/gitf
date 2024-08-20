package configs

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

//go:embed language_extensions.json
var langConfig []byte

type LanguageJSON struct {
	Name       string
	Type       string
	Extensions []string
}

func LimitExtensions(files, languages, extensions []string, onlyProgramming bool) ([]string, error) {
	if len(extensions) == 0 && len(languages) == 0 {
		return files, nil
	}
	extensions_ := extensions

	if len(languages) != 0 {
		err := viper.ReadConfig(bytes.NewReader(langConfig))

		var languageJSONs []LanguageJSON

		if err = json.Unmarshal(langConfig, &languageJSONs); err != nil {
			return nil, err
		}

		langMap := make(map[string]int, len(languageJSONs))

		for i, langJSON := range languageJSONs {
			langMap[strings.ToLower(langJSON.Name)] = i
		}

		for _, lang := range languages {
			if i, ok := langMap[lang]; ok {
				extensions_ = append(extensions_, languageJSONs[i].Extensions...)
			}
		}
	}

	extensionsMap := make(map[string]bool, len(extensions_))

	for _, extension := range extensions_ {
		extensionsMap[extension] = true
	}

	files_ := make([]string, 0, len(files))

	for _, file := range files {
		if _, ok := extensionsMap[filepath.Ext(file)]; ok {
			files_ = append(files_, file)
		}
	}

	return files_, nil
}
