package db

import (
	"context"

	"github.com/etnperlong/clash2sfa/model"
)

type DB interface {
	GetArg(cxt context.Context, blake3 string) (model.ConvertArg, error)
	PutArg(cxt context.Context, blake3 string, arg model.ConvertArg) error
}
