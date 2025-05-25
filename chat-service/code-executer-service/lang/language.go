package lang

func SupportedLanguages() []string {
	return []string{"python", "go", "nodejs"}
}

func IsSupported(lang string) bool {
	for _, l := range SupportedLanguages() {
		if l == lang {
			return true
		}
	}
	return false
}
