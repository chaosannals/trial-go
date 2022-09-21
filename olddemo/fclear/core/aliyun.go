package core

import (
	"os"
	"fmt"
	"errors"
	"strings"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliOssClient struct {
	client *oss.Client
	bucket *oss.Bucket
	objectACL oss.Option
	storageType oss.Option
	config *FClearAliyunOssConfig
}

func MatchAliOssClient(cfgs []FClearAliyunOssConfig, tag string) (*AliOssClient, error) {
	for i:= 0; i < len(cfgs); i+=1 {
		cfg := &cfgs[i]
		if strings.EqualFold(cfg.Tag, tag) {
			return NewAliOssClient(cfg)
		}
	}
	return nil, errors.New(fmt.Sprintf("没有匹配的配置： %s", tag))
}

func NewAliOssClient(cfg *FClearAliyunOssConfig) (*AliOssClient, error) {
	client, err := oss.New(cfg.Endpoint, cfg.Key, cfg.Secret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, err
	}

	var st oss.Option
	if cfg.IsStorageIA {
		st = oss.ObjectStorageClass(oss.StorageIA)
	} else {
		st = oss.ObjectStorageClass(oss.StorageStandard)
	}

	var oacl oss.Option
	if cfg.IsPrivate {
		oacl = oss.ObjectACL(oss.ACLPrivate)
	} else {
		oacl = oss.ObjectACL(oss.ACLDefault)
	}

	return &AliOssClient{
		client: client,
		bucket: bucket,
		objectACL: oacl,
		storageType: st,
		config: cfg,
	}, nil
}

func (i *AliOssClient)PutFromFile(osspath string, localpath string) error {
	src, err := os.Open(localpath)
	if err != nil {
		return err
	}
	dist := fmt.Sprintf("%s/%s", i.config.Prefix, osspath)
	return i.bucket.PutObject(dist, src, i.storageType, i.objectACL)
}
