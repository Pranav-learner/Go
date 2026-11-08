package repository

import (
	"employee-management/employee"
	"errors"
)

var (
	ErrNotFound      = errors.New("employee not found")
	ErrAlreadyExists = errors.New("employee with this ID already exists")
)

// EmployeeRepository defines the contract for storage operations
type EmployeeRepository interface {
	Add(emp employee.Employee) error
	GetByID(id int) (employee.Employee, error)
	UpdateSalary(id int, newSalary float64) error
	Delete(id int) error
	GetAll() []employee.Employee
}