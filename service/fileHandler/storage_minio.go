package fileHandler

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type fileStorageMinio struct {
	DefaultClient *minio.Client
}

//storagePath = /{clusterName}/{buctketName}/{objectName}
//buctketName:namespace-yearMonth
//objectName:prefixName-hourMinSec-filename
//storagePath = /{clusterName}/{buctketName}/{objectName}
//buctketName:namespace-yearMonth
//objectName:prefixName-hourMinSec-filename
func (fsm *fileStorageMinio) Upload(fileHeader *multipart.FileHeader, xContentType, clusterName, namespace, prefixName string) (storagePath, hash string, err error) {

	var miniClient *minio.Client = nil
	if clusterName == "default" {
		miniClient = fsm.DefaultClient
	}
	if miniClient == nil {
		return "", "", fmt.Errorf("clusterName is invalid")
	}

	// 检查存储桶是否已经存在
	buctketName := namespace + "-" + time.Now().Format("200601")
	exists, err := miniClient.BucketExists(buctketName)
	if err != nil {
		logrus.Error(fmt.Errorf("minioClient.BucketExists(): %s", err.Error()))
		return "", "", fmt.Errorf("minioClient.BucketExists(): %s", err.Error())
	}
	if exists == false {
		// 创建一个叫xxxx的存储桶
		err = miniClient.MakeBucket(buctketName, "")
		if err != nil {
			logrus.Error(fmt.Errorf("Make Bucket:%s", err.Error()))
			return "", "", fmt.Errorf("Make Bucket:%s", err.Error())
		}
	}

	src, err1 := fileHeader.Open()
	if err1 != nil {
		logrus.Error(fmt.Errorf("Failed upload file, fileHeaderOpen: %s", err.Error()))
		return "", "", fmt.Errorf("Failed upload file, fileHeaderOpen: %s", err.Error())
	}
	defer src.Close()

	// 计算文件hash
	h := sha256.New()
	if _, err := io.Copy(h, src); err != nil {
		logrus.Error(fmt.Errorf("Failed hash256: %s", err.Error()))
		return "", "", fmt.Errorf("Failed hash256: %s", err.Error())
	}
	hash256 := fmt.Sprintf("%x", h.Sum(nil))
	logrus.Println(hash256)

	objectName := prefixName + "-" + time.Now().Format("150405") + "-" + fileHeader.Filename
	_, err2 := miniClient.PutObject(buctketName, objectName, src, fileHeader.Size, minio.PutObjectOptions{ContentType: xContentType})
	if err2 != nil {
		logrus.Error(fmt.Errorf("Failed upload file, PutObject: %s", err2.Error()))
		return "", "", fmt.Errorf("Failed upload file, PutObject: %s", err2.Error())
	}

	return fsm.GenStoragePath(clusterName, buctketName, objectName), hash256, nil
}

func (fsm *fileStorageMinio) GenStoragePath(clusterName, bucketName, objectName string) (storagePath string) {
	return "/" + clusterName + "/" + bucketName + "/" + objectName
}

func (fsm *fileStorageMinio) CutStoragePath(storagePath string) (clusterName, bucketName, objectName string) {
	// 按照分隔符'/'对路径进行分解，获取平台名称platformName
	pathArray := strings.Split(storagePath, "/")
	length := len(pathArray)
	clusterName = pathArray[1]
	objectName = pathArray[length-1]
	bucketName = strings.Join(pathArray[2:length-1], "/")
	return
}
func (fsm *fileStorageMinio) Download(ctx *gin.Context, storagePath string) (err error) {
	if storagePath == "" {
		return errors.New("no filePath")
	}
	clusterName, bucketName, objectName := fsm.CutStoragePath(storagePath)
	var miniClient *minio.Client = nil
	if clusterName == "default" {
		miniClient = fsm.DefaultClient
	}
	if miniClient == nil {
		logrus.Error(fmt.Errorf("clusterName is invalid"))
		return fmt.Errorf("clusterName is invalid")
	}

	fileInfo, err := miniClient.StatObject(bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		logrus.Error(fmt.Errorf("fileInfo failed: %s", err.Error()))
		return fmt.Errorf("fileInfo failed: %s", err.Error())
	}

	src, err := miniClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logrus.Error(fmt.Errorf("doawload upload file: %s", err.Error()))
		return fmt.Errorf("doawload upload file: %s", err.Error())
	}
	defer src.Close()

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + objectName + `"`,
	}
	ctx.DataFromReader(http.StatusOK, fileInfo.Size, fileInfo.ContentType, src, extraHeaders)
	return nil
}

func (fsm *fileStorageMinio) View(ctx *gin.Context, storagePath string) (err error) {
	if storagePath == "" {
		return errors.New("no filePath")
	}
	clusterName, bucketName, objectName := fsm.CutStoragePath(storagePath)
	var miniClient *minio.Client = nil
	if clusterName == "default" {
		miniClient = fsm.DefaultClient
	}
	if miniClient == nil {
		logrus.Error(fmt.Errorf("clusterName is invalid"))
		return fmt.Errorf("clusterName is invalid")
	}

	fileInfo, err := miniClient.StatObject(bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		logrus.Error(fmt.Errorf("fileInfo failed: %s", err.Error()))
		return fmt.Errorf("fileInfo failed: %s", err.Error())
	}

	src, err := miniClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logrus.Error(fmt.Errorf("doawload upload file: %s", err.Error()))
		return fmt.Errorf("doawload upload file: %s", err.Error())
	}
	defer src.Close()

	rangeString := ctx.GetHeader("Range")
	if rangeString != "" {
		// 分块传输
		ranges, err := parseRange(rangeString, fileInfo.Size)
		if err != nil {
			return fmt.Errorf("parseRange: %s", err.Error())
		}
		start := ranges[0].start
		length := ranges[0].length
		startByte, err := src.Seek(start, 0)
		if err != nil {
			return fmt.Errorf("src.Seek: %s", err.Error())
		}
		logrus.Info("start byte:", startByte)
		extraHeaders := map[string]string{
			"Content-Range": getRange(start, start+length-1, fileInfo.Size),
		}
		ctx.DataFromReader(http.StatusPartialContent, length, fileInfo.ContentType, src, extraHeaders)
		return nil
	} else {
		ctx.DataFromReader(http.StatusOK, fileInfo.Size, fileInfo.ContentType, src, nil)
		return nil
	}

}
func (fsm *fileStorageMinio) Delete(storagePath string) (err error) {
	clusterName, bucketName, objectName := fsm.CutStoragePath(storagePath)
	var miniClient *minio.Client = nil
	if clusterName == "default" {
		miniClient = fsm.DefaultClient
	}
	if miniClient == nil {
		logrus.Error(fmt.Errorf("clusterName is invalid"))
		return fmt.Errorf("clusterName is invalid")
	}

	return miniClient.RemoveObject(bucketName, objectName)
}

func (fsm *fileStorageMinio) Stat(storagePath string) (hash string, err error) {
	clusterName, bucketName, objectName := fsm.CutStoragePath(storagePath)
	var minioClient *minio.Client = nil
	if clusterName == "default" {
		minioClient = fsm.DefaultClient
	}
	if minioClient == nil {
		logrus.Error(fmt.Errorf("clusterName is invalid"))
		return "", fmt.Errorf("clusterName is invalid")
	}

	fileInfo, err := minioClient.StatObject(bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		logrus.Error(fmt.Errorf("fileInfo failed: %s", err.Error()))
		return "", fmt.Errorf("fileInfo failed: %s", err.Error())
	}
	return fileInfo.ETag, nil
}

func (fsm *fileStorageMinio) FileUpload(filePath, clusterName, namespace, prefixName, contentType string) (storagePath string, err error) {

	var miniClient *minio.Client = nil
	if clusterName == "default" {
		miniClient = fsm.DefaultClient
	}
	if miniClient == nil {
		logrus.Error(fmt.Errorf("clusterName is invalid"))
		return "", fmt.Errorf("clusterName is invalid")
	}

	// 检查存储桶是否已经存在
	buctketName := namespace + "-" + time.Now().Format("200601")
	exists, err := miniClient.BucketExists(buctketName)
	if err != nil {
		logrus.Error(fmt.Errorf("minioClient.BucketExists(): %s", err.Error()))
		return "", fmt.Errorf("minioClient.BucketExists(): %s", err.Error())
	}
	if exists == false {
		// 创建一个叫xxxx的存储桶
		err = miniClient.MakeBucket(buctketName, "")
		if err != nil {
			logrus.Error(fmt.Errorf("Make Bucket:%s", err.Error()))
			return "", fmt.Errorf("Make Bucket:%s", err.Error())
		}
	}

	objectName := prefixName + "-" + time.Now().Format("150406") + "-" + path.Base(filePath)
	_, err2 := miniClient.FPutObject(buctketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err2 != nil {
		logrus.Error(fmt.Errorf("Failed upload file, PutObject: %s", err2.Error()))
		return "", fmt.Errorf("Failed upload file, PutObject: %s", err2.Error())
	}

	return fsm.GenStoragePath(clusterName, buctketName, objectName), nil
}

func (fsm *fileStorageMinio) DownLoadUrl(storagePath string) (string, error) {
	clusterName, bucketName, objectName := fsm.CutStoragePath(storagePath)
	var miniClient *minio.Client = nil
	if clusterName == "default" {
		miniClient = fsm.DefaultClient
	}
	if miniClient == nil {
		logrus.Error(fmt.Errorf("clusterName is invalid"))
		return "", fmt.Errorf("clusterName is invalid")
	}

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", `"attachment; filename="`+objectName+"`")

	// Generates a presigned url which expires in a day.
	presignedURL, err := miniClient.PresignedGetObject(bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		logrus.Error(fmt.Errorf("doawload upload file: %s", err.Error()))
		return "", fmt.Errorf("doawload upload file: %s", err.Error())
	}
	return presignedURL.String(), err
}
