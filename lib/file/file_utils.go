package file

import "os"

// 检查文件是否存在
// check file exists
func fileUtilsExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// 创建文件
// create file
func fileUtilsCreateFile(filePath string) bool {
	// create file
	f, err := os.Create(filePath)
	if err != nil {
		return false
	}
	err = f.Close()
	return err == nil
}

// 如果文件不存在创建文件
// create file if not exists
func fileUtilsCreateFileIfNotExist(filePath string) bool {
	if fileUtilsExists(filePath) {
		return true
	}
	return fileUtilsCreateFile(filePath)
}
