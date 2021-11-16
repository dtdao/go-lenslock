package models

import "strings"

var (
	ErrorNotFound          modelError = "models: resource not found"
	ErrorPasswordIncorrect modelError = "modals: Incorrect password provided"
	ErrorEmailRequired     modelError = "Email address is required"
	ErrorEmailInvalid      modelError = "Email address is not valid"
	ErrorEmailIsTaken      modelError = "models: email address is already taken"
	ErrorPasswordTooShort  modelError = "models: passwords must be at least 8 characters long"
	ErrorPasswordRequired  modelError = "models: passwords is required"
	ErrTitleRequired       modelError = "models: title is required"
	ErrorRememberTooShort  privateError = "models: remember token must be at least 32 bytes"
	ErrorRememberRequired  privateError = "models: remember token is required"
	ErrUserIDRequired      privateError = "models: User ID is required"
	ErrorIDInvalid         privateError = "models: ID provided is invalid"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}