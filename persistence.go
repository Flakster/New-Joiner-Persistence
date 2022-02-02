package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Employee struct {
	gorm.Model
	IdentificationNumber int
	Name string
	LastName string
	Stack string
	Role string
	EnglishLevel string
	DomainExperience string
}

type Task struct {
	gorm.Model
	ParentTaskId *Task
	Employee_Id *Employee
	Name string
	Description string
	EstimatedRequiredHours int
	Stack string
	MinRole []string
}

type Name struct {
	first_name string
	middle_name string
	last_name string
}

var db *gorm.DB
var err error

func DBConnect() {

	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5444 user=user1 dbname=NewJoinerTasks sslmode=disable password=pass123")
	if err != nil {
		panic("Failed to connect database")
	}
	defer db.Close()
}

func AddJoiner(w http.ResponseWriter, r *http.Request){
	//soon
}

func insertEmployee(employee []byte){
	var joinerMap interface{}
	json.Unmarshal(employee, &joinerMap)
	joiner := joinerMap.(map[string]interface{})
	skillsMap := joiner["Core Technical Skills"]
	skills := skillsMap.(map[string]interface{})
	var newJoiner Employee
	name := SplitName(fmt.Sprintf("%v", joiner["Name"]))
	newJoiner.Name = name.first_name + " " + name.middle_name
	newJoiner.LastName = name.last_name
	newJoiner.Role = fmt.Sprintf("%v", joiner["Role"])
	languages := fmt.Sprintf("%v", skills["Languages & Technologies"]) 
	databases := fmt.Sprintf("%v", skills["Databases"]) 
	newJoiner.Stack = languages + " "+ databases
	newJoiner.DomainExperience = fmt.Sprintf("%s", joiner["Domain Experience"])
	newJoiner.EnglishLevel = "Proficient"

	log.Printf("Joiner: %v",newJoiner)

	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5444 user=user1 dbname=NewJoinerTasks sslmode=disable password=pass123")
	if err != nil {
		panic("Failed to connect database")
	}
	defer db.Close()

	result:= db.Create(&newJoiner)
	CheckError(result.Error)
	log.Println("Record successfully created")
}

func SplitName(wholeName string)(name Name){
	words := strings.Fields(wholeName)
	name.first_name = words[0]
	if len(words) > 3 {
		name.middle_name= words[1]
		name.last_name = words[2] + " " + words[3]
		return
	}
	if len(words) > 2 {
		name.middle_name = words[2]
		name.last_name = words[3]
		return
	}
	if len(words) > 1 {
		name.last_name = words[len(words) -1]
		return
	}
	return
}

func CheckError(err error){
	if err !=nil {
		panic(err)
	}
}