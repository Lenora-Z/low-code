package file

import (
	"github.com/jinzhu/gorm"
)

type FileService interface {
	// CreateFileItem 创建文件
	CreateFileItem(path, name, fileHash, hash string, appId uint64) (error, *File)
	//获取文件by hash值
	GetFileByHash(hash string) (bool, *File)
}

type fileService struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) FileService {
	u := new(fileService)
	u.db = db
	return u
}

func (srv *fileService) CreateFileItem(path, name, fileHash, hash string, appId uint64) (error, *File) {
	return createFile(srv.db, File{
		Path:     path,
		Name:     name,
		IsDelete: &FALSE,
		AppID:    appId,
		FileHash: fileHash,
		Hash:     hash,
	})
}

func (srv *fileService) GetFileByHash(hash string) (bool, *File) {
	column := make(map[string]interface{})
	column["hash"] = hash
	return fileDetailByColumn(srv.db, column)
}
