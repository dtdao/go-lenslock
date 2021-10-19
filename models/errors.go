package models

import "strings"

var (
	ErrorNotFound          modelError = "models: resource not found"
	ErrorIDInvalid         modelError = "models: ID provided is invalid"
	ErrorPasswordIncorrect modelError = "modals: Incorrect password provided"
	ErrorEmailRequired     modelError = "Email address is required"
	ErrorEmailInvalid      modelError = "Email address is not valid"
	ErrorEmailIsTaken      modelError = "models: email address is already taken"
	ErrorPasswordTooShort  modelError = "models: passwords must be at least 8 characters long"
	ErrorPasswordRequired  modelError = "models: passwords is required"
	ErrorRememberTooShort  modelError = "models: remember token must be at least 32 bytes"
	ErrorRememberRequired  modelError = "models: remember token is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}
