package storage

import (
	errcode "OneDisk/def/err_code"
	"OneDisk/lib/log"
	"OneDisk/module/database"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

type PlatformInterfaceLocal struct {
	PlatformInterface
	Root string
}

var tag = "PlatformInterfaceLocal"

func (storage *PlatformInterfaceLocal) List(path string) ([]File, database.OperationResult) {
	// 1、目录拼装并检查
	var tmpFilePath = storage.Root + "/" + path
	if len(tmpFilePath) <= 0 {
		return []File{}, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "路径不存在",
		}
	}
	// 2、错误校验
	// 2.1 清除路径中的标识符
	tmpFilePath = filepath.Clean(tmpFilePath)
	tmpRootPath := filepath.Clean(storage.Root)
	// 2.2 判断路径是否为 storage.Root 的子目录
	if len(tmpRootPath) < len(tmpFilePath) || strings.Index(tmpRootPath, tmpFilePath) != 0 {
		return []File{}, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "路径不存在",
		}
	}
	// 3、读取目录
	tmpFileStat, err := os.Stat(tmpFilePath)
	if tmpFileStat == nil || err != nil {
		return []File{}, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "路径不存在",
		}
	}
	// 4、读取目录下的文件
	tmpFilePathEntity, err := os.ReadDir(tmpFilePath)
	if err != nil {
		return []File{}, database.OperationResult{
			Code:    errcode.FileCanNotReadDirectory,
			Message: "无法读取路径",
		}
	}
	// 5、遍历组装文件
	var tmpFiles []File
	for _, tmpFile := range tmpFilePathEntity {
		entityPath := filepath.Join(tmpFilePath, tmpFile.Name())
		tmpFileInfo, err := tmpFile.Info()
		if err != nil {
			log.Warming(tag, "Error occur while tmpFile.Info", zap.Error(err))
			continue
		}
		tmpFiles = append(tmpFiles, File{
			Name: tmpFile.Name(),
			Path: entityPath,
			Size: tmpFileInfo.Size(),
			Type: func() int {
				if tmpFile.IsDir() {
					return FileTypeDirectory
				}
				return FileTypeFile
			}(),
			UpdateTime: tmpFileInfo.ModTime().Unix(),
		})
	}
	return tmpFiles, database.OperationResult{Code: errcode.OK}
}
