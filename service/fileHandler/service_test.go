package fileHandler

import (
	"fmt"
	"github.com/minio/minio-go"
	"testing"
)

var client *minio.Client

func init() {
	minioClient, err := minio.New("192.168.3.47:31058", "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", false)
	if err != nil {
		fmt.Errorf("minio: %+v", err)
		return
	}
	client = minioClient
}

func TestFileStorageMinio_DownLoadUrl(t *testing.T) {
	path := "/default/low-code-202104/20210414105611-105621-u=825057118,3516313570&fm=193&f=GIF.jpeg"
	storage := NewFileStorageMinio(client)
	srv := NewFileService(storage)
	url, err := srv.DownLoadUrl(path)
	t.Log(url, err)
}
