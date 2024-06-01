package main

import (
	"fmt"

	"github.com/DylanSp/go-sql-pbt/pkg/storage"
)

func main() {
	dbCfg := storage.DBConfig{
		Host:   "localhost",
		Port:   "5432",
		DBName: "school",

		Username: "postgres",
		Password: "devpassword",
	}

	store, err := storage.NewStore(dbCfg)
	if err != nil {
		fmt.Println("Unable to initialize store")
		panic(err)
	}

	studentName := "Alice"
	createdStudent, err := store.CreateStudent(studentName)
	if err != nil {
		fmt.Println("Unable to create student")
		panic(err)
	}

	fetchedStudent, err := store.GetStudentByID(createdStudent.ID)
	if err != nil {
		fmt.Println("Unable to fetch student")
		panic(err)
	}

	fmt.Println("Fetched student", fetchedStudent.Name)
}
