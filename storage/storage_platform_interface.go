package storage

import "OneDisk/module/database"

type PlatformInterface interface {

	// List
	// 列出指定路径下的文件
	List(path string) ([]File, database.OperationResult)

	// CreateFile
	// 创建一个文件
	CreateFile(path string) (*File, database.OperationResult)

	// CreateDirectory
	// 创建一个目录
	CreateDirectory(path string) (*File, database.OperationResult)

	// Delete
	// 删除一个文件或目录
	Delete(path string) database.OperationResult

	// Move
	// 移动一个文件或目录
	Move(path string, newPath string) (*File, database.OperationResult)
}
