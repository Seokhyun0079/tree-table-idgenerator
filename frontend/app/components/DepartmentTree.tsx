"use client";

import { useState } from "react";
import { DepartmentNode } from "../lib/utils";

interface DepartmentTreeProps {
  departments: DepartmentNode[];
  level?: number;
  onSelect?: (departmentId: number) => void;
}

export default function DepartmentTree({
  departments,
  level = 0,
  onSelect,
}: DepartmentTreeProps) {
  const [expanded, setExpanded] = useState<Set<number>>(new Set());

  const toggleExpand = (id: number) => {
    setExpanded((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  return (
    <div className="space-y-1">
      {departments.map((dept) => (
        <div key={dept.id} style={{ marginLeft: `${level * 1.5}rem` }}>
          <div className="flex items-center py-1">
            {dept.children.length > 0 && (
              <button
                onClick={() => toggleExpand(dept.id)}
                className="w-6 h-6 flex items-center justify-center text-gray-500 hover:text-gray-700"
              >
                {expanded.has(dept.id) ? "▼" : "▶"}
              </button>
            )}
            <div
              className={`flex-1 px-2 py-1 rounded cursor-pointer hover:bg-gray-100 ${
                onSelect ? "cursor-pointer" : ""
              }`}
              onClick={() => onSelect?.(dept.id)}
            >
              <div className="font-medium">{dept.name}</div>
              <div className="text-sm text-gray-500">
                ID: {dept.id} | Parent ID: {dept.parent_id}
              </div>
            </div>
          </div>
          {expanded.has(dept.id) && dept.children.length > 0 && (
            <DepartmentTree
              departments={dept.children}
              level={level + 1}
              onSelect={onSelect}
            />
          )}
        </div>
      ))}
    </div>
  );
}
