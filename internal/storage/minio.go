// Package storage 提供 MinIO 对象存储功能
package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOConfig MinIO 配置
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// MinIOStorage MinIO 存储客户端
type MinIOStorage struct {
	client *minio.Client
	bucket string
}

// NewMinIOStorage 创建 MinIO 存储
func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("MinIO 客户端创建失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 确保 bucket 存在
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("检查 bucket 失败: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("创建 bucket 失败: %w", err)
		}
	}

	return &MinIOStorage{client: client, bucket: cfg.Bucket}, nil
}

// PresignedPutURL 生成预签名上传 URL
func (s *MinIOStorage) PresignedPutURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedPutObject(ctx, s.bucket, objectName, expires)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// PresignedGetURL 生成预签名下载 URL
func (s *MinIOStorage) PresignedGetURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// Upload 直接上传文件
func (s *MinIOStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Download 下载文件
func (s *MinIOStorage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
}

// Delete 删除文件
func (s *MinIOStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}

// Exists 检查文件是否存在
func (s *MinIOStorage) Exists(ctx context.Context, objectName string) bool {
	_, err := s.client.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
	return err == nil
}

// FileInfo 文件信息
type FileInfo struct {
	Name         string
	Size         int64
	ContentType  string
	LastModified time.Time
}

// Stat 获取文件信息
func (s *MinIOStorage) Stat(ctx context.Context, objectName string) (*FileInfo, error) {
	info, err := s.client.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Name:         info.Key,
		Size:         info.Size,
		ContentType:  info.ContentType,
		LastModified: info.LastModified,
	}, nil
}

// ListFiles 列出文件
func (s *MinIOStorage) ListFiles(ctx context.Context, prefix string, maxKeys int) ([]FileInfo, error) {
	var files []FileInfo
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	count := 0
	for obj := range s.client.ListObjects(ctx, s.bucket, opts) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		files = append(files, FileInfo{
			Name:         obj.Key,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
			LastModified: obj.LastModified,
		})
		count++
		if maxKeys > 0 && count >= maxKeys {
			break
		}
	}
	return files, nil
}
