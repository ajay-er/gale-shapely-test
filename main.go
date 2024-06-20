package main

import (
	"fmt"
	"sort"
)

// Employee struct represents an employee
type Employee struct {
	ID             string
	Skills         []string
	Location       string
	Age            int
	Experience     int
	ExpectedSalary int
}

// Job struct represents a job
type Job struct {
	ID         string
	Skills     []string
	Location   string
	Experience int
	Budget     int
	Vacancies  int
}

// Match struct represents a match between an employee and a job
type Match struct {
	EmployeeID string
	Score      float64
}

// CalculateScore calculates the score for an employee-job pair
func CalculateScore(employee Employee, job Job) float64 {
	score := 0.0

	// Skills match
	skillsMatch := float64(len(intersect(employee.Skills, job.Skills))) / float64(len(job.Skills))
	score += skillsMatch

	// Location match
	locationMatch := 0.0
	if employee.Location == job.Location {
		locationMatch = 1.0
	}
	score += locationMatch

	// Experience match
	experienceMatch := float64(employee.Experience) / float64(job.Experience)
	if experienceMatch > 1.0 {
		experienceMatch = 1.0
	}
	score += experienceMatch

	// Salary match (inverse since we prefer lower expected salary)
	salaryMatch := 1.0 - float64(employee.ExpectedSalary)/float64(job.Budget)
	if salaryMatch < 0.0 {
		salaryMatch = 0.0
	} else if salaryMatch > 1.0 {
		salaryMatch = 1.0
	}
	score += salaryMatch

	return score
}

// GaleShapley performs the Gale-Shapley algorithm to match employees to jobs
func GaleShapley(employees []Employee, jobs []Job) map[string][]Match {
	employeePreferences := make(map[string][]Match)
	jobPreferences := make(map[string][]Match)

	// Create preference lists based on scores
	for _, employee := range employees {
		var preferences []Match
		for _, job := range jobs {
			preferences = append(preferences, Match{
				EmployeeID: employee.ID,
				Score:      CalculateScore(employee, job),
			})
		}
		sort.Slice(preferences, func(i, j int) bool {
			return preferences[i].Score > preferences[j].Score // Sort by score descending
		})
		employeePreferences[employee.ID] = preferences
	}

	for _, job := range jobs {
		var preferences []Match
		for _, employee := range employees {
			preferences = append(preferences, Match{
				EmployeeID: employee.ID,
				Score:      CalculateScore(employee, job),
			})
		}
		sort.Slice(preferences, func(i, j int) bool {
			return preferences[i].Score > preferences[j].Score // Sort by score descending
		})
		jobPreferences[job.ID] = preferences
	}

	jobMatches := make(map[string][]Match)

	// Populate jobMatches based on employee preferences
	for _, job := range jobs {
		for _, preference := range jobPreferences[job.ID] {
			employeeID := preference.EmployeeID
			jobMatches[job.ID] = append(jobMatches[job.ID], Match{
				EmployeeID: employeeID,
				Score:      preference.Score,
			})
		}
	}

	// Sort employees within each job by score descending
	for _, matches := range jobMatches {
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score > matches[j].Score
		})
	}

	return jobMatches
}

// Function to find intersection of two slices
func intersect(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	var intersection []string

	for _, item := range slice1 {
		m[item] = true
	}

	for _, item := range slice2 {
		if m[item] {
			intersection = append(intersection, item)
		}
	}

	return intersection
}

func main() {
	// Original datasets
	employees := []Employee{
		{ID: "Alice", Skills: []string{"JavaScript", "React"}, Location: "New York", Age: 30, Experience: 5, ExpectedSalary: 60000},
		{ID: "Bob", Skills: []string{"Python", "Django"}, Location: "San Francisco", Age: 28, Experience: 3, ExpectedSalary: 70000},
		{ID: "Charlie", Skills: []string{"Java", "Spring"}, Location: "New York", Age: 35, Experience: 10, ExpectedSalary: 90000},
		{ID: "David", Skills: []string{"JavaScript", "Node.js"}, Location: "Los Angeles", Age: 26, Experience: 2, ExpectedSalary: 55000},
		{ID: "Eve", Skills: []string{"Ruby", "Rails"}, Location: "San Francisco", Age: 29, Experience: 4, ExpectedSalary: 65000},
		{ID: "Frank", Skills: []string{"C#", ".NET"}, Location: "Texas", Age: 32, Experience: 6, ExpectedSalary: 80000},
		{ID: "Grace", Skills: []string{"PHP", "Laravel"}, Location: "Florida", Age: 27, Experience: 3, ExpectedSalary: 60000},
		{ID: "Henry", Skills: []string{"Python", "Flask"}, Location: "New York", Age: 34, Experience: 8, ExpectedSalary: 85000},
		{ID: "Ivy", Skills: []string{"Java", "Spring"}, Location: "San Francisco", Age: 31, Experience: 7, ExpectedSalary: 90000},
		{ID: "Jack", Skills: []string{"JavaScript", "Angular"}, Location: "Texas", Age: 29, Experience: 5, ExpectedSalary: 70000},
	}

	jobs := []Job{
		{ID: "Google - Web Developer", Skills: []string{"JavaScript", "React"}, Location: "New York", Experience: 4, Budget: 65000, Vacancies: 2},
		{ID: "Facebook - Python Developer", Skills: []string{"Python", "Django"}, Location: "San Francisco", Experience: 3, Budget: 75000, Vacancies: 1},
		{ID: "Amazon - Java Developer", Skills: []string{"Java", "Spring"}, Location: "New York", Experience: 7, Budget: 85000, Vacancies: 1},
		{ID: "Netflix - Node.js Developer", Skills: []string{"JavaScript", "Node.js"}, Location: "Los Angeles", Experience: 2, Budget: 60000, Vacancies: 1},
		{ID: "Twitter - Ruby on Rails Developer", Skills: []string{"Ruby", "Rails"}, Location: "San Francisco", Experience: 4, Budget: 70000, Vacancies: 1},
		{ID: "Microsoft - C#/.NET Developer", Skills: []string{"C#", ".NET"}, Location: "Texas", Experience: 5, Budget: 80000, Vacancies: 1},
		{ID: "Uber - PHP/Laravel Developer", Skills: []string{"PHP", "Laravel"}, Location: "Florida", Experience: 3, Budget: 65000, Vacancies: 1},
		{ID: "Airbnb - Python/Flask Developer", Skills: []string{"Python", "Flask"}, Location: "New York", Experience: 7, Budget: 90000, Vacancies: 1},
		{ID: "Salesforce - Java/Spring Developer", Skills: []string{"Java", "Spring"}, Location: "San Francisco", Experience: 8, Budget: 95000, Vacancies: 1},
		{ID: "Tesla - JavaScript/Angular Developer", Skills: []string{"JavaScript", "Angular"}, Location: "Texas", Experience: 5, Budget: 75000, Vacancies: 1},
	} 

	matches := GaleShapley(employees, jobs)

	// Print matches
	for jobID, matches := range matches {
		fmt.Printf("Job %s:\n", jobID)
		for rank, match := range matches {
			var employee Employee
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
