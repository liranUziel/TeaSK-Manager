package task

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Task struct {
    Id      int
    Title   string
    Status  bool
    DueDate string
}

const (
    FileName  = "tasks.gt"
    FolderName = "TeaSK"
)

func LoadTasks() ([]Task, error) {
    documentsPath := os.Getenv("USERPROFILE") + "\\Documents" // For Windows
    if documentsPath == "" {
        documentsPath = os.Getenv("HOME") // For Unix-based systems
    }

    taskFolder := fmt.Sprintf("%s\\%s", documentsPath, FolderName)
    taskFilePath := fmt.Sprintf("%s\\%s", taskFolder, FileName)

    if _, err := os.Stat(taskFolder); os.IsNotExist(err) {
        err := os.MkdirAll(taskFolder, os.ModePerm)
        if err != nil {
            return nil, err
        }
    }

    if _, err := os.Stat(taskFilePath); os.IsNotExist(err) {
        // create file if not exists
        _, err := os.Create(taskFilePath)
        if err != nil {
            return nil, err
        }
    }

    data, err := os.ReadFile(taskFilePath)
    if err != nil {
        return nil, err
    }

    taskStrings := strings.Split(string(data), ";")
    var tasks []Task

    for _, taskString := range taskStrings {
        taskParts := strings.Split(taskString, ":")
        if len(taskParts) != 4 {
            continue
        }
        id := 0
        fmt.Sscanf(taskParts[0], "%d", &id)
        status := false
        if taskParts[2] == "Done" {
            status = true
        }

        task := Task{
            Id:      id,
            Title:   taskParts[1],
            Status:  status,
            DueDate: taskParts[3],
        }
        tasks = append(tasks, task)
    }
    return tasks, nil
}

func SaveTasks(tasks []Task) error {
    // Create or open the log file
    logFile, err := os.OpenFile("task_manager.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer logFile.Close()

    // Set log output to the log file
    log.SetOutput(logFile)

    documentsPath := os.Getenv("USERPROFILE") + "\\Documents"
    if documentsPath == "" {
        documentsPath = os.Getenv("HOME") + "\\Documents"
    }

    taskFolder := fmt.Sprintf("%s\\%s", documentsPath, FolderName)
    taskFilePath := fmt.Sprintf("%s\\%s", taskFolder, FileName)

    if _, err := os.Stat(taskFolder); os.IsNotExist(err) {
        err := os.MkdirAll(taskFolder, os.ModePerm)
        if err != nil {
            log.Println("Error creating task folder:", err)
            return err
        }
    }

    file, err := os.OpenFile(taskFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        log.Println("Error opening task file:", err)
        return err
    }
    defer file.Close()

    if len(tasks) == 0 {
        log.Println("No tasks to save")
    } else {
        log.Printf("Saving %d tasks\n", len(tasks))
    }

    var sb strings.Builder
    for _, task := range tasks {
        status := "In-Progress"
        if task.Status {
            status = "Done"
        }
        sb.WriteString(fmt.Sprintf("%d:%s:%s:%s;", task.Id, task.Title, status, task.DueDate))
    }

    data := sb.String()
    log.Println("Saving tasks:", data) // Log the data being saved

    _, err = file.WriteString(data)
    if err != nil {
        log.Println("Error writing to file:", err) // Log any write errors
    }
    return err
}