package models

import (
	"fmt"
	"testing"
	"time"
)

func testingUserService() (*UserService, error){
	const (
		host = "localhost"
		port = 5432
		user = "dong"
		password = "password"
		dbname = "lenslocked_test"
	)	
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}

	us.DestructiveReset()

	return us, nil
}

func TestCreateUser(t *testing.T){
	us, err := testingUserService()
	if err !=nil {
		t.Fatal(err)
	}
	user := User{
		Name: "Michael Scott",
		Email:  "michael@dunermifflin.com",
	}
	err = us.CreateUser(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID==0{
		t.Errorf("Expected Id >0. Received %d", user.ID)
	}
	if time.Since(user.CreatedAt) > time.Duration(5*time.Second){
		t.Errorf("Expected create at to be recent. Received %s", user.CreatedAt)
	}
	if time.Since(user.UpdatedAt) > time.Duration(5*time.Second){
		t.Errorf("Expected updated at to be recent. Received %s", user.UpdatedAt)
	}
}