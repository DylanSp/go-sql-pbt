package main

import (
	"fmt"

	"github.com/DylanSp/go-sql-pbt/pkg/storage"
	"github.com/google/uuid"
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

	fetchedStudent, found, err := store.GetStudentByID(createdStudent.ID)
	if err != nil {
		fmt.Println("Unable to fetch student")
		panic(err)
	}

	if found {
		fmt.Println("Fetched student", fetchedStudent.Name)
	} else {
		fmt.Println("Created student doesn't exist")
		return
	}

	nonexistentStudent, found, err := store.GetStudentByID(uuid.New())
	if err != nil {
		fmt.Println("Error trying to fetch nonexistent student")
		panic(err)
	}

	if found {
		fmt.Println("Somehow fetched student with random ID")
		fmt.Println("Name:", nonexistentStudent.Name)
	} else {
		fmt.Println("Correctly failed to find nonexistent student")
	}

}
