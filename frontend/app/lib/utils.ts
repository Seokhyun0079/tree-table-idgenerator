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

export function findEmployeesInDepartmentTree(
  departments: Department[],
  employees: Employee[],
  targetDepartmentId: number
): Employee[] {
  const departmentMap = new Map<number, Department>();
  const childDepartments = new Map<number, number[]>();

  // 부서 맵과 자식 부서 관계 설정
  departments.forEach(dept => {
    departmentMap.set(dept.id, dept);
    if (dept.parent_id) {
      if (!childDepartments.has(dept.parent_id)) {
        childDepartments.set(dept.parent_id, []);
      }
      childDepartments.get(dept.parent_id)!.push(dept.id);
    }
  });

  // 재귀적으로 하위 부서 ID 수집
  const getAllChildDepartmentIds = (deptId: number): number[] => {
    const result = [deptId];
    const children = childDepartments.get(deptId) || [];
    children.forEach(childId => {
      result.push(...getAllChildDepartmentIds(childId));
    });
    return result;
  };

  // 대상 부서와 모든 하위 부서의 ID 수집
  const allDepartmentIds = getAllChildDepartmentIds(targetDepartmentId);

  // 해당 부서들의 모든 직원 반환
  return employees.filter(emp => allDepartmentIds.includes(emp.department_id));
} 