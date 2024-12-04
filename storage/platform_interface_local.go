package storage

import (
	errcode "OneDisk/def/err_code"
	"OneDisk/lib/log"
	"OneDisk/module/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type PlatformInterfaceLocal struct {
	PlatformInterface
	Root string
}

var tag = "PlatformInterfaceLocal"

func convertBoolToFileType(isDir bool) int {
	if isDir {
		return FileTypeDirectory
	}
	return FileTypeFile
}

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
			Name:       tmpFile.Name(),
			Path:       entityPath,
			Size:       tmpFileInfo.Size(),
			Type:       convertBoolToFileType(tmpFileInfo.IsDir()),
			UpdateTime: tmpFileInfo.ModTime().Unix(),
		})
	}
	return tmpFiles, database.OperationResult{Code: errcode.OK}
}

func (storage *PlatformInterfaceLocal) CreateFile(path string) (*File, database.OperationResult) {
	// 1、路径校验
	// 1.1、拼接路径
	tmpFilePath := filepath.Clean(storage.Root + "/" + path)
	tmpRootPath := filepath.Clean(storage.Root)
	// 1.2、校验是否为子目录
	if len(tmpRootPath) < len(tmpFilePath) || strings.Index(tmpRootPath, tmpFilePath) != 0 {
		return nil, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "File path is not exist",
		}
	}
	// 1.3、检查目录是否存在
	tmpFileStat, err := os.Stat(tmpFilePath)
	if tmpFileStat != nil || err == nil {
		return nil, database.OperationResult{
			Code:    errcode.FileCanNotReadDirectory,
			Message: "File path is already exist",
		}
	}
	// 2、遍历创建目录
	err = os.MkdirAll(filepath.Dir(tmpFilePath), os.ModePerm)
	if err != nil {
		return nil, database.OperationResult{
			Code:    errcode.FileCanNotCreateParentDirectory,
			Message: "Can not create parent directory",
		}
	}
	// 3、创建文件
	createFile, err := os.Create(tmpFilePath)
	if err != nil || createFile == nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not create file",
		}
	}
	err = createFile.Close()
	if err != nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not close file",
		}
	}
	// 4、返回结果
	createFileInfo, err := createFile.Stat()
	if err != nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not get file info",
		}
	}
	return &File{
		Name:       createFile.Name(),
		Path:       tmpFilePath,
		Type:       convertBoolToFileType(createFileInfo.IsDir()),
		Size:       createFileInfo.Size(),
		UpdateTime: createFileInfo.ModTime().Unix(),
	}, database.OperationResult{Code: errcode.OK}
}

func (storage *PlatformInterfaceLocal) CreateDirectory(path string) (*File, database.OperationResult) {
	// 1、路径校验
	// 1.1、拼接路径
	tmpFilePath := filepath.Clean(storage.Root + "/" + path)
	tmpRootPath := filepath.Clean(storage.Root)
	// 1.2、校验是否为子目录
	if len(tmpRootPath) < len(tmpFilePath) || strings.Index(tmpRootPath, tmpFilePath) != 0 {
		return nil, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "File path is not exist",
		}
	}
	// 1.3、检查目录是否存在
	tmpFileStat, err := os.Stat(tmpFilePath)
	if tmpFileStat != nil || err == nil {
		return nil, database.OperationResult{
			Code:    errcode.FileCanNotReadDirectory,
			Message: "File path is already exist",
		}
	}
	// 2、遍历创建目录
	err = os.MkdirAll(tmpFilePath, os.ModePerm)
	if err != nil {
		return nil, database.OperationResult{
			Code:    errcode.FileCanNotCreateParentDirectory,
			Message: "Can not create parent directory",
		}
	}
	// 3、返回结果
	createDirectory, err := os.Stat(tmpFilePath)
	if err != nil || createDirectory == nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not create file",
		}
	}
	return &File{
		Name:       createDirectory.Name(),
		Path:       tmpFilePath,
		Type:       convertBoolToFileType(createDirectory.IsDir()),
		Size:       createDirectory.Size(),
		UpdateTime: createDirectory.ModTime().Unix(),
	}, database.OperationResult{Code: errcode.OK}
}

func (storage *PlatformInterfaceLocal) Delete(path string) database.OperationResult {
	// 1、路径校验
	// 1.1、拼接路径
	tmpFilePath := filepath.Clean(storage.Root + "/" + path)
	tmpRootPath := filepath.Clean(storage.Root)
	// 1.2、校验是否为子目录
	if len(tmpRootPath) < len(tmpFilePath) || strings.Index(tmpRootPath, tmpFilePath) != 0 {
		return database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "File path is not exist",
		}
	}
	// 1.3、检查目录是否存在
	tmpFileStat, err := os.Stat(tmpFilePath)
	if tmpFileStat != nil || err == nil {
		return database.OperationResult{
			Code:    errcode.FileCanNotReadDirectory,
			Message: "File path is already exist",
		}
	}
	// 2、删除目录 Or 文件
	err = os.RemoveAll(tmpFilePath)
	if err != nil {
		return database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not delete file",
		}
	}
	// 3、返回结果
	return database.OperationResult{Code: errcode.OK}
}

func (storage *PlatformInterfaceLocal) Move(fromFilePath string, toFilePath string) (*File, database.OperationResult) {
	// 1、路径校验
	// 1.1、拼接路径
	tmpFromFilePath := filepath.Clean(storage.Root + "/" + fromFilePath)
	tmpToPath := filepath.Clean(storage.Root + "/" + toFilePath)
	// 1.2、校验是否为子目录
	if len(tmpToPath) < len(tmpFromFilePath) || strings.Index(tmpToPath, tmpFromFilePath) != 0 {
		return nil, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "File path is not exist",
		}
	}
	if len(tmpToPath) < len(tmpFromFilePath) || strings.Index(tmpToPath, tmpFromFilePath) != 0 {
		return nil, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "File path is not exist",
		}
	}
	// 1.3、检查源文件是否存在
	if _, err := os.Stat(tmpFromFilePath); err != nil {
		return nil, database.OperationResult{
			Code:    errcode.FileNotExist,
			Message: "File is not exist",
		}
	}
	// 1.4、检查目标文件是否存在
	if _, err := os.Stat(tmpToPath); err == nil {
		return nil, database.OperationResult{
			Code:    errcode.FileAlreadyExist,
			Message: "File is already exist",
		}
	}
	// 2、移动文件
	err := os.Rename(tmpFromFilePath, tmpToPath)
	if err != nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not move file",
		}
	}
	// 3、返回结果
	createDirectory, err := os.Stat(tmpToPath)
	if err != nil || createDirectory == nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not move file",
		}
	}
	return &File{
		Name:       createDirectory.Name(),
		Path:       tmpToPath,
		Type:       convertBoolToFileType(createDirectory.IsDir()),
		Size:       createDirectory.Size(),
		UpdateTime: createDirectory.ModTime().Unix(),
	}, database.OperationResult{Code: errcode.OK}
}

func (storage *PlatformInterfaceLocal) Upload(context *gin.Context, requestFile *multipart.FileHeader, path string) (*File, database.OperationResult) {
	// 1、检查 Path 是否存在于目录下
	// 1.1、拼接路径
	tmpFilePath := filepath.Clean(storage.Root + "/" + path)
	tmpRootPath := filepath.Clean(storage.Root)
	// 1.2、校验是否为子目录
	if len(tmpRootPath) < len(tmpFilePath) || strings.Index(tmpRootPath, tmpFilePath) != 0 {
		return nil, database.OperationResult{
			Code:    errcode.ParamsError,
			Message: "File path is not exist",
		}
	}
	// 1.3、检查目录是否存在
	tmpFileStat, err := os.Stat(tmpFilePath)
	if tmpFileStat != nil || err != nil {
		return nil, database.OperationResult{
			Code:    errcode.FileCanNotReadDirectory,
			Message: "File path is already exist",
		}
	}
	// 1.4、检查父目录是否创建，如果没有则创建
	tmpFileParentPath := filepath.Dir(tmpFilePath)
	if _, err := os.Stat(tmpFileParentPath); err != nil {
		err = os.MkdirAll(tmpFileParentPath, os.ModePerm)
		if err != nil {
			return nil, database.OperationResult{
				Code:    errcode.FileCanNotCreateParentDirectory,
				Message: "Can not create parent directory",
			}
		}
	}
	// 2、保存文件
	err = context.SaveUploadedFile(requestFile, path)
	if err != nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not save file",
		}
	}
	// 3、获取文件信息并返回
	tmpFile, err := os.Stat(tmpFilePath)
	if err != nil {
		return nil, database.OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Can not get file info",
		}
	}
	return &File{
		Name:       tmpFile.Name(),
		Path:       path,
		Type:       FileTypeFile,
		Size:       tmpFile.Size(),
		UpdateTime: tmpFile.ModTime().Unix(),
	}, database.OperationResult{Code: errcode.OK}
}
