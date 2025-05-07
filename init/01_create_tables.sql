-- Create tables
CREATE TABLE IF NOT EXISTS departments (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_id INT,
    FOREIGN KEY (parent_id) REFERENCES departments(id)
);

CREATE TABLE IF NOT EXISTS employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    employee_number VARCHAR(10) NOT NULL UNIQUE DEFAULT '',
    name VARCHAR(100) NOT NULL,
    position VARCHAR(50) NOT NULL,
    department_id INT NOT NULL,
    hire_date DATE NOT NULL,
    FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- Create top-level departments (1000~10000)
INSERT INTO departments (id, name, parent_id) VALUES
(1000, 'Management Support Division', NULL),
(2000, 'Sales Division', NULL),
(3000, 'Production Division', NULL),
(4000, 'Research & Development Division', NULL),
(5000, 'IT Division', NULL),
(6000, 'HR & General Affairs Division', NULL),
(7000, 'Finance & Accounting Division', NULL),
(8000, 'Quality Control Division', NULL),
(9000, 'International Business Division', NULL),
(10000, 'Strategic Planning Division', NULL);

-- Create first-level departments (3 for each division)
INSERT INTO departments (id, name, parent_id) VALUES
-- Management Support Division
(900, 'General Affairs Team', 1000),
(800, 'Legal Team', 1000),
(700, 'Planning Team', 1000),

-- Sales Division
(1900, 'Domestic Sales Team', 2000),
(1800, 'International Sales Team', 2000),
(1700, 'Sales Support Team', 2000),

-- Production Division
(2900, 'Production Management Team', 3000),
(2800, 'Facility Management Team', 3000),
(2700, 'Material Management Team', 3000),

-- R&D Division
(3900, 'Product Development Team', 4000),
(3800, 'Technical Research Team', 4000),
(3700, 'Quality Improvement Team', 4000),

-- IT Division
(4900, 'System Development Team', 5000),
(4800, 'Infrastructure Team', 5000),
(4700, 'Security Team', 5000),

-- HR & General Affairs Division
(5900, 'HR Team', 6000),
(5800, 'Training Team', 6000),
(5700, 'Benefits Team', 6000),

-- Finance & Accounting Division
(6900, 'Finance Team', 7000),
(6800, 'Accounting Team', 7000),
(6700, 'Tax Team', 7000),

-- Quality Control Division
(7900, 'Quality Assurance Team', 8000),
(7800, 'Inspection Team', 8000),
(7700, 'Environment & Safety Team', 8000),

-- International Business Division
(8900, 'Asia Region Team', 9000),
(8800, 'Europe Region Team', 9000),
(8700, 'Americas Region Team', 9000),

-- Strategic Planning Division
(9900, 'Strategic Planning Team', 10000),
(9800, 'Business Development Team', 10000),
(9700, 'Investment Management Team', 10000);

-- Create second-level departments (2-3 for each first-level department)
INSERT INTO departments (id, name, parent_id) VALUES
-- Management Support Division second-level
(910, 'General Affairs 1 Team', 900),
(920, 'Office Management Team', 900),
(810, 'Contract Management Team', 800),
(820, 'Legal Support Team', 800),
(710, 'Management Planning Team', 700),
(720, 'Performance Management Team', 700),

-- Sales Division second-level
(1910, 'Domestic Sales 1 Team', 1900),
(1920, 'Domestic Sales 2 Team', 1900),
(1810, 'Asia Sales Team', 1800),
(1820, 'Europe Sales Team', 1800),
(1710, 'Sales Support 1 Team', 1700),
(1720, 'Sales Support 2 Team', 1700),

-- Production Division second-level
(2910, 'Production Planning Team', 2900),
(2920, 'Production Control Team', 2900),
(2810, 'Facility Maintenance Team', 2800),
(2820, 'Facility Improvement Team', 2800),
(2710, 'Material Procurement Team', 2700),
(2720, 'Material Control Team', 2700),

-- R&D Division second-level
(3910, 'Product Development 1 Team', 3900),
(3920, 'Product Development 2 Team', 3900),
(3810, 'Technical Research 1 Team', 3800),
(3820, 'Technical Research 2 Team', 3800),
(3710, 'Quality Improvement 1 Team', 3700),
(3720, 'Quality Improvement 2 Team', 3700),

-- IT Division second-level
(4910, 'Web Development Team', 4900),
(4920, 'Mobile Development Team', 4900),
(4810, 'Server Management Team', 4800),
(4820, 'Network Team', 4800),
(4710, 'Security Management Team', 4700),
(4720, 'Security Audit Team', 4700),

-- HR & General Affairs Division second-level
(5910, 'HR Management Team', 5900),
(5920, 'Recruitment Team', 5900),
(5810, 'Training Planning Team', 5800),
(5820, 'Training Operations Team', 5800),
(5710, 'Benefits 1 Team', 5700),
(5720, 'Benefits 2 Team', 5700),

-- Finance & Accounting Division second-level
(6910, 'Financial Planning Team', 6900),
(6920, 'Treasury Team', 6800),
(6810, 'Accounting Management Team', 6800),
(6820, 'Financial Accounting Team', 6800),
(6710, 'Tax Management Team', 6700),
(6720, 'Tax Planning Team', 6700),

-- Quality Control Division second-level
(7910, 'Quality Assurance 1 Team', 7900),
(7920, 'Quality Assurance 2 Team', 7900),
(7810, 'Inspection Management Team', 7800),
(7820, 'Inspection Operations Team', 7800),
(7710, 'Environment Management Team', 7700),
(7720, 'Safety Management Team', 7700),

-- International Business Division second-level
(8910, 'East Asia Team', 8900),
(8920, 'Southeast Asia Team', 8900),
(8810, 'Western Europe Team', 8800),
(8820, 'Eastern Europe Team', 8800),
(8710, 'North America Team', 8700),
(8720, 'South America Team', 8700),

-- Strategic Planning Division second-level
(9910, 'Strategic Planning 1 Team', 9900),
(9920, 'Strategic Planning 2 Team', 9900),
(9810, 'Business Development 1 Team', 9800),
(9820, 'Business Development 2 Team', 9800),
(9710, 'Investment Management 1 Team', 9700),
(9720, 'Investment Management 2 Team', 9700);

-- Create employees (1000 employees distributed across departments)
INSERT INTO employees (employee_number, name, position, department_id, hire_date) VALUES
-- Management Support Division (100 employees)
('MS001', 'John Smith', 'Manager', 1000, '2020-01-01'),
('MS002', 'Sarah Johnson', 'Senior Staff', 900, '2020-02-15'),
('MS003', 'Michael Brown', 'Staff', 910, '2020-03-01'),
('MS004', 'Emily Davis', 'Staff', 920, '2020-04-15'),
('MS005', 'David Wilson', 'Senior Staff', 800, '2020-05-01'),
('MS006', 'Lisa Anderson', 'Staff', 810, '2020-06-15'),
('MS007', 'Robert Taylor', 'Staff', 820, '2020-07-01'),
('MS008', 'Jennifer Martinez', 'Senior Staff', 700, '2020-08-15'),
('MS009', 'William Thomas', 'Staff', 710, '2020-09-01'),
('MS010', 'Patricia Garcia', 'Staff', 720, '2020-10-15');

-- Add more employees to reach 1000 (using a stored procedure)
DELIMITER //
CREATE PROCEDURE AddMoreEmployees()
BEGIN
    DECLARE i INT DEFAULT 11;
    DECLARE dept_id INT;
    DECLARE emp_num VARCHAR(10);
    DECLARE emp_name VARCHAR(100);
    DECLARE emp_pos VARCHAR(50);
    DECLARE hire_date DATE;
    DECLARE dept_prefix CHAR(2);
    
    WHILE i <= 1000 DO
        -- Select a random department
        SELECT id INTO dept_id FROM departments ORDER BY RAND() LIMIT 1;
        
        -- Generate employee number based on department
        SET dept_prefix = CASE
            WHEN dept_id BETWEEN 1000 AND 1999 THEN 'MS'
            WHEN dept_id BETWEEN 2000 AND 2999 THEN 'SL'
            WHEN dept_id BETWEEN 3000 AND 3999 THEN 'PD'
            WHEN dept_id BETWEEN 4000 AND 4999 THEN 'RD'
            WHEN dept_id BETWEEN 5000 AND 5999 THEN 'IT'
            WHEN dept_id BETWEEN 6000 AND 6999 THEN 'HR'
            WHEN dept_id BETWEEN 7000 AND 7999 THEN 'FA'
            WHEN dept_id BETWEEN 8000 AND 8999 THEN 'QC'
            WHEN dept_id BETWEEN 9000 AND 9999 THEN 'IB'
            ELSE 'SP'
        END;
        
        SET emp_num = CONCAT(dept_prefix, LPAD(i, 4, '0'));
        
        -- Generate random name
        SET emp_name = CONCAT(
            CASE FLOOR(RAND() * 10)
                WHEN 0 THEN 'James'
                WHEN 1 THEN 'Mary'
                WHEN 2 THEN 'John'
                WHEN 3 THEN 'Patricia'
                WHEN 4 THEN 'Robert'
                WHEN 5 THEN 'Linda'
                WHEN 6 THEN 'Michael'
                WHEN 7 THEN 'Barbara'
                WHEN 8 THEN 'William'
                ELSE 'Elizabeth'
            END,
            ' ',
            CASE FLOOR(RAND() * 10)
                WHEN 0 THEN 'Smith'
                WHEN 1 THEN 'Johnson'
                WHEN 2 THEN 'Williams'
                WHEN 3 THEN 'Brown'
                WHEN 4 THEN 'Jones'
                WHEN 5 THEN 'Garcia'
                WHEN 6 THEN 'Miller'
                WHEN 7 THEN 'Davis'
                WHEN 8 THEN 'Rodriguez'
                ELSE 'Martinez'
            END
        );
        
        -- Generate random position
        SET emp_pos = CASE FLOOR(RAND() * 5)
            WHEN 0 THEN 'Manager'
            WHEN 1 THEN 'Senior Staff'
            WHEN 2 THEN 'Staff'
            WHEN 3 THEN 'Assistant'
            ELSE 'Intern'
        END;
        
        -- Generate random hire date between 2020-01-01 and 2024-12-31
        SET hire_date = DATE_ADD('2020-01-01', 
            INTERVAL FLOOR(RAND() * DATEDIFF('2024-12-31', '2020-01-01')) DAY);
        
        -- Insert the employee
        INSERT INTO employees (employee_number, name, position, department_id, hire_date)
        VALUES (emp_num, emp_name, emp_pos, dept_id, hire_date);
        
        SET i = i + 1;
    END WHILE;
END //
DELIMITER ;

-- Execute the stored procedure
CALL AddMoreEmployees();

-- Drop the stored procedure
DROP PROCEDURE IF EXISTS AddMoreEmployees; 