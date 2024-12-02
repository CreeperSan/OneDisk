package defstorage

type Config struct{}

type ConfigLocalPath struct {
	Config
	Path string `json:"path"`
}
