package services

import "context"

type (
	IShopUser interface {
		RegisterShop(ctx context.Context)
	}

	IShopAdmin interface{}
)
