package repository

import (
	"employee-management/employee"
)

// InMemoryRepository stores employees in memory using a map.
type InMemoryRepository struct {
	employees map[int]employee.Employee
}

// NewInMemoryRepository is a constructor that initializes the map.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		employees: make(map[int]employee.Employee),
	}
}

// Add inserts a new employee.
// Returns ErrAlreadyExists if the employee ID already exists.
func (r *InMemoryRepository) Add(emp employee.Employee)error{
	if _,exists := r.employees[emp.ID];exists{
		return ErrAlreadyExists
	}
	r.employees[emp.ID] = emp
	return nil
}

// GetByID looks up an employee by ID.
func (r *InMemoryRepository) GetByID(id int) (employee.Employee, error){
	emp,exists := r.employees[id]
	if !exists{
		return employee.Employee{}, ErrNotFound
	}
	return emp, nil
}

// GetAll returns a slice of all employees.
func (r *InMemoryRepository) GetAll() []employee.Employee{
	emps := make([]employee.Employee, 0, len(r.employees))
	for _,emp := range r.employees{
		emps = append(emps, emp)
	}
	return emps
}

// UpdateSalary upates an existing employee's salary

func (r *InMemoryRepository) UpdateSalary(id int, newSalary float64) error {
	emp,exists := r.employees[id]
	if !exists {
		return ErrNotFound
	}
	emp.Salary = newSalary
	r.employees[id] = emp
	return nil
}

// Delete removes an employee by ID.
func (r *InMemoryRepository) Delete(id int) error {
	if _, exists := r.employees[id]; !exists {
		return ErrNotFound
	}
	delete(r.employees, id)
	return nil
}