package lines

import (
	"strings"
)

func Parser(input string) []string {
	var args []string
	var currentArg strings.Builder

	inQuotes := false
	escapeNext := false
	inSimpleQuotes := false

	for _, r := range input {
		if escapeNext {
			currentArg.WriteRune(r)
			escapeNext = false
			continue
		}

		if r == '\\' {
			escapeNext = true
			currentArg.WriteRune(r)
			continue
		}

		if r == '"' {
			inQuotes = !inQuotes
			currentArg.WriteRune(r)
			continue
		}

		if r == '\'' {
			inSimpleQuotes = !inSimpleQuotes
			continue
		}

		if r == ' ' && !inQuotes && !inSimpleQuotes {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}

		currentArg.WriteRune(r)
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}
