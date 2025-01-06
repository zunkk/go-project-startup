package util

import "unicode"

// CleanInput removes non-printable and invalid UTF-8 characters from the input string
func CleanInput(input string) string {
	cleanedRunes := make([]rune, 0, len(input))
	for _, r := range input {
		if unicode.IsPrint(r) {
			cleanedRunes = append(cleanedRunes, r)
		}
	}
	return string(cleanedRunes)
}
