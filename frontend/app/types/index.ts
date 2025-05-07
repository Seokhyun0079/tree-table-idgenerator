export interface Department {
  id: number;
  parent_id: number | null;
  name: string;
}

export interface Employee {
  id: number;
  name: string;
  department_id: number;
} 