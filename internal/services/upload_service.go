package services

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type IUploadService interface {
	UploadFile(
		ctx context.Context,
		bucketName string,
		objectName string,
		file *multipart.FileHeader,
		opts minio.PutObjectOptions,
	) error
}

var (
	localUploadService IUploadService
)

func UploadService() IUploadService {
	if localUploadService == nil {
		panic("IUploadService interface not implemented yet")
	}
	return localUploadService
}

func InitUploadService(i IUploadService) {
	localUploadService = i
}
