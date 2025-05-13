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

// Employee struct definition
type Employee struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DepartmentID   int    `json:"department_id"`
	Position       string `json:"position"`
	HireDate       string `json:"hire_date"`
	EmployeeNumber string `json:"employee_number"`
	LargeText      string `json:"large_text"`
}

// Department struct definition
type Department struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	ParentID sql.NullInt64  `json:"parent_id"`
}

type DepartmentRequest struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	ParentID *int    `json:"parent_id"`  // Using pointer to handle null possibility
}

func initDB() {
	var err error
	// Get database connection information from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "rootpassword")
	dbName := getEnv("DB_NAME", "mydatabase")
	
	// Convert string to number
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

	// Create database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// Retry logic
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

// Get environment variable (with default value)
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

	// Add CORS middleware
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

	// Add health check endpoint
	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Database connection failed"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Set up API routes
	api := r.Group("/api")
	{
		// Department related APIs
		api.GET("/departments", getDepartments)
		api.GET("/departments/:id", getDepartment)
		api.GET("/departments/:id/employees", getDepartmentEmployees)
		api.POST("/departments", createDepartment)

		// Employee related APIs
		api.GET("/employees", getEmployees)
		api.GET("/employees/:id", getEmployee)

		// Add new endpoint
		api.POST("/employees/by-departments", GetEmployeesByDepartmentIDs)

		// Department tree query API
		api.GET("/departments/tree-recursive", getDepartmentTree)
		api.GET("/departments/tree-comparison", getDepartmentTreeByComparison)
	}

	r.Run(":8080")
}

func getDepartmentTreeByComparison(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID "+ id + " " + err.Error()})
		return
	}
	if idInt < 10 {
		c.JSON(400, gin.H{"error": "ID must be greater than 10"})
		return
	}

	increment := 10
	for i := len(id)- 1; i > 0 && id[i-1] == '0'; i-- {
		increment *= 10
	}
	

	query :=`
	SELECT id, name, parent_id,
	REPEAT(CONCAT(id, name, parent_id), 200000) as virtual_column
	FROM departments
	WHERE id <= ? and id > ?
	`
	log.Printf("Executing query: %s with parameters: id=%s, idInt-increment=%d", query, id, idInt-increment)
	rows, err := db.Query(query, id, idInt-increment)
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
		var virtualColumn sql.NullString
		if err := rows.Scan(&id, &name, &parentID, &virtualColumn); err != nil {
			log.Printf("Error scanning department row: %v", err)
			c.JSON(500, gin.H{"error": "Failed to scan department row"})
			return
		}
		departments = append(departments, gin.H{
			"id":        id,
			"name":      name,
			"parent_id": parentID.Int64,
			"virtual_column": virtualColumn.String,
		})
	}

	c.JSON(200, departments)
}

func getDepartmentTree(c *gin.Context) {
	parentId := c.Query("parentId")
	if parentId == "" {
		c.JSON(400, gin.H{"error": "Parent ID is required"})
		return
	}

	query := `
		WITH RECURSIVE department_tree AS (
			-- Base case: selected parent department
			SELECT id, name, parent_id, 0 as level, CAST(id AS CHAR(100)) as path,
			REPEAT(CONCAT(id, name, parent_id), 200000) as virtual_column
			FROM departments
			WHERE id = ?
			
			UNION ALL
			
			-- Recursive case: child departments
			SELECT d.id, d.name, d.parent_id, dt.level + 1, CONCAT(dt.path, ',', d.id),
			REPEAT(CONCAT(d.id, d.name, d.parent_id), 200000) as virtual_column
			FROM departments d
			INNER JOIN department_tree dt ON d.parent_id = dt.id
		)
		SELECT DISTINCT d1.id, d1.name, d1.parent_id, d1.level, d1.virtual_column
		FROM department_tree d1
		LEFT JOIN department_tree d2 ON d1.id = d2.id AND d1.path > d2.path
		WHERE d2.id IS NULL
		ORDER BY d1.level, d1.parent_id, d1.id;
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
		var virtualColumn sql.NullString
		if err := rows.Scan(&id, &name, &parentID, &level, &virtualColumn); err != nil {
			log.Printf("Error scanning department row: %v", err)
			c.JSON(500, gin.H{"error": "Failed to scan department row"})
			return
		}
		departments = append(departments, gin.H{
			"id":        id,
			"name":      name,
			"parent_id": parentID.Int64,
			"level":     level,
			"virtual_column": virtualColumn.String,
		})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating department rows: %v", err)
		c.JSON(500, gin.H{"error": "Failed to iterate department rows"})
		return
	}

	c.JSON(200, departments)
}

// Get department list
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

// Get specific department
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

// Get employee list for specific department
func getDepartmentEmployees(c *gin.Context) {
	deptID := c.Param("id")
	
	// Query for execution plan check
	explainQuery := `
		EXPLAIN FORMAT=JSON
		WITH RECURSIVE subdepartments AS (
			-- Base department
			SELECT id, parent_id
			FROM departments
			WHERE id = ?
			
			UNION ALL
			
			-- Child departments
			SELECT d.id, d.parent_id
			FROM departments d
			INNER JOIN subdepartments sd ON d.parent_id = sd.id
		)
		SELECT e.id, e.name, e.department_id, e.position, e.hire_date, e.employee_number, e.large_text
		FROM employees e
		INNER JOIN subdepartments sd ON e.department_id = sd.id
		ORDER BY e.department_id, e.name
	`
	
	// Print execution plan
	rows, err := db.Query(explainQuery, deptID)
	if err != nil {
		log.Printf("Error explaining query: %v", err)
	} else {
		defer rows.Close()
		var explainResult string
		for rows.Next() {
			if err := rows.Scan(&explainResult); err != nil {
				log.Printf("Error scanning explain result: %v", err)
			} 
		}
	}
	
	// Execute actual query
	query := `
		WITH RECURSIVE subdepartments AS (
			-- Base department
			SELECT id, parent_id
			FROM departments
			WHERE id = ?
			
			UNION ALL
			
			-- Child departments
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
const max_id_num int = 10000

func createDepartment(c *gin.Context) {
	var req DepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dept := Department{
		ID:   req.ID,
		Name: req.Name,
	}

	if req.ParentID != nil {
		dept.ParentID = sql.NullInt64{
			Int64: int64(*req.ParentID),
			Valid: true,
		}
	} else {
		dept.ParentID = sql.NullInt64{
			Valid: false,
		}
	}

	log.Printf("dept: %v", dept)

	// Handle parent ID array
	var parentIDs []interface{}
	var newID int
	if dept.ParentID.Valid && dept.ParentID.Int64 != 0 {
		// Determine increment based on parent ID digits
		increment := int64(1)
		tempID := int(dept.ParentID.Int64)
		strID := strconv.Itoa(tempID)  
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
			calcId := dept.ParentID.Int64 + (int64(i) * increment)
			parentIDs = append(parentIDs, calcId)
		}

		// Create placeholders for IN clause
		placeholders := make([]string, len(parentIDs))
		for i := range parentIDs {
			placeholders[i] = "?"
		}

		// Execute query
		query := fmt.Sprintf("SELECT id FROM departments WHERE id IN (%s) ORDER BY id", strings.Join(placeholders, ","))
		rows, err := db.Query(query, parentIDs...)
		if err != nil {
			log.Printf("Error querying departments: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query departments: %v", err)})
			return
		}
		defer rows.Close()

		for i := 1; i < max_id_length; i++ {
			var calcId int = int(dept.ParentID.Int64 + (int64(i) * increment))
			var id int	
			if !rows.Next() {
				// If no more IDs to compare, use the last calculated ID
				newID = calcId
				break
			}
			err := rows.Scan(&id)
			if err != nil {
				log.Printf("Error scanning count: %v", err)
				c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to scan count: %v", err)})
				return
			}

			// If calculated ID and queried ID are different, it means there's a gap in IDs
			// Example: calculated ID 1100, queried ID 1100 OK
			// Example: calculated ID 1300, queried ID 1400 means 1300 is available
			log.Printf("calcId: %v, id: %v", calcId, id)
			if calcId != id {
				newID = calcId
				break
			}
		}
	}else{
		// If parent ID is 0, query maximum ID
		query := fmt.Sprintf("SELECT max(id) FROM departments")
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error querying departments: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query departments: %v", err)})
			return
		}
		defer rows.Close()
		var maxID int
		if rows.Next() {
			err := rows.Scan(&maxID)
			if err != nil {
				log.Printf("Error scanning max ID: %v", err)
				c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to scan max ID: %v", err)})
				return
			}
		}
		log.Printf("maxID: %v", maxID)
		// if maxID = 2345 -> "2345"
		stringID := strconv.Itoa(maxID)
		// if maxID = 2345 -> "1000" -> 1000
		increment, err := strconv.Atoi("1" + strings.Repeat("0", len(stringID)-1))
		// if maxID = 2345 -> "2000" -> 2000
		highestDigit, err := strconv.Atoi(string(stringID[0]) + strings.Repeat("0", len(stringID)-1))
		// if maxID = 2345 -> 2000(highestDigit) + 1000(increment) = 3000
		newID = increment + highestDigit
		if err != nil {
			log.Printf("Error converting max ID: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to convert max ID: %v", err)})
			return
		}
	}
	log.Printf("newID: %v", newID)
	if newID >= max_id_num || newID == 0 {
		c.JSON(400, gin.H{"error": "can't create department"})
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

// Get employee list
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

// Get specific employee
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

	// Create placeholders for IN clause
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