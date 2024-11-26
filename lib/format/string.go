package format

import "fmt"

func String(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
