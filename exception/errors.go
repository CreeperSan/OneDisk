package exception

import "OneDisk/lib/format/formatstring"

type InterruptException struct {
	Message       string
	Code          int
	OriginalError error
}

func (e *InterruptException) Error() string {
	return formatstring.String("Code:%d, Message=%s, OriginalError=%v", e.Code, e.Message, e.OriginalError)
}
