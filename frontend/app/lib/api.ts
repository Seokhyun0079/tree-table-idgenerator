import { Department, Employee } from '../types';

const API_BASE_URL = 'http://localhost:8080/api';

// 부서 관련 API
export async function getDepartments(): Promise<Department[]> {
  const response = await fetch(`${API_BASE_URL}/departments`);
  if (!response.ok) {
    throw new Error('Failed to fetch departments');
  }
  return response.json();
}

export async function getDepartment(id: number): Promise<Department> {
  const response = await fetch(`${API_BASE_URL}/departments/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch department');
  }
  return response.json();
}

export async function createDepartment(department: Omit<Department, 'id'>): Promise<Department> {
  const response = await fetch(`${API_BASE_URL}/departments`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(department),
  });
  if (!response.ok) {
    throw new Error('Failed to create department');
  }
  return response.json();
}

export async function updateDepartment(department: Department): Promise<Department> {
  const response = await fetch(`${API_BASE_URL}/departments/${department.id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(department),
  });
  if (!response.ok) {
    throw new Error('Failed to update department');
  }
  return response.json();
}

export async function deleteDepartment(id: number): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/departments/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Failed to delete department');
  }
}

// 직원 관련 API
export async function getEmployees(page: number = 1, pageSize: number = 100): Promise<Employee[]> {
  const response = await fetch(`${API_BASE_URL}/employees?page=${page}&pageSize=${pageSize}`);
  if (!response.ok) {
    throw new Error('Failed to fetch employees');
  }
  return response.json();
}

export async function getEmployee(id: number): Promise<Employee> {
  const response = await fetch(`${API_BASE_URL}/employees/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch employee');
  }
  return response.json();
}

export async function getDepartmentEmployees(departmentId: number): Promise<Employee[]> {
  const response = await fetch(`${API_BASE_URL}/departments/${departmentId}/employees`);
  if (!response.ok) {
    throw new Error('Failed to fetch department employees');
  }
  return response.json();
}

export async function createEmployee(employee: Omit<Employee, 'id'>): Promise<Employee> {
  const response = await fetch(`${API_BASE_URL}/employees`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(employee),
  });
  if (!response.ok) {
    throw new Error('Failed to create employee');
  }
  return response.json();
}

export async function updateEmployee(employee: Employee): Promise<Employee> {
  const response = await fetch(`${API_BASE_URL}/employees/${employee.id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(employee),
  });
  if (!response.ok) {
    throw new Error('Failed to update employee');
  }
  return response.json();
}

export async function deleteEmployee(id: number): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/employees/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Failed to delete employee');
  }
}

export async function getEmployeesByDepartmentIDs(departmentIds: number[]): Promise<Employee[]> {
  const response = await fetch(`${API_BASE_URL}/employees/by-departments`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(departmentIds),
  });
  if (!response.ok) {
    throw new Error('Failed to fetch employees by department IDs');
  }
  return response.json();
} 