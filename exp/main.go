package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	data := struct {
		Name   string
		Place  string
		Time   int
		Nested struct {
			Name  string
			Level int
		}
	}{Name: "John Smith", Place: "Tokyo", Nested: struct {
		Name  string
		Level int
	}{"TEST", 3}}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
