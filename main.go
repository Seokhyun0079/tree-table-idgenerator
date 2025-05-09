package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Employee 구조체 정의
type Employee struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DepartmentID   int    `json:"department_id"`
	Position       string `json:"position"`
	HireDate       string `json:"hire_date"`
	EmployeeNumber string `json:"employee_number"`
	LargeText      string `json:"large_text"`
}

// Department 구조체 정의
type Department struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
}


func initDB() {
	var err error
	// 환경 변수에서 데이터베이스 연결 정보 가져오기
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "rootpassword")
	dbName := getEnv("DB_NAME", "mydatabase")
	
	// 문자열을 숫자로 변환
	retryIntervalStr := getEnv("DB_RETRY_INTERVAL", "10")
	maxRetriesStr := getEnv("DB_MAX_RETRIES", "30")
	
	retryInterval, err := strconv.Atoi(retryIntervalStr)
	if err != nil {
		log.Printf("Invalid DB_RETRY_INTERVAL, using default value: %v", err)
		retryInterval = 10
	}
	
	maxRetries, err := strconv.Atoi(maxRetriesStr)
	if err != nil {
		log.Printf("Invalid DB_MAX_RETRIES, using default value: %v", err)
		maxRetries = 30
	}

	// 데이터베이스 연결 문자열 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// 재시도 로직
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(time.Duration(retryInterval) * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(time.Duration(retryInterval) * time.Second)
			continue
		}

		fmt.Println("Successfully connected to database!")
		return
	}

	log.Fatal("Failed to connect to database after maximum retries")
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

	// Health check 엔드포인트 추가
	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Database connection failed"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 라우트 설정
	api := r.Group("/api")
	{
		// 부서 관련 API
		api.GET("/departments", getDepartments)
		api.GET("/departments/:id", getDepartment)
		api.GET("/departments/:id/employees", getDepartmentEmployees)
		api.POST("/departments", createDepartment)

		// 직원 관련 API
		api.GET("/employees", getEmployees)
		api.GET("/employees/:id", getEmployee)

		// 새로운 엔드포인트 추가
		api.POST("/employees/by-departments", GetEmployeesByDepartmentIDs)

		// 부서 트리 조회 API
		api.GET("/departments/tree", func(c *gin.Context) {
			parentId := c.Query("parentId")
			if parentId == "" {
				c.JSON(400, gin.H{"error": "Parent ID is required"})
				return
			}

			query := `
				WITH RECURSIVE department_tree AS (
					-- 기본 케이스: 선택된 부모 부서
					SELECT id, name, parent_id, 0 as level
					FROM departments
					WHERE id = ?
					
					UNION ALL
					
					-- 재귀 케이스: 하위 부서들
					SELECT d.id, d.name, d.parent_id, dt.level + 1
					FROM departments d
					INNER JOIN department_tree dt ON d.parent_id = dt.id
				)
				SELECT * FROM department_tree
				ORDER BY level, id;
			`

			rows, err := db.Query(query, parentId)
			if err != nil {
				log.Printf("Error querying department tree: %v", err)
				c.JSON(500, gin.H{"error": "Failed to fetch department tree"})
				return
			}
			defer rows.Close()

			var departments []gin.H
			for rows.Next() {
				var id int
				var name string
				var parentID sql.NullInt64
				var level int
				if err := rows.Scan(&id, &name, &parentID, &level); err != nil {
					log.Printf("Error scanning department row: %v", err)
					c.JSON(500, gin.H{"error": "Failed to scan department row"})
					return
				}
				departments = append(departments, gin.H{
					"id":        id,
					"name":      name,
					"parent_id": parentID.Int64,
					"level":     level,
				})
			}

			if err = rows.Err(); err != nil {
				log.Printf("Error iterating department rows: %v", err)
				c.JSON(500, gin.H{"error": "Failed to iterate department rows"})
				return
			}

			c.JSON(200, departments)
		})
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
	
	// 실행 계획 확인을 위한 쿼리
	explainQuery := `
		EXPLAIN FORMAT=JSON
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
		SELECT e.id, e.name, e.department_id, e.position, e.hire_date, e.employee_number, e.large_text
		FROM employees e
		INNER JOIN subdepartments sd ON e.department_id = sd.id
		ORDER BY e.department_id, e.name
	`
	
	// 실행 계획 출력
	rows, err := db.Query(explainQuery, deptID)
	if err != nil {
		log.Printf("Error explaining query: %v", err)
	} else {
		defer rows.Close()
		var explainResult string
		for rows.Next() {
			if err := rows.Scan(&explainResult); err != nil {
				log.Printf("Error scanning explain result: %v", err)
			} else {
				log.Printf("Query execution plan for department %s: %s", deptID, explainResult)
			}
		}
	}
	
	// 실제 쿼리 실행
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
		SELECT e.id, e.name, e.department_id, e.position, e.hire_date, e.employee_number, e.large_text
		FROM employees e
		INNER JOIN subdepartments sd ON e.department_id = sd.id
		ORDER BY e.department_id, e.name
	`
	
	rows, err = db.Query(query, deptID)
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
		var largeText sql.NullString
		if err := rows.Scan(&id, &name, &deptID, &position, &hireDate, &employeeNumber, &largeText); err != nil {
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
			"large_text":     largeText.String,
		})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating employee rows: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to iterate employee rows: %v", err)})
		return
	}

	c.JSON(200, employees)
}
const max_id_length int = 9

func createDepartment(c *gin.Context) {
	var dept Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 부모 ID가 배열인 경우를 처리
	var parentIDs []interface{}
	var newID int
	if dept.ParentID != 0 {
		// 부모 ID의 자리수에 따라 increment 결정
		increment := 1
		tempID := dept.ParentID
		strID := strconv.Itoa(tempID)  // "123"
		log.Printf("strID: %v", strID)
		
		if strID[len(strID)-1] != '0' {
			c.JSON(400, gin.H{"error": "Failed to create department"})
			return
		}

		for i := len(strID)- 1; i > 0 && strID[i-1] == '0'; i-- {
			increment *= 10
		}
		
		log.Printf("increment: %v", increment)

		for i := 1; i < max_id_length; i++ {
			// 부모 ID에 증가값을 곱해서 더함
			// 예: 900 -> 910, 920, 930...
			// 예: 1000 -> 1100, 1200, 1300...
			calcId := dept.ParentID + (i * increment)
			parentIDs = append(parentIDs, calcId)
		}

		// IN 절을 위한 플레이스홀더 생성
		placeholders := make([]string, len(parentIDs))
		for i := range parentIDs {
			placeholders[i] = "?"
		}

		// 쿼리 실행
		query := fmt.Sprintf("SELECT id FROM departments WHERE id IN (%s) ORDER BY id", strings.Join(placeholders, ","))
		log.Printf("query: %v", query)
		rows, err := db.Query(query, parentIDs...)
		if err != nil {
			log.Printf("Error querying departments: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query departments: %v", err)})
			return
		}
		defer rows.Close()
		for i := 1; i < max_id_length; i++ {
			var calcId int = dept.ParentID + (i * increment)
			var id int	
			if !rows.Next() {
				// 더 이상 비교할 ID가 없으면 마지막 계산된 ID를 사용
				newID = calcId
				break
			}
			err := rows.Scan(&id)
			if err != nil {
				log.Printf("Error scanning count: %v", err)
				c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to scan count: %v", err)})
				return
			}


			//계산한 아이디와 조회한 아이디가 다르다는 것은 중간에 비어있는 아이디가 있다는 것이므로 그 아이디를 할당함
			//예 : 계산한 아이디 1100, 조회한 아이디 1100 OK
			//예 : 계산한 아이디 1300, 조회한 아이디 1400 중간에 1300이 비어있음 그래서 1300을 할당
			log.Printf("calcId: %v, id: %v", calcId, id)
			if calcId != id {
				newID = calcId
				break
			}
		}
		log.Printf("newID: %v", newID)
	}


	log.Printf("dept.Name: %v", dept)
	if newID == 0 {
		c.JSON(400, gin.H{"error": "Failed to create department"})
		return
	}

	result, err := db.Exec("INSERT INTO departments (id, name, parent_id) VALUES (?, ?, ?)", newID, dept.Name, dept.ParentID)
	if err != nil {
		log.Printf("Error creating department: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create department: %v", err)})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get last insert ID: %v", err)})
		return
	}

	c.JSON(200, gin.H{
		"message": "Department created successfully",
		"id": id,
	})
}


// 직원 목록 조회
func getEmployees(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, department_id, large_text FROM employees")
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
		var largeText sql.NullString
		if err := rows.Scan(&id, &name, &deptID, &largeText); err != nil {
			log.Printf("Error scanning employee row: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to scan employee row: %v", err)})
			return
		}
		employees = append(employees, gin.H{
			"id":            id,
			"name":          name,
			"department_id": deptID,
			"large_text":    largeText.String,
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
	var largeText sql.NullString
	err := db.QueryRow("SELECT id, name, department_id, large_text FROM employees WHERE id = ?", id).
		Scan(&empID, &name, &deptID, &largeText)
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
		"large_text":    largeText.String,
	}
	c.JSON(200, employee)
}

// GetEmployeesByDepartmentIDs handles GET request for employees by department IDs
func GetEmployeesByDepartmentIDs(c *gin.Context) {
	var departmentIDs []int
	if err := c.ShouldBindJSON(&departmentIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// IN 절을 위한 플레이스홀더 생성
	placeholders := make([]string, len(departmentIDs))
	args := make([]interface{}, len(departmentIDs))
	for i, id := range departmentIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT e.id, e.name, e.department_id, e.position, e.hire_date, e.employee_number, e.large_text
		FROM employees e
		WHERE e.department_id IN (%s)
		ORDER BY e.department_id, e.id
	`, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.DepartmentID, &emp.Position, &emp.HireDate, &emp.EmployeeNumber, &emp.LargeText); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, emp)
	}

	c.JSON(http.StatusOK, employees)
} 