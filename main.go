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

	if !found {
		fmt.Println("Created student doesn't exist")
		return
	}

	fmt.Println("Fetched student", fetchedStudent.Name)

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

	studentNameChanged := models.Student{
		ID:   createdStudent.ID,
		Name: "Bob",
	}

	updatedStudent, found, err := store.UpdateStudent(&studentNameChanged)
	if err != nil {
		fmt.Println("Error trying to update student")
		panic(err)
	}

	if !found {
		fmt.Println("Didn't find student when trying to update")
		return
	}

	fmt.Println("Updated student, new name is", updatedStudent.Name)

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

	found, err = store.DeleteStudentByID(updatedStudent.ID)
	if err != nil {
		fmt.Println("Error trying to delete student")
		panic(err)
	}

	if !found {
		fmt.Println("Didn't find student when trying to delete")
		return
	}

	fmt.Println("Deleted student successfully")

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
