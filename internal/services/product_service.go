package services

import "context"

type (
	IProductUser interface {
		ViewProductDetails(ctx context.Context, productId int64) (code int, err error)
	}

	IProductOwner interface {
	}
)