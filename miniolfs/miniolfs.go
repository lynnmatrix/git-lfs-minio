package miniolfs

import (
	"log"
	"net/url"
	"time"
	
	"context"
	minio "github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

type MinioLFSInitParams struct {
	Host       string
	AccessKey  string
	SecretKey  string
	Bucket     string
	URLExpires uint64
}

type MinioLFS struct {
	ctx				 context.Context
	api        *minio.Client
	Bucket     string
	URLExpires time.Duration
}

func NewMinioLFS(p MinioLFSInitParams) *MinioLFS {
	m := new(MinioLFS)
	api, err := minio.New(p.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(p.AccessKey, p.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		m.api = api
	}
	m.ctx = context.Background()
	m.Bucket = p.Bucket
	m.URLExpires = time.Duration(p.URLExpires) * time.Second
	return m
}

func (m *MinioLFS) IsExist(oid string) bool {
	if _, err := m.api.StatObject(m.ctx, m.Bucket, oid, minio.StatObjectOptions{}); err != nil {
		res := minio.ToErrorResponse(err)
		switch res.Code {
		case "NoSuchBucket":
		case "NoSuchKey":
			return false
		default:
			log.Fatal(err)
		}
	}
	return true
}

func (m *MinioLFS) DownloadURL(oid string) *url.URL {
	reqParams := make(url.Values)
	presignedURL, err := m.api.PresignedGetObject(m.ctx, m.Bucket, oid, m.URLExpires, reqParams)
	if err != nil {
		log.Fatal(err)
	}
	return presignedURL
}

func (m *MinioLFS) UploadURL(oid string) *url.URL {
	presignedURL, err := m.api.PresignedPutObject(m.ctx, m.Bucket, oid, m.URLExpires)
	if err != nil {
		log.Fatal(err)
	}
	return presignedURL
}
