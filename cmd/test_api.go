package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	baseURL := "http://localhost:8080/api"
	fmt.Println("=== Memulai Pengujian API Todo List ===")

	// 1. Health Check
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Health check gagal: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Health check sukses (200 OK)")

	// 2. Create Task (POST /api/tasks)
	newTask := map[string]interface{}{
		"title":       "Testing Task via Automated Script",
		"description": "Pengujian otomatis untuk API",
		"priority":    "high",
		"status":      "todo",
		"due_date":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	jsonBody, _ := json.Marshal(newTask)
	resp, err = http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil || resp.StatusCode != http.StatusCreated {
		fmt.Printf("❌ Create task gagal: %v, status: %v\n", err, resp.StatusCode)
		os.Exit(1)
	}
	var createdTask Task
	json.NewDecoder(resp.Body).Decode(&createdTask)
	resp.Body.Close()
	fmt.Printf("✅ Create task sukses (201 Created), ID: %s\n", createdTask.ID)

	taskID := createdTask.ID

	// 3. Get All Tasks (GET /api/tasks)
	resp, err = http.Get(baseURL + "/tasks")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Get all tasks gagal: %v\n", err)
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Println("✅ Get all tasks sukses (200 OK)")

	// 4. Get Task By ID (GET /api/tasks/{id})
	resp, err = http.Get(baseURL + "/tasks/" + taskID)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Get task by ID gagal: %v\n", err)
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Printf("✅ Get task by ID sukses (200 OK) untuk ID: %s\n", taskID)

	// 5. Update Task Status (PATCH /api/tasks/{id}/status)
	updateStatusBody, _ := json.Marshal(map[string]string{"status": "in_progress"})
	req, _ := http.NewRequest(http.MethodPatch, baseURL+"/tasks/"+taskID+"/status", bytes.NewBuffer(updateStatusBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Update task status gagal: %v, status: %v\n", err, resp.StatusCode)
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Println("✅ Update task status sukses (200 OK)")

	// 6. Update Task (PUT /api/tasks/{id})
	updateTaskBody, _ := json.Marshal(map[string]interface{}{
		"title":       "Updated Testing Task",
		"description": "Updated description",
		"priority":    "medium",
		"status":      "completed",
		"due_date":    time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req, _ = http.NewRequest(http.MethodPut, baseURL+"/tasks/"+taskID, bytes.NewBuffer(updateTaskBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Update task gagal: %v, status: %v\n", err, resp.StatusCode)
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Println("✅ Update task sukses (200 OK)")

	// 7. Delete Task (DELETE /api/tasks/{id})
	req, _ = http.NewRequest(http.MethodDelete, baseURL+"/tasks/"+taskID, nil)
	resp, err = client.Do(req)
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent) {
		fmt.Printf("❌ Delete task gagal: %v, status: %v\n", err, resp.StatusCode)
		os.Exit(1)
	}
	resp.Body.Close()
	fmt.Println("✅ Delete task sukses (204/200 No Content/OK)")

	fmt.Println("\n🎉 SEMUA PENGUJIAN API BERHASIL DIJALANKAN DENGAN SUKSES!")
}
