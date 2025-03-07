package impl

import (
	"context"
	"errors"

	"github.com/phongnd2802/ezy-mark/internal/cache"
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/mapper"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
	"github.com/rs/zerolog/log"
)

type sUserInfo struct {
	store db.Store
	cache cache.Cache
}

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
func (s *sUserInfo) UpdateUserProfile(ctx context.Context, params *models.UpdateProfileUserReq) (code int, res *models.UpdateProfileUserRes, err error) {
	panic("unimplemented")
}

func NewUserInfoImpl(store db.Store, cache cache.Cache) *sUserInfo {
	return &sUserInfo{
		store: store,
		cache: cache,
	}
}

var _ services.IUserInfo = (*sUserInfo)(nil)
