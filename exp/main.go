package main

import (
	"fmt"

	"lenslocked.com/models"
)

const (
	host = "localhost"
	port = 5432
	user = "dong"
	password = "password"
	dbname = "lenslocked_dev"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
 host, port, user, password, dbname)


	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

	// us.DestructiveReset()
	// user := models.User {
	// 	Name: "Michael Scott",
	// 	Email:  "michael@dunermifflin.com",
	// }
	// if err := us.CreateUser(&user); err != nil {
	// 	panic(err)
	// }
	user, err := us.ById(1)
 
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

	// t, err := template.ParseFiles("hello.gohtml")
	// if err != nil {
	// 	panic(err)
	// }

	// data := struct {
	// 	Name   string
	// 	Place  string
	// 	Time   int
	// 	Nested struct {
	// 		Name  string
	// 		Level int
	// 	}
	// }{Name: "John Smith", Place: "Tokyo", Nested: struct {
	// 	Name  string
	// 	Level int
	// }{"TEST", 3}}

	// err = t.Execute(os.Stdout, data)
	// if err != nil {
	// 	panic(err)
	// }
}
