"use client";

import { useEffect, useMemo, useState } from "react";
import { Department, Employee } from "./types";
import {
  getDepartments,
  getDepartmentEmployees,
  getEmployeesByDepartmentIDs,
  createEmployee,
  updateEmployee,
  deleteEmployee,
  createDepartment,
  updateDepartment,
  deleteDepartment,
} from "./lib/api";
import {
  buildDepartmentTree,
  DepartmentNode,
  findEmployeesInDepartmentTree,
} from "./lib/utils";
import DepartmentTree from "./components/DepartmentTree";
import EmployeeList from "./components/EmployeeList";
import DepartmentManager from "./components/DepartmentManager";
import DeparmentItem from "./components/DeparmentItem";

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
  const [useClientSideProcessing, setUseClientSideProcessing] = useState(false);
  const [
    useDepartmentTreeClientSideProcessing,
    setUseDepartmentTreeClientSideProcessing,
  ] = useState(false);

  const rootTree = useMemo(() => {
    return departmentTree.filter((item) => !item.parent_id);
  }, [departmentTree]);
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
  }, [useClientSideProcessing]);

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
        let employees: Employee[];

        if (useClientSideProcessing) {
          // Collect sub-department IDs
          const departmentIds = findEmployeesInDepartmentTree(
            departments,
            [],
            selectedDepartmentId
          ).map((dept) => dept.id);
          employees = await getEmployeesByDepartmentIDs(departmentIds);
        } else {
          employees = await getDepartmentEmployees(selectedDepartmentId);
        }

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
  }, [selectedDepartmentId, useClientSideProcessing]);

  useEffect(() => {
    if (departmentEmployees.length > 0) {
      const renderStartTime = performance.now();
      // Measure time after rendering is complete
      requestAnimationFrame(() => {
        const renderEndTime = performance.now();
        setRenderTime(renderEndTime - renderStartTime);
      });
    }
  }, [departmentEmployees]);

  const handleAddEmployee = async (employee: Omit<Employee, "id">) => {
    try {
      const newEmployee = await createEmployee(employee);
      if (selectedDepartmentId === employee.department_id) {
        setDepartmentEmployees((prev) => [...prev, newEmployee]);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add employee");
    }
  };

  const handleEditEmployee = async (employee: Employee) => {
    try {
      const updatedEmployee = await updateEmployee(employee);
      if (selectedDepartmentId === employee.department_id) {
        setDepartmentEmployees((prev) =>
          prev.map((emp) =>
            emp.id === updatedEmployee.id ? updatedEmployee : emp
          )
        );
      }
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
        <h1 className="text-3xl font-bold">Organization Chart</h1>
        <div className="flex gap-4">
          <div className="flex items-center gap-2">
            <label htmlFor="processingMode" className="text-sm">
              Processing Mode:
            </label>
            <select
              id="processingMode"
              value={useClientSideProcessing ? "client" : "server"}
              onChange={(e) =>
                setUseClientSideProcessing(e.target.value === "client")
              }
              className="px-2 py-1 border rounded"
            >
              <option value="server">Recurcive SQL</option>
              <option value="client">IN SQL</option>
            </select>
          </div>
          <button
            onClick={() => setShowDepartmentManager(!showDepartmentManager)}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            {showDepartmentManager
              ? "Close Department Manager"
              : "Department Manager"}
          </button>
        </div>
      </div>

      {showDepartmentManager ? (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Department Management</h2>
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
            <h2 className="text-2xl font-semibold mb-4">
              Department Structure
            </h2>
            <div className="bg-white rounded-lg shadow p-4">
              <DepartmentTree
                departments={departmentTree}
                onSelect={setSelectedDepartmentId}
              />
            </div>

            <h2 className="text-2xl font-semibold mb-4">
              Department Structure Performance
            </h2>
            <div className="flex items-center gap-2 mb-4">
              <label htmlFor="processingMode" className="text-sm">
                Department Tree Processing Mode:
              </label>
              <select
                id="departmentTreeProcessingMode"
                value={
                  useDepartmentTreeClientSideProcessing
                    ? "recursive"
                    : "comparison"
                }
                onChange={(e) =>
                  setUseDepartmentTreeClientSideProcessing(
                    e.target.value === "recursive"
                  )
                }
                className="px-2 py-1 border rounded"
              >
                <option value="recursive">Recursive SQL</option>
                <option value="comparison">Comparison SQL</option>
              </select>
            </div>
            <div className="bg-white rounded-lg shadow p-4">
              <DeparmentItem
                departments={rootTree}
                useClientSideProcessing={useDepartmentTreeClientSideProcessing}
              />
            </div>
          </div>

          <div>
            <h2 className="text-2xl font-semibold mb-4">
              {selectedDepartmentId
                ? `Employee List (Department ID: ${selectedDepartmentId})`
                : "Employee List"}
            </h2>
            {(apiCallTime || renderTime) && (
              <div className="mb-4 p-3 bg-gray-100 rounded-lg">
                <h3 className="text-lg font-medium mb-2">
                  Performance Metrics:
                </h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-gray-600">API Call Time:</p>
                    <p className="font-semibold">
                      {apiCallTime?.toFixed(2) ?? "-"} ms
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Rendering Time:</p>
                    <p className="font-semibold">
                      {renderTime?.toFixed(2) ?? "-"} ms
                    </p>
                  </div>
                  <div className="col-span-2">
                    <p className="text-sm text-gray-600">Total Time:</p>
                    <p className="font-semibold">
                      {apiCallTime && renderTime
                        ? (apiCallTime + renderTime).toFixed(2)
                        : "-"}{" "}
                      ms
                    </p>
                  </div>
                </div>
              </div>
            )}
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
