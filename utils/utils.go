package utils

import (
	"bytes"
	"fmt"
	"os"
)

func print(args ...any) (int, error) {
	w := os.Stdout
	var buf bytes.Buffer

	for i, arg := range args {
		if i > 0 {
			buf.WriteByte(' ')
		}
		_, err := fmt.Fprint(&buf, arg)
		if err != nil {
			return 0, err
		}
	}

	if buf.Len() > 0 {
		buf.WriteByte('\n')
		return w.Write(buf.Bytes())
	}

	return 0, nil
}

func System(args ...any) (int, error) {
	return print(append([]any{"\033[35m|System>\033[0m"}, args...)...)
}
