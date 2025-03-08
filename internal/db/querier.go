// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CheckRefreshTokenUsed(ctx context.Context, refreshTokenUsed pgtype.Text) (int64, error)
	CheckUserBaseExists(ctx context.Context, userEmail string) (int64, error)
	CreateUserBase(ctx context.Context, arg CreateUserBaseParams) (UserBase, error)
	CreateUserProfile(ctx context.Context, arg CreateUserProfileParams) (UserProfile, error)
	CreateUserSession(ctx context.Context, arg CreateUserSessionParams) (UserSession, error)
	DeleteSessionBySubToken(ctx context.Context, subToken string) error
	DeleteSessionByUserId(ctx context.Context, userID int64) error
	GetSessionByRefreshTokenUsed(ctx context.Context, refreshTokenUsed pgtype.Text) (GetSessionByRefreshTokenUsedRow, error)
	GetSessionBySubToken(ctx context.Context, subToken string) (GetSessionBySubTokenRow, error)
	GetUserBaseByEmail(ctx context.Context, userEmail string) (UserBase, error)
	GetUserBaseById(ctx context.Context, userID int64) (UserBase, error)
	GetUserByUserHash(ctx context.Context, userHash string) (UserBase, error)
	GetUserProfile(ctx context.Context, userID int64) (GetUserProfileRow, error)
	UpdateSession(ctx context.Context, arg UpdateSessionParams) error
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (UserProfile, error)
	UpdateUserVerify(ctx context.Context, userHash string) (UserBase, error)
}

var _ Querier = (*Queries)(nil)
