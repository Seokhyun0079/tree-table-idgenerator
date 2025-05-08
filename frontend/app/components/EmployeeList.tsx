"use client";

import { useState } from "react";
import { Employee } from "../types";

interface EmployeeListProps {
  employees: Employee[];
  departmentId: number | null;
  onAdd?: (employee: Omit<Employee, "id">) => Promise<void>;
  onEdit?: (employee: Employee) => Promise<void>;
  onDelete?: (employeeId: number) => Promise<void>;
}

export default function EmployeeList({
  employees,
  departmentId,
  onAdd,
  onEdit,
  onDelete,
}: EmployeeListProps) {
  const [isAdding, setIsAdding] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [newEmployee, setNewEmployee] = useState<Omit<Employee, "id">>({
    name: "",
    department_id: departmentId || 0,
    employee_number: "",
    position: "",
    hire_date: new Date().toISOString().split("T")[0],
  });

  const handleAdd = async () => {
    if (
      !newEmployee.name.trim() ||
      !newEmployee.employee_number.trim() ||
      !newEmployee.position.trim()
    )
      return;
    try {
      await onAdd?.(newEmployee);
      setNewEmployee({
        name: "",
        department_id: departmentId || 0,
        employee_number: "",
        position: "",
        hire_date: new Date().toISOString().split("T")[0],
      });
      setIsAdding(false);
    } catch (error) {
      console.error("Failed to add employee:", error);
    }
  };

  const handleEdit = async (employee: Employee) => {
    try {
      await onEdit?.(employee);
      setEditingId(null);
    } catch (error) {
      console.error("Failed to edit employee:", error);
    }
  };

  const handleDelete = async (employeeId: number) => {
    if (!confirm("정말로 이 직원을 삭제하시겠습니까?")) return;
    try {
      await onDelete?.(employeeId);
    } catch (error) {
      console.error("Failed to delete employee:", error);
    }
  };

  return (
    <div className="space-y-4">
      {departmentId && (
        <div className="flex justify-end">
          <button
            onClick={() => setIsAdding(true)}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            직원 추가
          </button>
        </div>
      )}

      {isAdding && (
        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold mb-2">새 직원 추가</h3>
          <div className="space-y-2">
            <input
              type="text"
              value={newEmployee.employee_number}
              onChange={(e) =>
                setNewEmployee({
                  ...newEmployee,
                  employee_number: e.target.value,
                })
              }
              placeholder="사번"
              className="w-full px-3 py-2 border rounded"
            />
            <input
              type="text"
              value={newEmployee.name}
              onChange={(e) =>
                setNewEmployee({ ...newEmployee, name: e.target.value })
              }
              placeholder="직원 이름"
              className="w-full px-3 py-2 border rounded"
            />
            <input
              type="text"
              value={newEmployee.position}
              onChange={(e) =>
                setNewEmployee({ ...newEmployee, position: e.target.value })
              }
              placeholder="직위"
              className="w-full px-3 py-2 border rounded"
            />
            <input
              type="date"
              value={newEmployee.hire_date}
              onChange={(e) =>
                setNewEmployee({ ...newEmployee, hire_date: e.target.value })
              }
              className="w-full px-3 py-2 border rounded"
            />
            <div className="flex justify-end space-x-2">
              <button
                onClick={() => setIsAdding(false)}
                className="px-3 py-1 text-gray-600 hover:text-gray-800"
              >
                취소
              </button>
              <button
                onClick={handleAdd}
                className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                추가
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="space-y-2">
        {employees.map((emp) => (
          <div
            key={emp.id}
            className="bg-white p-4 rounded-lg shadow-sm border"
          >
            {editingId === emp.id ? (
              <div className="space-y-2">
                <input
                  type="text"
                  value={emp.employee_number}
                  onChange={(e) =>
                    handleEdit({ ...emp, employee_number: e.target.value })
                  }
                  placeholder="사번"
                  className="w-full px-3 py-2 border rounded"
                />
                <input
                  type="text"
                  value={emp.name}
                  onChange={(e) => handleEdit({ ...emp, name: e.target.value })}
                  placeholder="직원 이름"
                  className="w-full px-3 py-2 border rounded"
                />
                <input
                  type="text"
                  value={emp.position}
                  onChange={(e) =>
                    handleEdit({ ...emp, position: e.target.value })
                  }
                  placeholder="직위"
                  className="w-full px-3 py-2 border rounded"
                />
                <input
                  type="date"
                  value={emp.hire_date}
                  onChange={(e) =>
                    handleEdit({ ...emp, hire_date: e.target.value })
                  }
                  className="w-full px-3 py-2 border rounded"
                />
                <div className="flex justify-end space-x-2">
                  <button
                    onClick={() => setEditingId(null)}
                    className="px-3 py-1 text-gray-600 hover:text-gray-800"
                  >
                    취소
                  </button>
                  <button
                    onClick={() => handleEdit(emp)}
                    className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                  >
                    저장
                  </button>
                </div>
              </div>
            ) : (
              <div className="flex justify-between items-center">
                <div>
                  <div className="font-medium">{emp.name}</div>
                  <div className="text-sm text-gray-500">
                    ID: {emp.id} | 부서 ID: {emp.department_id}
                  </div>
                  <div className="text-sm text-gray-500">
                    사번: {emp.employee_number} | 직위: {emp.position} | 입사일:{" "}
                    {emp.hire_date}
                  </div>
                </div>
                <div className="flex space-x-2">
                  <button
                    onClick={() => setEditingId(emp.id)}
                    className="px-3 py-1 text-blue-500 hover:text-blue-600"
                  >
                    수정
                  </button>
                  <button
                    onClick={() => handleDelete(emp.id)}
                    className="px-3 py-1 text-red-500 hover:text-red-600"
                  >
                    삭제
                  </button>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {employees.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          {departmentId
            ? "해당 부서에 소속된 직원이 없습니다."
            : "부서를 선택하여 직원 목록을 확인하세요."}
        </div>
      )}
    </div>
  );
}
