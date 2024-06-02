package main

import (
	"fmt"

	"github.com/DylanSp/go-sql-pbt/pkg/models"
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

	nonexistentStudent, found, err = store.UpdateStudent(&models.Student{
		ID:   uuid.New(),
		Name: "Michael McDoesntExist",
	})
	if err != nil {
		fmt.Println("Error trying to update nonexistent student")
		panic(err)
	}

	if found {
		fmt.Println("Somehow updated student with random ID")
		fmt.Println("Name:", nonexistentStudent.Name)
	} else {
		fmt.Println("Correctly failed to find and update nonexistent student")
	}

	found, err = store.DeleteStudentByID(uuid.New())
	if err != nil {
		fmt.Println("Error trying to delete nonexistent student")
		panic(err)
	}

	if found {
		fmt.Println("Somehow deleted student with random ID")
	} else {
		fmt.Println("Correctly failed to find nonexistent student when trying to delete")
	}
}
