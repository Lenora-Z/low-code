//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 8:47 下午
package server

import (
	"github.com/Lenora-Z/low-code/service/file"
	"github.com/Lenora-Z/low-code/service/fileHandler"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

type fileStruct struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type fileFormat struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (ds *defaultServer) fileUpload(files []*multipart.FileHeader, appId uint64) (error, []fileStruct) {
	storage := fileHandler.NewFileStorageMinio(ds.minioClient)
	fileSrv := fileHandler.NewFileService(storage)
	path := make([]fileStruct, 0, cap(files))
	srv := file.NewFileService(ds.db)
	for _, item := range files {
		name := item.Filename
		pathItem, hash, err := fileSrv.Upload(item, item.Header.Get("Content-Type"), "default", "low-code", utils.GetRandomStringSec(14))
		if err != nil {
			logrus.Error("file(", name, ") upload failed:", err)
			return err, nil
		}
		md5hash := utils.MD5(hash, utils.GetRandomStringSec(4))
		if err, _ := srv.CreateFileItem(pathItem, name, hash, md5hash, appId); err != nil {
			logrus.Error("put file(", name, ") in db error:", err)
			return err, nil
		}
		path = append(path, fileStruct{
			Name: name,
			Hash: md5hash,
		})
		logrus.Info("file(", name, ") upload success,path:", pathItem)
	}
	return nil, path
}

func (ds *defaultServer) getLanguage(c *gin.Context) string {
	return c.Request.Header.Get("Lang")
}

var validate *validator.Validate

func validatorInstance() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	return validate
}

func (ds *defaultServer) validateClient(moduleId, formId, userId uint64) bool {
	//获取用户组织
	userSrv := user.NewUserService(ds.db)
	status, item := userSrv.GetUser(userId)
	if status {
		logrus.Error("validatePermission:user not found")
		return false
	}
	err, _, roles := userSrv.GetRelationByOrgId(item.GroupID)
	if err != nil {
		logrus.Error("get organization roles failed:", err)
		return false
	}

	status, _ = userSrv.SearchPermission(roles.RoleIds(), moduleId, formId)

	return !status
}
