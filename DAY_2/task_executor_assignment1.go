package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents an individual unit of work
type Task struct {
	ID       int
	Name     string
	Duration time.Duration // Simulated work time
}

// TaskResult holds the output from a finished task
type TaskResult struct {
	TaskID    int
	TaskName  string
	Duration  time.Duration
	Output    string
}

// executeTask runs on its own goroutine
func executeTask(task Task, results chan<- TaskResult, wg *sync.WaitGroup) {
	// Rule: Always defer wg.Done() so it is guaranteed to execute
	defer wg.Done()

	fmt.Printf("[Started]  Task #%d: %s (will take %v)\n", task.ID, task.Name, task.Duration)

	// Simulate work using the task's specific duration
	time.Sleep(task.Duration)

	// Send the result into the channel
	results <- TaskResult{
		TaskID:   task.ID,
		TaskName: task.Name,
		Duration: task.Duration,
		Output:   fmt.Sprintf("Completed %s successfully", task.Name),
	}

	fmt.Printf("[Finished] Task #%d: %s\n", task.ID, task.Name)
}

// executeTask runs on its own goroutine
func executeTask(task Task, results chan<- TaskResult, wg *sync.WaitGroup) {
	// Rule: Always defer wg.Done() so it is guaranteed to execute
	defer wg.Done()

	fmt.Printf("[Started]  Task #%d: %s (will take %v)\n", task.ID, task.Name, task.Duration)

	// Simulate work using the task's specific duration
	time.Sleep(task.Duration)

	// Send the result into the channel
	results <- TaskResult{
		TaskID:   task.ID,
		TaskName: task.Name,
		Duration: task.Duration,
		Output:   fmt.Sprintf("Completed %s successfully", task.Name),
	}

	fmt.Printf("[Finished] Task #%d: %s\n", task.ID, task.Name)
}

func main() {
	startTime := time.Now()

	fmt.Println("==================================================")
	fmt.Println("       CONCURRENT TASK EXECUTOR (10 TASKS)        ")
	fmt.Println("==================================================")

	// 1. Define 10 different tasks with varying durations
	tasks := []Task{
		{ID: 1, Name: "Fetch User Profile", Duration: 200 * time.Millisecond},
		{ID: 2, Name: "Download Avatar", Duration: 350 * time.Millisecond},
		{ID: 3, Name: "Query Database", Duration: 150 * time.Millisecond},
		{ID: 4, Name: "Generate PDF Report", Duration: 500 * time.Millisecond},
		{ID: 5, Name: "Send Welcome Email", Duration: 300 * time.Millisecond},
		{ID: 6, Name: "Compress Log Files", Duration: 250 * time.Millisecond},
		{ID: 7, Name: "Calculate Analytics", Duration: 400 * time.Millisecond},
		{ID: 8, Name: "Call Payment Gateway", Duration: 180 * time.Millisecond},
		{ID: 9, Name: "Sync Search Index", Duration: 220 * time.Millisecond},
		{ID: 10, Name: "Push Notification", Duration: 120 * time.Millisecond},
	}

	// 2. Create a buffered channel for all 10 results
	results := make(chan TaskResult, len(tasks))

	var wg sync.WaitGroup

	// 3. Launch 10 goroutines concurrently
	fmt.Println("\n--- Launching 10 Goroutines ---")
	for _, task := range tasks {
		wg.Add(1) // Call Add(1) BEFORE launching the goroutine!
		go executeTask(task, results, &wg)
	}

	// 4. Monitor goroutine: Wait for all tasks to finish, then close the channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 5. Collect and print all results as they arrive
	fmt.Println("\n--- Collecting Results ---")
	for res := range results {
		fmt.Printf("   [Result] Task #%-2d | %-22s | Time: %v\n", res.TaskID, res.TaskName, res.Duration)
	}

	totalTime := time.Since(startTime)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("All 10 tasks completed in: %v\n", totalTime)
	fmt.Println("(Notice: Sequentially this would take ~2.67s, but concurrently only ~500ms!)")
	fmt.Println("==================================================")
}