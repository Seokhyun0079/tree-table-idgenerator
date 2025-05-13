import { Department, Employee } from '../types';

export interface DepartmentNode extends Department {
  children: DepartmentNode[];
  employees: number[];
}

export function buildDepartmentTree(departments: Department[]): DepartmentNode[] {
  const departmentMap = new Map<number, DepartmentNode>();
  const rootDepartments: DepartmentNode[] = [];

  // First add all departments to the map
  departments.forEach((dept) => {
    departmentMap.set(dept.id, {
      ...dept,
      children: [],
      employees: [],
    });
  });

  // Create department tree structure
  departments.forEach((dept) => {
    const node = departmentMap.get(dept.id)!;
    if (dept.parent_id === null) {
      rootDepartments.push(node);
    } else {
      const parent = departmentMap.get(dept.parent_id);
      if (parent) {
        parent.children.push(node);
      } else {
        // If no parent, add as root
        rootDepartments.push(node);
      }
    }
  });

  return rootDepartments;
}

export function findDepartmentPath(
  departments: Department[],
  targetId: number
): number[] {
  const path: number[] = [];
  let currentId = targetId;

  while (currentId !== 0) {
    const dept = departments.find(d => d.id === currentId);
    if (!dept || dept.parent_id === null) break;

    path.unshift(currentId);
    currentId = dept.parent_id;
  }

  return path;
}

export function findEmployeesInDepartmentTree(
  departments: Department[],
  employees: Employee[],
  targetDepartmentId: number
): Department[] {
  const departmentMap = new Map<number, Department>();
  const childDepartments = new Map<number, number[]>();

  // Set up department map and child department relationships
  departments.forEach(dept => {
    departmentMap.set(dept.id, dept);
    if (dept.parent_id) {
      if (!childDepartments.has(dept.parent_id)) {
        childDepartments.set(dept.parent_id, []);
      }
      childDepartments.get(dept.parent_id)!.push(dept.id);
    }
  });

  // Recursively collect child department IDs
  const getAllChildDepartmentIds = (deptId: number): number[] => {
    const result = [deptId];
    const children = childDepartments.get(deptId) || [];
    children.forEach(childId => {
      result.push(...getAllChildDepartmentIds(childId));
    });
    return result;
  };

  // Collect IDs of target department and all its child departments
  const allDepartmentIds = getAllChildDepartmentIds(targetDepartmentId);

  // Return all department objects for the collected IDs
  return departments.filter(dept => allDepartmentIds.includes(dept.id));
} 