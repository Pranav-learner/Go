package service

import (
	"errors"
	"fmt"
	"strings"
	"employee-management/employee"
	"employee-management/repository"
)
// Business validation sentinel errors
var (
	ErrInvalidID     = errors.New("employee ID must be greater than zero")
	ErrEmptyName     = errors.New("employee name cannot be empty")
	ErrEmptyDept     = errors.New("department cannot be empty")
	ErrInvalidSalary = errors.New("salary must be greater than zero")
)

type EmployeeService struct{
	repo repository.EmployeeRepository
}

// NewEmployeeService creates a new service instance
func NewEmployeeService(repo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) CreateEmployee(id int, name, dept string, salary float64) error {
	if id <= 0 {
		return ErrInvalidID
	}

	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}

	if strings.TrimSpace(dept) == "" {
		return ErrEmptyDept
	}

	if salary <= 0 {
		return ErrInvalidSalary
	} 

	emp := employee.Employee{
		ID:         id,
		Name:       strings.TrimSpace(name),
		Department: strings.TrimSpace(dept),
		Salary:     salary,
	}

	if err := s.repo.Add(emp); err != nil {
		// Uses %w to wrap repository error while adding context
		return fmt.Errorf("failed to create employee: %w", err)	
	}
	return nil
}

// GetEmployee retrieves an employee by ID
func (s *EmployeeService) GetEmployee(id int) (employee.Employee, error) {
	if id <= 0 {
		return employee.Employee{}, ErrInvalidID
	}
	emp, err := s.repo.GetByID(id)
	if err != nil {
		return employee.Employee{}, fmt.Errorf("lookup failed for ID %d: %w", id, err)
	}
	return emp, nil
}

// UpdateSalary validates the salary and updates through repository
func (s *EmployeeService) UpdateSalary(id int, newSalary float64) error {
	if id <= 0 {
		return ErrInvalidID
	}
	if newSalary <= 0 {
		return ErrInvalidSalary
	}
	if err := s.repo.UpdateSalary(id, newSalary); err != nil {
		return fmt.Errorf("failed to update salary for ID %d: %w", id, err)
	}
	return nil
}


// DeleteEmployee validates ID and deletes through repository
func (s *EmployeeService) DeleteEmployee(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete employee ID %d: %w", id, err)
	}
	return nil
}
// ListEmployees returns all employees
func (s *EmployeeService) ListEmployees() []employee.Employee {
	return s.repo.GetAll()
}