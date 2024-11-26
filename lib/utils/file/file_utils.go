package fileutils

import "os"

// Exists
// 检查文件是否存在
// check file exists
func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// CreateFile
// 创建文件
// create file
func CreateFile(filePath string) bool {
	// create file
	f, err := os.Create(filePath)
	if err != nil {
		return false
	}
	err = f.Close()
	return err == nil
}

// CreateFileIfNotExist
// 如果文件不存在创建文件
// create file if not exists
func CreateFileIfNotExist(filePath string) bool {
	if Exists(filePath) {
		return true
	}
	return CreateFile(filePath)
}
