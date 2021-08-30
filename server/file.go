package server

import (
	"github.com/Lenora-Z/low-code/service/file"
	"github.com/Lenora-Z/low-code/service/fileHandler"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strings"
)

// @Summary 上传文件
// @Tags file
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param X-content-type header string true "图片编码格式(e.g.image/jpeg)"
// @param app-id header string true "应用id"
// @param file[] body object true "文件"
// @Success 200 {object} ApiResponse{result=[]string}
// @Router /file/upload [post]
func (ds *defaultServer) FileUpload(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	forms, _ := ctx.MultipartForm()
	if forms == nil {
		ds.InvalidParametersError(ctx)
		return
	} else {
		err, path := ds.uploadFile(forms, claims.AppId)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		ds.ResponseSuccess(ctx, ds.getFullUrl(path.Hash))
	}
}

// ClientFileUpload
// @Summary 上传文件
// @Tags api/file
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param file[] body object true "文件"
// @Success 200 {object} ApiResponse{result=[]string}
// @Router /api/file/upload [post]
func (ds *defaultServer) ClientFileUpload(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	forms, _ := ctx.MultipartForm()
	if forms == nil {
		ds.InvalidParametersError(ctx)
		return
	} else {
		err, path := ds.uploadFile(forms, claims.AppId)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		ds.ResponseSuccess(ctx, ds.getFullUrl(path.Hash))
	}
}

func (ds *defaultServer) uploadFile(header *multipart.Form, appId uint64) (error, *fileStruct) {
	err, path := ds.fileUpload(header.File["file[]"], appId)
	if err != nil {
		return err, nil
	}
	if len(path) <= 0 {
		return nil, nil
	}
	return nil, &path[0]
}

func (ds *defaultServer) FilePreView(ctx *gin.Context) {
	hash := ctx.Query("code")

	storage := fileHandler.NewFileStorageMinio(ds.minioClient)
	fileSrv := fileHandler.NewFileService(storage)
	srv := file.NewFileService(ds.db)
	_, item := srv.GetFileByHash(hash)
	if err := fileSrv.View(ctx, item.Path); err != nil {
		ds.InternalServiceError(ctx, err.Error())
	}

}

func (ds *defaultServer) getFullUrl(path string) string {
	if strings.Contains(path, "http://") || strings.Contains(path, "https://") {
		return path
	}
	url := "http://" + ds.conf.BaseUrl + "/file/show?code=" + path
	return url
}
