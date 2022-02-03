package gtsasamiserver

import (
	"context"

	"github.com/snowzach/queryp"
)

type GTStore interface {
	AccountGetByID(ctx context.Context, id string) (*Account, error)
	AccountSave(ctx context.Context, account *Account) error
	AccountDeleteByID(ctx context.Context, id string) error
	AccountsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*Account, int64, error)
}
