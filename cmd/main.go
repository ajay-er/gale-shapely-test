package main

import (
	"fmt"
	"gale-shapley/handler"
	"gale-shapley/models"
)




func main() {
	employees, err := handler.ReadEmployees("csv/employees.csv")
	if err != nil {
		fmt.Println("Error reading employees:", err)
		return
	}

	jobs, err :=handler.ReadJobs("csv/jobs.csv")
	if err != nil {
		fmt.Println("Error reading jobs:", err)
		return
	}

	matches :=handler.GaleShapley(employees, jobs)

	// Print matches
	for jobID, matches := range matches {
		fmt.Printf("Job %s:\n", jobID)
		for rank, match := range matches {
			var employee models.Employee
			for _, emp := range employees {
				if emp.ID == match.EmployeeID {
					employee = emp
					break
				}
			}
			fmt.Printf("%d. %s (Score: %.2f)\n", rank+1, employee.ID, match.Score)
		}
		fmt.Println()
	}
}
