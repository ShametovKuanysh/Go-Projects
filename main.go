package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Db struct {
	Tasks []Task `json:"tasks,omitempty"`
}

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func createTask(tasks *[]Task, description string) Task {
	return Task{
		Id:          len(*tasks) + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func updateTask(tasks *[]Task, id int, description string) {
	for i, task := range *tasks {
		if task.Id == id {
			(*tasks)[i].Description = description
			(*tasks)[i].UpdatedAt = time.Now()
		}
	}
}

func deleteTask(tasks *[]Task, id int) {
	for i, task := range *tasks {
		if task.Id == id {
			*tasks = append((*tasks)[:i], (*tasks)[i+1:]...)
			break
		}
	}
}

func listTasks(tasks []Task) {

	for _, task := range tasks {
		fmt.Printf("ID: %d, Description: %s, Status: %s, Created At: %s, Updated At: %s\n", task.Id, task.Description, task.Status, task.CreatedAt.Format(time.RFC1123), task.UpdatedAt.Format(time.RFC1123))
	}
}

func listTasksByStatus(tasks []Task, status string) {
	for _, task := range tasks {
		if task.Status == status {
			fmt.Printf("ID: %d, Description: %s, Status: %s, Created At: %s, Updated At: %s\n", task.Id, task.Description, task.Status, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339))
		}
	}
}

func updateDB(tasks []Task) {
	data, err := json.MarshalIndent(Db{Tasks: tasks}, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("db.json", data, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	var db Db

	b, err := os.ReadFile("db.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, &db)

	fmt.Println("-- Welcome to CLI ToDoList --------------------------------")

	fmt.Println("Current tasks:")
	listTasks(db.Tasks)

	fmt.Println("Enter command:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()

		parts := strings.Split(command, " ")
		switch parts[0] {
		case "add":
			db.Tasks = append(db.Tasks, createTask(&db.Tasks, strings.Join(parts[1:], " ")))
			fmt.Printf("-- Task added with ID: %d\n", db.Tasks[len(db.Tasks)-1].Id)
			updateDB(db.Tasks)
		case "update":
			id, _ := strconv.Atoi(parts[1])
			updateTask(&db.Tasks, id, strings.Join(parts[1:], ""))
			fmt.Printf("-- Task updated with ID: %d\n", id)
			updateDB(db.Tasks)
		case "delete":
			id, _ := strconv.Atoi(parts[1])
			deleteTask(&db.Tasks, id)
			fmt.Printf("-- Task deleted with ID: %d\n", id)
			updateDB(db.Tasks)
		case "list":
			if len(parts) > 1 {
				listTasksByStatus(db.Tasks, parts[1])
			} else {
				listTasks(db.Tasks)
			}
		}

		if parts[0] == "exit" {
			fmt.Println("Exiting...")
			break
		}
		fmt.Println("Enter command:")

	}

	if scanner.Err() != nil {
		return
	}
}
