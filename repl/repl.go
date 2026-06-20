package repl

import "strings"

func CleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
