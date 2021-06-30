package models

import (
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrorNotFound = errors.New("models: resource not found")

type UserService struct {
	db *gorm.DB
}

func (us *UserService) Test() string{
	return "hello"
}

func (us *UserService) DestructiveReset()  {
	us.db.Migrator().DropTable(&User{})
	us.db.AutoMigrate(&User{})
}

func NewUserService(connectionInfo string) (*UserService, error){
	db, err := gorm.Open(postgres.Open(connectionInfo),&gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	} )
	if err != nil {
		panic(err)
	}

	return &UserService{
		db: db,
	}, nil
}


func (us *UserService) ById(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil: 
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrorNotFound
	default: 
		return nil, err
	}
}

func (us *UserService) Close() error {
	return us.Close()
}

type User struct {
	gorm.Model
	Name string
	Email string `gorm:"not null; unique_index"`
}