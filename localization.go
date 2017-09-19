package sprig

// Translator interface
type Translator interface {
	T(string) string
}

var locale = InitLocale("EN")

// Localization type for serving the
// localization functions,
// implements Translator
type Localization struct {
	locale       string
	dictionaries map[string]map[string]string
}

// Localization initialization
func InitLocale(lang string) Localization {
	var localization = Localization{}
	localization.SetLocale(lang)

	return localization
}

// Sets default locale
func (l *Localization) SetLocale(lang string) {
	l.locale = lang
}

// Returns default locale
func (l *Localization) GetLocale() string {
	return l.locale
}

// Sets localization dictionary
func (l *Localization) SetTranslation(locale string, dictionary map[string]string) Localization {
	if l.dictionaries == nil {
		l.dictionaries = make(map[string]map[string]string)
	}
	l.dictionaries[locale] = make(map[string]string)
	l.dictionaries[locale] = dictionary
	return *l
}

// Translation function
func (l *Localization) T(word string) string {
	if l.dictionaries[l.locale] != nil {
		if _, ok := l.dictionaries[l.locale][word]; ok {
			return l.dictionaries[l.locale][word]
		}
	}

	return ""
}
