import { Department } from '../types';

export interface DepartmentNode extends Department {
  children: DepartmentNode[];
  employees: number[];
}

export function buildDepartmentTree(departments: Department[]): DepartmentNode[] {
  const departmentMap = new Map<number, DepartmentNode>();
  const rootDepartments: DepartmentNode[] = [];

  // 먼저 모든 부서를 맵에 추가
  departments.forEach((dept) => {
    departmentMap.set(dept.id, {
      ...dept,
      children: [],
      employees: [],
    });
  });

  // 부서 트리 구조 생성
  departments.forEach((dept) => {
    const node = departmentMap.get(dept.id)!;
    if (!dept.parent_id) {
      rootDepartments.push(node);
    } else {
      const parent = departmentMap.get(dept.parent_id);
      if (parent) {
        parent.children.push(node);
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