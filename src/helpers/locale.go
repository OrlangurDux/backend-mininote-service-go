package helpers

import (
	"net/http"

	"golang.org/x/text/language"
)

type Locale string

const (
	LocaleRu Locale = "ru"
	LocaleEn Locale = "en"
)

var supportedLocales = []struct {
	Tag    language.Tag
	Locale Locale
}{
	{language.Russian, LocaleRu},
	{language.English, LocaleEn},
}

var languageMatcher = func() language.Matcher {
	tags := make([]language.Tag, 0, len(supportedLocales))
	for _, item := range supportedLocales {
		tags = append(tags, item.Tag)
	}

	return language.NewMatcher(tags)
}()

// LocaleFromRequest -> get locale from request
func LocaleFromRequest(req *http.Request) Locale {
	header := req.Header.Get("Accept-Language")
	if header == "" {
		return LocaleEn
	}

	tags, _, err := language.ParseAcceptLanguage(header)
	if err != nil || len(tags) == 0 {
		return LocaleEn
	}

	_, index, _ := languageMatcher.Match(tags...)

	if index < 0 || index >= len(supportedLocales) {
		return LocaleEn
	}

	return supportedLocales[index].Locale
}
