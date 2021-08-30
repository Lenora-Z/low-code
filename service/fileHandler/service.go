//Created by Goland
//@User: lenora
//@Date: 2021/1/18
//@Time: 2:49 下午
package fileHandler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

type FileStorage interface {
	DownloadFile(storagePath, name string) (error, *minio.ObjectInfo, *minio.Object, map[string]string)
	DownLoadUrl(storagePath string) (string, error)
	Upload(fileHeader *multipart.FileHeader, xContentType, clusterName, namespace, prefixName string) (storagePath, hash string, err error)
	View(ctx *gin.Context, storagePath string) (err error)
}

func NewFileStorageMinio(minoClient *minio.Client) FileStorage {
	fms := new(fileStorageMinio)
	fms.DefaultClient = minoClient
	return fms
}

type FileService interface {
	FileStorage
}

type fileService struct {
	FileStorage
}

func NewFileService(storage FileStorage) FileService {
	fs := new(fileService)
	fs.FileStorage = storage
	return fs
}

func (h *fileStorageMinio) DownloadFile(storagePath, name string) (error, *minio.ObjectInfo, *minio.Object, map[string]string) {
	clusterName, bucketName, objectName := h.CutStoragePath(storagePath)
	var miniClient *minio.Client = nil
	if clusterName == "default" {
		miniClient = h.DefaultClient
	}
	if miniClient == nil {
		logrus.Error(fmt.Errorf("clusterName is invalid"))
		return errors.New("clusterName is invalid"), nil, nil, nil
	}

	fileInfo, err := miniClient.StatObject(bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		logrus.Error(fmt.Errorf("fileInfo failed: %s", err.Error()))
		return err, nil, nil, nil
	}

	src, err := miniClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logrus.Error(fmt.Errorf("doawload upload file: %s", err.Error()))
		return err, nil, nil, nil
	}

	if name != "" {
		objectName = name
	}

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + objectName + `"`,
	}
	return nil, &fileInfo, src, extraHeaders
}
