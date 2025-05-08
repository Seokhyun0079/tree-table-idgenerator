package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	// 환경 변수에서 데이터베이스 연결 정보 가져오기
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "rootpassword")
	dbName := getEnv("DB_NAME", "mydatabase")

	// 데이터베이스 연결 문자열 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to database!")
}

// 환경 변수 가져오기 (기본값 설정)
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// CORS 미들웨어 추가
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API 라우트 설정
	api := r.Group("/api")
	{
		// 부서 관련 API
		api.GET("/departments", getDepartments)
		api.GET("/departments/:id", getDepartment)
		api.GET("/departments/:id/employees", getDepartmentEmployees)

		// 직원 관련 API
		api.GET("/employees", getEmployees)
		api.GET("/employees/:id", getEmployee)
	}

	r.Run(":8080")
}

// 부서 목록 조회
func getDepartments(c *gin.Context) {
	rows, err := db.Query("SELECT id, parent_id, name FROM departments")
	if err != nil {
		log.Printf("Error querying departments: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query departments: %v", err)})
		return
	}
	defer rows.Close()

	var departments []gin.H
	for rows.Next() {
		var id int
		var parentID sql.NullInt64
		var name string
		if err := rows.Scan(&id, &parentID, &name); err != nil {
			log.Printf("Error scanning department row: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to scan department row: %v", err)})
			return
		}
		departments = append(departments, gin.H{
			"id":        id,
			"parent_id": parentID.Int64,
			"name":      name,
		})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating department rows: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to iterate department rows: %v", err)})
		return
	}

	c.JSON(200, departments)
}

// 특정 부서 조회
func getDepartment(c *gin.Context) {
	id := c.Param("id")
	var deptID int
	var parentID sql.NullInt64
	var name string
	err := db.QueryRow("SELECT id, parent_id, name FROM departments WHERE id = ?", id).
		Scan(&deptID, &parentID, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Department not found"})
		} else {
			log.Printf("Error querying department: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query department: %v", err)})
		}
		return
	}
	department := gin.H{
		"id":        deptID,
		"parent_id": parentID.Int64,
		"name":      name,
	}
	c.JSON(200, department)
}

// 특정 부서의 직원 목록 조회
func getDepartmentEmployees(c *gin.Context) {
	deptID := c.Param("id")
	
	// 재귀적으로 하위 부서 ID들을 가져오는 쿼리
	query := `
		WITH RECURSIVE subdepartments AS (
			-- 기본 부서
			SELECT id, parent_id
			FROM departments
			WHERE id = ?
			
			UNION ALL
			
			-- 하위 부서들
			SELECT d.id, d.parent_id
			FROM departments d
			INNER JOIN subdepartments sd ON d.parent_id = sd.id
		)
		SELECT e.id, e.name, e.department_id, e.position, e.hire_date, e.employee_number
		FROM employees e
		INNER JOIN subdepartments sd ON e.department_id = sd.id
		ORDER BY e.department_id, e.name
	`
	
	rows, err := db.Query(query, deptID)
	if err != nil {
		log.Printf("Error querying department employees: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query department employees: %v", err)})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, deptID int
		var name, position, employeeNumber string
		var hireDate string
		if err := rows.Scan(&id, &name, &deptID, &position, &hireDate, &employeeNumber); err != nil {
			log.Printf("Error scanning employee row: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to scan employee row: %v", err)})
			return
		}
		employees = append(employees, gin.H{
			"id":             id,
			"name":           name,
			"department_id":  deptID,
			"position":       position,
			"hire_date":      hireDate,
			"employee_number": employeeNumber,
		})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating employee rows: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to iterate employee rows: %v", err)})
		return
	}

	c.JSON(200, employees)
}

// 직원 목록 조회
func getEmployees(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, department_id FROM employees")
	if err != nil {
		log.Printf("Error querying employees: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query employees: %v", err)})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, deptID int
		var name string
		if err := rows.Scan(&id, &name, &deptID); err != nil {
			log.Printf("Error scanning employee row: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to scan employee row: %v", err)})
			return
		}
		employees = append(employees, gin.H{
			"id":            id,
			"name":          name,
			"department_id": deptID,
		})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating employee rows: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to iterate employee rows: %v", err)})
		return
	}

	c.JSON(200, employees)
}

// 특정 직원 조회
func getEmployee(c *gin.Context) {
	id := c.Param("id")
	var empID, deptID int
	var name string
	err := db.QueryRow("SELECT id, name, department_id FROM employees WHERE id = ?", id).
		Scan(&empID, &name, &deptID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Employee not found"})
		} else {
			log.Printf("Error querying employee: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query employee: %v", err)})
		}
		return
	}
	employee := gin.H{
		"id":            empID,
		"name":          name,
		"department_id": deptID,
	}
	c.JSON(200, employee)
} 