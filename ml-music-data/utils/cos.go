package utils

import (
	"context"
	"github.com/sta-golang/music-recommend/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
)

var cosClient *cos.Client
var cosUrl string

func InitCosClient(cfg config.COSConfig) error {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return err
	}
	cosUrl = cfg.URL
	b := &cos.BaseURL{BucketURL: u}
	cosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
	return nil
}

func CosPutObject(ctx context.Context, path string, reader io.Reader, op *cos.ObjectPutOptions) error {
	_, err := cosClient.Object.Put(ctx, path, reader, op)
	if err != nil {
		return err
	}
	return nil
}

func CosURL() string {
	return cosUrl
}