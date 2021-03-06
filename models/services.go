package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.Config.Logger.LogMode(logger.Info)
	return &Services{
		User: NewUserService(db),
		GalleryService: NewGalleryService(db),
	}, nil
}

type Services struct {
	GalleryService GalleryService
	User           UserService
}
