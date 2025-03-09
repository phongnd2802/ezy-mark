package initialize

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/phongnd2802/ezy-mark/internal/consts"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/rs/zerolog/log"
)

func initMinIO() {
	endpoint := global.Config.MinIOEndPoint
	accessKeyID := global.Config.MinIOAccessKey
	secretAccessKey := global.Config.MinIOSecretKey
	useSSL := global.Config.MinIOUseSSL

	log.Info().Msg("Connecting to MinIO...")
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect MinIO")
	}

	global.Minio = minioClient

	log.Info().Msg("Connected MinIO successfully")

	initBuckets()
}

func initBuckets() {
	buckets := []string{consts.BucketUserAvatar, consts.BucketShopLogo, consts.BucketShopBusinessLicense}

	for _, bucket := range buckets {
		exists, err := global.Minio.BucketExists(context.Background(), bucket)
		if err != nil {
			log.Error().Err(err).Msgf("Error checking if bucket %s exists", bucket)
			continue
		}

		if !exists {
			err = global.Minio.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{
				Region:        "us-east-1",
				ObjectLocking: true,
			})
			if err != nil {
				log.Error().Err(err).Msgf("Error creating bucket %s", bucket)
			} else {
				log.Info().Msgf("Bucket %s created successfully", bucket)
			}
		} else {
			log.Info().Msgf("Bucket %s already exists", bucket)
		}
	}
}
