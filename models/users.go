package models

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"lenslocked.com/hash"
	"lenslocked.com/rand"
)

var (
	ErrorNotFound        = errors.New("models: resource not found")
	ErrorInvalidId       = errors.New("models: ID provided is invalid")
	ErrorInvalidPassword = errors.New("modals: Incorrect password provided")
	//ErrorInvalidEmail = errors.New("models: Incorerect email provided")
)

const userPwPepper = "randomPepperForThePizza"
const hmacSecretKey = "secret-hmac-key"

type UserDB interface {
	// methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)
	ByAge(age uint8) (*User, error)
	InAgeRange(min uint8, max uint8) ([]User, error)

	// methods for altering users
	CreateUser(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// used to close a db connection
	Close() error
}

// UserService is a set of methods used to manipulate
// and work with the user modal.
type UserService interface {
	// Authenticate will verify the provided email address and password
	// are correct.  If they are correct, the user
	// corresponding to that email will be returned.  Otherwise
	// you will received either:
	// error not found , error invalid, or another error if something
	// goes wrong
	Authenticate(email, password string) (*User, error)
	UserDB
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	hmac := hash.NewHMAC(hmacSecretKey)
	if err != nil {
		panic(err)
	}

	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

type userService struct {
	UserDB
}

type userValidator struct {
	UserDB
}

func (ug *userGorm) Test() string {
	return "hello"
}

func (ug *userGorm) DestructiveReset() {
	ug.db.Migrator().DropTable(&User{})
	ug.db.AutoMigrate(&User{})
}

func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

func (ug *userGorm) AutoMigrate() error {
	ug.db.Migrator().DropTable(&User{})
	// if len(err) != 0 {
	// 	return errors.New(err)
	// }
	ug.db.AutoMigrate(&User{})
	return nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with a given remember token
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	hashedToken := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", hashedToken), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))

	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrorInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

func (ug *userGorm) ByAge(age uint8) (*User, error) {
	var user User
	db := ug.db.Where("age = ?", age)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) InAgeRange(min uint8, max uint8) ([]User, error) {
	var users []User
	db := ug.db.Where("age >= ? AND age <= ?", min, max).Find(&users)
	err := all(db, &users)
	return users, err
}

func NewUserService(connectionInfo string) (*userService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

func (ug *userGorm) CreateUser(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)
	//if user.Remember != "" {
	//	user.RememberHash = us.hmac.Hash(user.Remember)
	//} else {
	//	token, err := rand.RememberToken()
	//	user.RememberHash = us.hmac.Hash(token)
	//}
	return ug.db.Create(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrorInvalidId
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorNotFound
	}
	return err
}

func all(db *gorm.DB, dst interface{}) error {
	err := db.Find(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorNotFound
	}
	return err
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) Close() error {
	return ug.Close()
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null; unique_index"`
	Age          uint8
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"size:60; not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}
