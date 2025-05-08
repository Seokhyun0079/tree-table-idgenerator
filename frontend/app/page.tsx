"use client";

import { useEffect, useState } from "react";
import { Department, Employee } from "./types";
import {
  getDepartments,
  getDepartmentEmployees,
  createEmployee,
  updateEmployee,
  deleteEmployee,
  createDepartment,
  updateDepartment,
  deleteDepartment,
} from "./lib/api";
import { buildDepartmentTree, DepartmentNode } from "./lib/utils";
import DepartmentTree from "./components/DepartmentTree";
import EmployeeList from "./components/EmployeeList";
import DepartmentManager from "./components/DepartmentManager";

export default function Home() {
  const [departments, setDepartments] = useState<Department[]>([]);
  const [departmentTree, setDepartmentTree] = useState<DepartmentNode[]>([]);
  const [selectedDepartmentId, setSelectedDepartmentId] = useState<
    number | null
  >(null);
  const [departmentEmployees, setDepartmentEmployees] = useState<Employee[]>(
    []
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showDepartmentManager, setShowDepartmentManager] = useState(false);
  const [apiCallTime, setApiCallTime] = useState<number | null>(null);
  const [renderTime, setRenderTime] = useState<number | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const departmentsData = await getDepartments();
        setDepartments(departmentsData);
        setDepartmentTree(buildDepartmentTree(departmentsData));
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const fetchDepartmentEmployees = async () => {
      if (!selectedDepartmentId) {
        setDepartmentEmployees([]);
        setApiCallTime(null);
        setRenderTime(null);
        return;
      }

      const startTime = performance.now();
      try {
        const employees = await getDepartmentEmployees(selectedDepartmentId);
        const apiEndTime = performance.now();
        setApiCallTime(apiEndTime - startTime);
        setDepartmentEmployees(employees);
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Failed to fetch department employees"
        );
        setApiCallTime(null);
        setRenderTime(null);
      }
    };

    fetchDepartmentEmployees();
  }, [selectedDepartmentId]);

  useEffect(() => {
    if (departmentEmployees.length > 0) {
      const renderStartTime = performance.now();
      // 렌더링이 완료된 후 시간 측정
      requestAnimationFrame(() => {
        const renderEndTime = performance.now();
        setRenderTime(renderEndTime - renderStartTime);
      });
    }
  }, [departmentEmployees]);

  const handleAddEmployee = async (employee: Omit<Employee, "id">) => {
    try {
      const newEmployee = await createEmployee(employee);
      setDepartmentEmployees((prev) => [...prev, newEmployee]);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add employee");
    }
  };

  const handleEditEmployee = async (employee: Employee) => {
    try {
      const updatedEmployee = await updateEmployee(employee);
      setDepartmentEmployees((prev) =>
        prev.map((emp) =>
          emp.id === updatedEmployee.id ? updatedEmployee : emp
        )
      );
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to update employee"
      );
    }
  };

  const handleDeleteEmployee = async (employeeId: number) => {
    try {
      await deleteEmployee(employeeId);
      setDepartmentEmployees((prev) =>
        prev.filter((emp) => emp.id !== employeeId)
      );
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to delete employee"
      );
    }
  };

  const handleAddDepartment = async (department: Omit<Department, "id">) => {
    try {
      const newDepartment = await createDepartment(department);
      setDepartments((prev) => [...prev, newDepartment]);
      setDepartmentTree(buildDepartmentTree([...departments, newDepartment]));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add department");
    }
  };

  const handleEditDepartment = async (department: Department) => {
    try {
      const updatedDepartment = await updateDepartment(department);
      const newDepartments = departments.map((dept) =>
        dept.id === updatedDepartment.id ? updatedDepartment : dept
      );
      setDepartments(newDepartments);
      setDepartmentTree(buildDepartmentTree(newDepartments));
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to update department"
      );
    }
  };

  const handleDeleteDepartment = async (departmentId: number) => {
    try {
      await deleteDepartment(departmentId);
      const newDepartments = departments.filter(
        (dept) => dept.id !== departmentId
      );
      setDepartments(newDepartments);
      setDepartmentTree(buildDepartmentTree(newDepartments));
      if (selectedDepartmentId === departmentId) {
        setSelectedDepartmentId(null);
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to delete department"
      );
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl text-red-500">Error: {error}</div>
      </div>
    );
  }

  return (
    <main className="min-h-screen p-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">조직도</h1>
        <button
          onClick={() => setShowDepartmentManager(!showDepartmentManager)}
          className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
        >
          {showDepartmentManager ? "부서 관리 닫기" : "부서 관리"}
        </button>
      </div>

      {showDepartmentManager ? (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">부서 관리</h2>
          <div className="bg-white rounded-lg shadow p-4">
            <DepartmentManager
              departments={departments}
              onAdd={handleAddDepartment}
              onEdit={handleEditDepartment}
              onDelete={handleDeleteDepartment}
            />
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div>
            <h2 className="text-2xl font-semibold mb-4">부서 구조</h2>
            <div className="bg-white rounded-lg shadow p-4">
              <DepartmentTree
                departments={departmentTree}
                onSelect={setSelectedDepartmentId}
              />
            </div>
          </div>

          <div>
            <h2 className="text-2xl font-semibold mb-4">
              {selectedDepartmentId
                ? `직원 목록 (부서 ID: ${selectedDepartmentId})${
                    apiCallTime || renderTime
                      ? ` - API 호출: ${
                          apiCallTime?.toFixed(2) ?? "-"
                        }ms, 렌더링: ${renderTime?.toFixed(2) ?? "-"}ms${
                          apiCallTime && renderTime
                            ? `, 총: ${(apiCallTime + renderTime).toFixed(2)}ms`
                            : ""
                        }`
                      : ""
                  }`
                : "직원 목록"}
            </h2>
            <div className="bg-white rounded-lg shadow p-4">
              <EmployeeList
                employees={departmentEmployees}
                departmentId={selectedDepartmentId}
                onAdd={handleAddEmployee}
                onEdit={handleEditEmployee}
                onDelete={handleDeleteEmployee}
              />
            </div>
          </div>
        </div>
      )}
    </main>
  );
}
