import { NextResponse } from "next/server";
import prisma from "@/lib/prisma";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const parentId = searchParams.get('parentId');

  if (!parentId) {
    return NextResponse.json({ error: 'Parent ID is required' }, { status: 400 });
  }

  try {
    const departments = await prisma.$queryRaw`
      WITH RECURSIVE department_tree AS (
        -- 기본 케이스: 선택된 부모 부서
        SELECT id, name, parent_id, 0 as level
        FROM departments
        WHERE id = ${parseInt(parentId)}
        
        UNION ALL
        
        -- 재귀 케이스: 하위 부서들
        SELECT d.id, d.name, d.parent_id, dt.level + 1
        FROM departments d
        INNER JOIN department_tree dt ON d.parent_id = dt.id
      )
      SELECT * FROM department_tree
      ORDER BY level, id;
    `;

    return NextResponse.json(departments);
  } catch (error) {
    console.error('Error fetching department tree:', error);
    return NextResponse.json(
      { error: 'Failed to fetch department tree' },
      { status: 500 }
    );
  }
} 