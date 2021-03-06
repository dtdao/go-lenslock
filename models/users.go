package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"lenslocked.com/hash"
	"lenslocked.com/rand"
	"regexp"
	"strings"
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

//func newUserGorm(connectionInfo string) (*userGorm, error) {
//	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info),
//	})
//	if err != nil {
//		panic(err)
//	}
//	return &userGorm{
//		db: db,
//	}, nil
//}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

type userService struct {
	UserDB
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	err := runUserValFuncs(&user, uv.hmacRemember)
	if err != nil {
		return nil, err
	}
	//hashedToken := uv.hmac.Hash(token)
	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) CreateUser(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.defaultRemember,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail,
	)
	if err != nil {
		return err
	}
	//user.RememberHash = uv.hmac.Hash(user.Remember)
	//if user.Remember != "" {
	//	user.RememberHash = us.hmac.Hash(user.Remember)
	//} else {
	//	token, err := rand.RememberToken()
	//	user.RememberHash = us.hmac.Hash(token)
	//}
	return uv.UserDB.CreateUser(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrorRememberRequired
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// setRememberIfUnset
func (uv *userValidator) defaultRemember(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

//func (uv *userValidator) idGreaterThanZero(user *User) error {
//	if user.ID <= 0 {
//		return ErrorIDInvalid
//	}
//	return nil
//}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrorRememberTooShort
	}
	return nil
}

func (uv *userValidator) idGreaterThan(u uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= u {
			return ErrorIDInvalid
		}
		return nil
	})
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrorPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrorPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrorPasswordRequired
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrorEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrorEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrorNotFound {
		// Email is not taken
		return nil
	}
	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		return ErrorEmailIsTaken
	}
	return nil
}

func (ug *userGorm) Test() string {
	return "hello"
}

func (ug *userGorm) DestructiveReset() {
	ug.db.Migrator().DropTable(&User{})
	ug.db.AutoMigrate(&User{})
}

func (ug *userGorm) Update(user *User) error {
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

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(email)
}

// ByRemember looks up a user with a given remember token
func (ug *userGorm) ByRemember(remberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", remberHash), &user)
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
			return nil, ErrorPasswordIncorrect
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

func NewUserService(db *gorm.DB) UserService  {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) CreateUser(user *User) error {
	return ug.db.Create(user).Error
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
