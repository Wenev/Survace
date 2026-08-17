package minio

import (
	"github.com/Acad600-TPA/WEB-WE-251/galactus/config"
	"os"
)

func MinIOConnect() (*config.MinIOClient, error) {
	user := os.Getenv("MINIO_ROOT_USER")
	password := os.Getenv("MINIO_ROOT_PASSWORD")
	minioUrl := os.Getenv("MINIO_URL")
	minio, err := config.NewMinIOClient(minioUrl, user, password, "video")
	if err != nil {
		return nil, err
	}
	return minio, nil
}
