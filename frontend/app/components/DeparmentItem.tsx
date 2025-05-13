"use client";

import { useEffect, useState } from "react";
import { DepartmentNode } from "../lib/utils";
import {
  getChildDepartmentComparison,
  getChildDepartmentRecursive,
} from "../lib/api";
import { Department } from "../types";

interface DepartmentItemProps {
  departments: DepartmentNode[];
  level?: number;
  onSelect?: (departmentId: number) => void;
  useClientSideProcessing?: boolean;
}

export default function DeparmentItem({
  departments,
  useClientSideProcessing,
}: DepartmentItemProps) {
  const [expandedDepartments, setExpandedDepartments] = useState<
    Record<number, Department[]>
  >({});

  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [apiCallTime, setApiCallTime] = useState<number | null>(null);
  const [renderTime, setRenderTime] = useState<number | null>(null);

  const onClickEvent = async (id: number) => {
    const willExpand = !expanded.has(id);
    toggleExpand(id);

    if (willExpand) {
      const startTime = performance.now();
      const data = useClientSideProcessing
        ? await getChildDepartmentRecursive(id)
        : await getChildDepartmentComparison(id);
      const apiEndTime = performance.now();
      setApiCallTime(apiEndTime - startTime);

      const newChildDepartments = data.map((dept: Department) => ({
        ...dept,
        children: [],
      }));
      setExpandedDepartments((prev) => ({
        ...prev,
        [id]: newChildDepartments,
      }));
    }
  };

  useEffect(() => {
    if (Object.keys(expandedDepartments).length > 0) {
      const renderStartTime = performance.now();
      requestAnimationFrame(() => {
        const renderEndTime = performance.now();
        setRenderTime(renderEndTime - renderStartTime);
      });
    }
  }, [expandedDepartments]);

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
    <div className="space-y-4">
      {(apiCallTime || renderTime) && (
        <div className="p-3 bg-gray-100 rounded-lg">
          <h3 className="text-lg font-medium mb-2">Performance Metrics:</h3>
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
      <div className="space-y-1">
        {departments.map((dept) => (
          <div key={dept.id}>
            <div
              className="flex items-center py-1"
              onClick={() => onClickEvent(dept.id)}
            >
              {dept.children.length > 0 && (
                <button
                  onClick={() => onClickEvent(dept.id)}
                  className="w-6 h-6 flex items-center justify-center text-gray-500 hover:text-gray-700"
                >
                  {expanded.has(dept.id) ? "▼" : "▶"}
                </button>
              )}
              <div
                className={`flex-1 px-2 py-1 rounded cursor-pointer hover:bg-gray-100`}
              >
                <div className="font-medium">{dept.name}</div>
                <div className="text-sm text-gray-500">
                  ID: {dept.id} | Parent ID: {dept.parent_id}
                </div>
              </div>
            </div>

            {expanded.has(dept.id) &&
              expandedDepartments[dept.id]?.length > 0 && (
                <div className="ml-4">
                  {expandedDepartments[dept.id].map((child) => (
                    <div key={child.id}>
                      <div className="flex items-center py-1">
                        <div className="text-sm text-gray-500">
                          ID: {child.id} | Parent ID: {child.parent_id}
                        </div>
                        <div className="font-medium">-{child.name}</div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
          </div>
        ))}
      </div>
    </div>
  );
}
