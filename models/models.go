package models

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
	Skills     []string
}
