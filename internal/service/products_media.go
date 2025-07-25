package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/lavatee/dresscode_backend/internal/model"
	"github.com/lavatee/dresscode_backend/internal/repository"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type ProductsMediaService struct {
	repo   *repository.Repository
	s3     *minio.Client
	bucket string
}

func NewProductsMediaService(repo *repository.Repository, s3 *minio.Client, bucket string) *ProductsMediaService {
	return &ProductsMediaService{
		repo:   repo,
		s3:     s3,
		bucket: bucket,
	}
}

func (s *ProductsMediaService) getMediaURL(key string) string {
	return fmt.Sprintf("https://5a1bc5f7-b5c2-4a61-969a-beacbd4d7999.selstorage.ru/%s", key)
}

func (s *ProductsMediaService) UploadOneProductMedia(ctx context.Context, userId int, productID int, media model.ProductMedia, fileName string, isProductMain bool, file multipart.File) (int, string, error) {
	if !s.repo.Auth.IsAdmin(userId) {
		return 0, "", errors.New("user is not admin")
	}
	media.ProductID = productID
	media.URL = s.getMediaURL(fileName)
	id, err := s.repo.ProductsMedia.CreateOneProductMedia(media, isProductMain)
	if err != nil {
		return 0, "", err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return 0, "", err
	}
	info, err := s.s3.PutObject(ctx, s.bucket, fileName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return 0, "", err
	}
	logrus.Infof("URL: %s, KEY: %s", media.URL, info.Key)
	return id, media.URL, nil
}

func (s *ProductsMediaService) GetProductMedia(productID int) ([]model.ProductMedia, error) {
	return s.repo.ProductsMedia.GetProductMedia(productID)
}

func (s *ProductsMediaService) DeleteOneProductMedia(userId int, mediaId int) error {
	if !s.repo.Auth.IsAdmin(userId) {
		return errors.New("user is not admin")
	}
	return s.repo.ProductsMedia.DeleteOneProductMedia(mediaId)
}
