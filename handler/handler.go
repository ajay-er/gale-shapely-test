package handler

import (
	"encoding/csv"
	"gale-shapley/models"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ReadEmployees(filename string) ([]models.Employee, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var employees []models.Employee
	for _, record := range records[1:] {
		age, _ := strconv.Atoi(record[3])
		experience, _ := strconv.Atoi(record[4])
		expectedSalary, _ := strconv.Atoi(record[5])
		employee := models.Employee{
			ID:             record[0],
			Skills:         strings.Split(record[1], ","),
			Location:       record[2],
			Age:            age,
			Experience:     experience,
			ExpectedSalary: expectedSalary,
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func ReadJobs(filename string) ([]models.Job, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var jobs []models.Job
	for _, record := range records[1:] {
		experience, _ := strconv.Atoi(record[3])
		budget, _ := strconv.Atoi(record[4])
		vacancies, _ := strconv.Atoi(record[5])
		job := models.Job{
			ID:         record[0],
			Skills:     strings.Split(record[1], ","),
			Location:   record[2],
			Experience: experience,
			Budget:     budget,
			Vacancies:  vacancies,
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GaleShapley performs the Gale-Shapley algorithm to match employees to jobs
func GaleShapley(employees []models.Employee, jobs []models.Job) map[string][]models.Match {
	employeePreferences := make(map[string][]models.Match)
	jobPreferences := make(map[string][]models.Match)

	// Create preference lists based on scores
	for _, employee := range employees {
		var preferences []models.Match
		for _, job := range jobs {
			preferences = append(preferences, models.Match{
				EmployeeID: employee.ID,
				Score:      CalculateScore(employee, job),
				Skills: employee.Skills,
			})
		}
		sort.Slice(preferences, func(i, j int) bool {
			return preferences[i].Score > preferences[j].Score // Sort by score descending
		})
		employeePreferences[employee.ID] = preferences
	}

	for _, job := range jobs {
		var preferences []models.Match
		for _, employee := range employees {
			preferences = append(preferences, models.Match{
				EmployeeID: employee.ID,
				Score:      CalculateScore(employee, job),
			})
		}
		sort.Slice(preferences, func(i, j int) bool {
			return preferences[i].Score > preferences[j].Score // Sort by score descending
		})
		jobPreferences[job.ID] = preferences
	}

	jobMatches := make(map[string][]models.Match)

	// Populate jobMatches based on employee preferences
	for _, job := range jobs {
		for _, preference := range jobPreferences[job.ID] {
			employeeID := preference.EmployeeID
			jobMatches[job.ID] = append(jobMatches[job.ID], models.Match{
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

// CalculateScore calculates the score for an employee-job pair
func CalculateScore(employee models.Employee, job models.Job) float64 {
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
