package impl

import (
	"context"
	"mime/multipart"

	minio "github.com/minio/minio-go/v7"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/phongnd2802/ezy-mark/internal/services"
)

type sUploadService struct{}

// UploadFile implements services.IUploadService.
func (s *sUploadService) UploadFile(ctx context.Context, bucketName string, objectName string, file *multipart.FileHeader, opts minio.PutObjectOptions) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = global.Minio.PutObject(ctx, bucketName, objectName, src, file.Size, opts)
	if err != nil {
		return err
	}
	
	return nil
}

func NewUploadServiceImpl() *sUploadService {
	return &sUploadService{}
}

var _ services.IUploadService = (*sUploadService)(nil)
