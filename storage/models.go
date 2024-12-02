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
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
