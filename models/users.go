package models

import (
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var  (
	ErrorNotFound = errors.New("models: resource not found")
	ErrorInvalidId = errors.New("models: ID provided is invalid")
)

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

func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

func (us *UserService) AutoMigrate() error {
	err := us.db.Migrator().DropTable(&User{}).Error()
	if len(err) != 0 {
		return errors.New(err)
	}
	us.db.AutoMigrate(&User{})
	return nil
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
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

func (us *UserService) CreateUser(user *User) error{
	
	return us.db.Create(user).Error
}

func (us *UserService) Delete(id uint) error {
	if id == 0{
		return ErrorInvalidId
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}


func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorNotFound
	}
	return err
}

func (us *UserService) ById(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

func (us *UserService) Close() error {
	return us.Close()
}

type User struct {
	gorm.Model
	Name string
	Email string `gorm:"not null; unique_index"`
}