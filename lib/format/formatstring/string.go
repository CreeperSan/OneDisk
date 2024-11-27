package formatstring

import "fmt"

func String(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
