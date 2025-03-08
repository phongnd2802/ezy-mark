package impl

import (
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/phongnd2802/ezy-mark/internal/cache"
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/phongnd2802/ezy-mark/internal/helpers"
	"github.com/phongnd2802/ezy-mark/internal/mapper"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/utils"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
	"github.com/phongnd2802/ezy-mark/internal/worker"
	"github.com/rs/zerolog/log"
)

type sUserInfo struct {
	store       db.Store
	cache       cache.Cache
	distributor worker.TaskDistributor
}

const (
	BucketUserAvatar string = "user-avatar"
)

// ChangePassword implements services.IUserInfo.
func (s *sUserInfo) ChangePassword(ctx context.Context, params *models.ChangePassword) (code int, err error) {
	panic("unimplemented")
}

// GetUserProfile implements services.IUserInfo.
func (s *sUserInfo) GetUserProfile(ctx context.Context, params *models.GetProfileParams) (int, *models.UserProfile, error) {
	var userProfile db.GetUserProfileRow
	err := s.cache.GetObj(ctx, params.SubToken, &userProfile)

	// Cache Hit
	if err == nil {
		log.Info().Msg("Get User Profile FROM CACHE...")
		return response.ErrCodeSuccess, mapper.MapUserProfile(userProfile), nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		return response.ErrCodeInternalServer, nil, err
	}

	// Cache Miss
	userProfile, err = s.store.GetUserProfile(ctx, params.UserId)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	log.Debug().Msg("Get User Profile FROM DBS...")
	return response.ErrCodeSuccess, mapper.MapUserProfile(userProfile), nil
}

// UpdateUserProfile implements services.IUserInfo.
func (s *sUserInfo) UpdateUserProfile(ctx context.Context, params *models.UpdateProfileUserReq) (int, *models.UpdateProfileUserRes, error) {
	var userBirthday *time.Time
	if params.UserBirthday != "" {
		date, _ := time.Parse("2006-01-02", params.UserBirthday)
		userBirthday = &date
	}

	var userProfile db.GetUserProfileRow
	// Get User Profile From Cache
	err := s.cache.GetObj(ctx, params.SubToken, &userProfile)

	// Inernal Server Error
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		return response.ErrCodeInternalServer, nil, err
	}
	// Cache Miss
	userProfile, err = s.store.GetUserProfile(ctx, params.UserId)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	updateArgs := db.UpdateUserProfileParams{
		UserNickname: params.UserNickname,
		UserFullname: pgtype.Text{String: params.UserFullname, Valid: params.UserFullname != ""},
		UserMobile:   pgtype.Text{String: params.UserMobile, Valid: params.UserMobile != ""},
		UserGender:   db.NullGenderEnum{GenderEnum: db.GenderEnum(*params.UserGender), Valid: params.UserGender != nil},
		UserBirthday: func() pgtype.Date {
			if userBirthday != nil {
				return pgtype.Date{
					Time:             *userBirthday,
					Valid:            true,
					InfinityModifier: pgtype.Finite,
				}
			}
			return pgtype.Date{Valid: false}
		}(),
		UserID: params.UserId,
	}
	if params.UserAvatar != nil {
		// Upload Avatar To MinIO
		objectName := helpers.GenerateFileAvatarName(params.UserId, params.UserAvatar.Filename)
		src, err := params.UserAvatar.Open()
		if err != nil {
			return response.ErrCodeInternalServer, nil, err
		}
		defer src.Close()
		_, err = global.Minio.PutObject(ctx, BucketUserAvatar, objectName, src, params.UserAvatar.Size, minio.PutObjectOptions{
			ContentType: params.UserAvatar.Header.Get("Content-Type"),
		})
		if err != nil {
			return response.ErrCodeInternalServer, nil, err
		}

		updateArgs.UserAvatar.String = utils.GetURLFileName(BucketUserAvatar, objectName)
		updateArgs.UserAvatar.Valid = true

		if userProfile.UserAvatar.Valid {
			// Remove Old Avatar
			payload := &worker.PayloadRemoveOldAvatar{
				ObjectName: utils.ExtractFileNameFromURL(userProfile.UserAvatar.String),
			}

			opts := []asynq.Option{
				asynq.MaxRetry(3),
				asynq.ProcessIn(3 * time.Second),
				asynq.Queue(worker.QueueDefault),
			}

			err := s.distributor.DistributeTaskRemoveOldAvatar(ctx, payload, opts...)
			if err != nil {
				log.Error().Err(err).Msg("Failed to distribute task remove old avatar")
			}
		}
	}

	updatedProfile, err := s.store.UpdateUserProfile(ctx, updateArgs)
	if err != nil {
		return response.ErrCodeInternalServer, nil, err
	}

	return response.ErrCodeSuccess, &models.UpdateProfileUserRes{
		UserProfile: models.UserProfile{
			UserNickname: updatedProfile.UserNickname,
			UserFullname: updatedProfile.UserFullname.String,
			UserMobile:   updatedProfile.UserMobile.String,
			UserGender:   (*string)(&updatedProfile.UserGender.GenderEnum),
			UserBirthday: updatedProfile.UserBirthday.Time.Format("2006-01-02"),
		},
		UserAvatar: updatedProfile.UserAvatar.String,
	}, nil

}

func NewUserInfoImpl(store db.Store, cache cache.Cache, distributor worker.TaskDistributor) *sUserInfo {
	return &sUserInfo{
		store: store,
		cache: cache,
		distributor: distributor,
	}
}

var _ services.IUserInfo = (*sUserInfo)(nil)
