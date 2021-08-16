package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"lenslocked.com/hash"
)

const (
	host = "localhost"
	port = 5432
	user = "dong"
	password = "password"
	dbname = "lenslocked_dev"
)

func main(){
	//fmt.Println(rand.String(10))
	//fmt.Println(rand.RememberToken())
	toHash := []byte("this is my string to has")
	h := hmac.New(sha256.New, []byte("my-secret-key"))

	h.Write(toHash)

	b := h.Sum(nil)

	fmt.Println(base64.URLEncoding.EncodeToString(b))

	hp := hash.NewHMAC("my-secret-key")
	fmt.Println(hp)
}

//func main() {
//
//	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// host, port, user, password, dbname)
//
//
//	us, err := models.NewUserService(psqlInfo)
//	if err != nil {
//		panic(err)
//	}
//
//	us.DestructiveReset()
//	user := models.User {
//		Name: "Michael Scott",
//		Email:  "michael@dunermifflin.com",
//	}
//	if err := us.CreateUser(&user); err != nil {
//		panic(err)
//	}
//
//	user.Email = "michael@michaelscottpaperco.com"
//	if err := us.Update(&user); err != nil {
//		panic(err)
//	}
//
//	// userByEmail, err := us.ByEmail("michael@michaelscottpaperco.com")
//	// userById, err := us.ById(user.ID)
//
//	if err := us.Delete(user.ID); err != nil {
//		panic(err)
//	}
//
//	userById, err := us.ById(user.ID)
//	if err != nil {
//		panic(err)
//	}
//
//
//
//	// fmt.Println(userByEmail)
//	fmt.Println(userById)
//
//	// t, err := template.ParseFiles("hello.gohtml")
//	// if err != nil {
//	// 	panic(err)
//	// }
//
//	// data := struct {
//	// 	Name   string
//	// 	Place  string
//	// 	Time   int
//	// 	Nested struct {
//	// 		Name  string
//	// 		Level int
//	// 	}
//	// }{Name: "John Smith", Place: "Tokyo", Nested: struct {
//	// 	Name  string
//	// 	Level int
//	// }{"TEST", 3}}
//
//	// err = t.Execute(os.Stdout, data)
//	// if err != nil {
//	// 	panic(err)
//	// }
//}
