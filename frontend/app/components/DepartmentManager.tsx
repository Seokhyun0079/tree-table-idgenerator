"use client";

import { useState, useEffect } from "react";
import { Department } from "../types";

interface DepartmentManagerProps {
  departments: Department[];
  onAdd?: (department: Omit<Department, "id">) => Promise<void>;
  onEdit?: (department: Department) => Promise<void>;
  onDelete?: (departmentId: number) => Promise<void>;
}

interface DepartmentWithLevel extends Department {
  level: number;
}

export default function DepartmentManager({
  departments,
  onAdd,
  onEdit,
  onDelete,
}: DepartmentManagerProps) {
  const [isAdding, setIsAdding] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [newDepartment, setNewDepartment] = useState<Omit<Department, "id">>({
    name: "",
    parent_id: null,
  });
  const [selectedLevel, setSelectedLevel] = useState<number>(0);

  const handleParentSelect = (parentId: number) => {
    setSelectedLevel(parentId);
    setNewDepartment((prev) => ({ ...prev, parent_id: parentId }));
  };

  const handleAdd = async () => {
    if (!newDepartment.name.trim()) return;
    try {
      await onAdd?.(newDepartment);
      setNewDepartment({ name: "", parent_id: null });
      setSelectedLevel(0);
      setIsAdding(false);
    } catch (error) {
      console.error("Failed to add department:", error);
    }
  };

  const handleEdit = async (department: Department) => {
    try {
      await onEdit?.(department);
      setEditingId(null);
    } catch (error) {
      console.error("Failed to edit department:", error);
    }
  };

  const handleDelete = async (departmentId: number) => {
    if (!confirm("Are you sure you want to delete this department?")) return;
    try {
      await onDelete?.(departmentId);
    } catch (error) {
      console.error("Failed to delete department:", error);
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <button
          onClick={() => setIsAdding(true)}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Add Department
        </button>
      </div>

      {isAdding && (
        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold mb-2">Add New Department</h3>
          <div className="space-y-2">
            <input
              type="text"
              value={newDepartment.name}
              onChange={(e) =>
                setNewDepartment({ ...newDepartment, name: e.target.value })
              }
              placeholder="Department Name"
              className="w-full px-3 py-2 border rounded"
            />
            <div className="space-y-2">
              <select
                value={selectedLevel || ""}
                onChange={(e) => handleParentSelect(parseInt(e.target.value))}
                className="w-full px-3 py-2 border rounded"
              >
                <option key={0} value={undefined}>
                  {"Select Parent Department"}
                </option>
                {departments.map((dept) => (
                  <option key={dept.id} value={dept.id}>
                    {dept.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex justify-end space-x-2">
              <button
                onClick={() => {
                  setIsAdding(false);
                  setNewDepartment({ name: "", parent_id: null });
                }}
                className="px-3 py-1 text-gray-600 hover:text-gray-800"
              >
                Cancel
              </button>
              <button
                onClick={handleAdd}
                className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                Add
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="space-y-2">
        {departments.map((dept) => (
          <div
            key={dept.id}
            className="bg-white p-4 rounded-lg shadow-sm border"
          >
            {editingId === dept.id ? (
              <div className="space-y-2">
                <input
                  type="text"
                  value={dept.name}
                  onChange={(e) =>
                    handleEdit({ ...dept, name: e.target.value })
                  }
                  className="w-full px-3 py-2 border rounded"
                />
                <select
                  value={dept.parent_id || ""}
                  onChange={(e) =>
                    handleEdit({
                      ...dept,
                      parent_id: e.target.value
                        ? parseInt(e.target.value)
                        : null,
                    })
                  }
                  className="w-full px-3 py-2 border rounded"
                >
                  {departments.map((d) => (
                    <option key={d.id} value={d.id}>
                      {d.name}
                    </option>
                  ))}
                </select>
                <div className="flex justify-end space-x-2">
                  <button
                    onClick={() => setEditingId(null)}
                    className="px-3 py-1 text-gray-600 hover:text-gray-800"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            ) : (
              <div className="flex justify-between items-center">
                <div>
                  <div className="font-medium">{dept.name}</div>
                  <div className="text-sm text-gray-500">
                    ID: {dept.id} | Parent Department ID: {dept.parent_id}
                  </div>
                </div>
                <div className="flex space-x-2">
                  <button
                    onClick={() => setEditingId(dept.id)}
                    className="px-3 py-1 text-blue-500 hover:text-blue-600"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(dept.id)}
                    className="px-3 py-1 text-red-500 hover:text-red-600"
                  >
                    Delete
                  </button>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
