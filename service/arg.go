package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/etnperlong/clash2sfa/db"
	"github.com/etnperlong/clash2sfa/model"
	"github.com/etnperlong/clash2singbox/httputils"
	"lukechampine.com/blake3"
)

func PutArg(cxt context.Context, arg model.ConvertArg, db db.DB) (string, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		return "", fmt.Errorf("PutArg: %w", err)
	}
	hash := blake3.Sum256(b)
	h := hex.EncodeToString(hash[:])
	err = db.PutArg(cxt, h, arg)
	if err != nil {
		return "", fmt.Errorf("PutArg: %w", err)
	}
	return h, nil
}

func GetSub(cxt context.Context, c *http.Client, db db.DB, id string, frontendByte []byte) ([]byte, error) {
	arg, err := db.GetArg(cxt, id)
	if err != nil {
		return nil, fmt.Errorf("GetSub: %w", err)
	}
	if arg.Config == "" && arg.ConfigUrl == "" {
		arg.Config = string(frontendByte)
	}
	if arg.ConfigUrl != "" {
		b, err := httputils.HttpGet(cxt, c, arg.ConfigUrl)
		if err != nil {
			return nil, fmt.Errorf("GetSub: %w", err)
		}
		arg.Config = string(b)
	}
	b, err := convert2sing(cxt, c, arg.Config, arg.Sub, arg.Include, arg.Exclude)
	if err != nil {
		return nil, fmt.Errorf("GetSub: %w", err)
	}
	return b, nil
}
