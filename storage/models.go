package storage

const (
	FileTypeFile      = 0
	FileTypeDirectory = 1
)

type File struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Type       int    `json:"type"`
	Size       int64  `json:"size"`
	UpdateTime int64  `json:"update_time"`
}
