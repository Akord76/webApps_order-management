package model

import "time"

type Employee struct {
	EmployeeID         int        `json:"employee_id"`
	EmployeeCardNumber int        `json:"employee_card_number"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	DateOfBirth        *time.Time `json:"date_of_birth"`
	Gender             string     `json:"gender"`
	Email              string     `json:"email"`
	DepartmentID       int        `json:"department_id"`
	PhoneNumber        string     `json:"phone_number"`
	BankNumber         string     `json:"bank_number"`
	PhotoPath          string     `json:"photo_path"`
}
