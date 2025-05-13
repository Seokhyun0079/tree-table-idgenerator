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
    large_text LONGTEXT,
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
(7000, 'Finance & Accounting Division', NULL);

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
(6700, 'Tax Team', 7000);

-- Create second-level departments (2-3 for each first-level department)
INSERT INTO departments (id, name, parent_id) VALUES
-- Management Support Division second-level
(890, 'General Affairs 1 Team', 900),
(880, 'Office Management Team', 900),
(790, 'Contract Management Team', 800),
(780, 'Legal Support Team', 800),
(690, 'Management Planning Team', 700),
(680, 'Performance Management Team', 700),

-- Sales Division second-level
(1890, 'Domestic Sales 1 Team', 1900),
(1880, 'Domestic Sales 2 Team', 1900),
(1790, 'Asia Sales Team', 1800),
(1780, 'Europe Sales Team', 1800),
(1690, 'Sales Support 1 Team', 1700),
(1680, 'Sales Support 2 Team', 1700),

-- Production Division second-level
(2890, 'Production Planning Team', 2900),
(2880, 'Production Control Team', 2900),
(2790, 'Facility Maintenance Team', 2800),
(2780, 'Facility Improvement Team', 2800),
(2690, 'Material Procurement Team', 2700),
(2680, 'Material Control Team', 2700),

-- R&D Division second-level
(3890, 'Product Development 1 Team', 3900),
(3880, 'Product Development 2 Team', 3900),
(3790, 'Technical Research 1 Team', 3800),
(3780, 'Technical Research 2 Team', 3800),
(3690, 'Quality Improvement 1 Team', 3700),
(3680, 'Quality Improvement 2 Team', 3700),

-- IT Division second-level
(4890, 'Web Development Team', 4900),
(4880, 'Mobile Development Team', 4900),
(4790, 'Server Management Team', 4800),
(4780, 'Network Team', 4800),
(4690, 'Security Management Team', 4700),
(4680, 'Security Audit Team', 4700),

-- HR & General Affairs Division second-level
(5890, 'HR Management Team', 5900),
(5880, 'Recruitment Team', 5900),
(5790, 'Training Planning Team', 5800),
(5780, 'Training Operations Team', 5800),
(5690, 'Benefits 1 Team', 5700),
(5680, 'Benefits 2 Team', 5700),

-- Finance & Accounting Division second-level
(6890, 'Financial Planning Team', 6900),
(6880, 'Treasury Team', 6900),
(6790, 'Accounting Management Team', 6800),
(6780, 'Financial Accounting Team', 6800),
(6690, 'Tax Management Team', 6700),
(6680, 'Tax Planning Team', 6700),

-- Quality Control Division second-level
(7890, 'Quality Assurance 1 Team', 7900),
(7880, 'Quality Assurance 2 Team', 7900),
(7790, 'Inspection Management Team', 7800),
(7780, 'Inspection Operations Team', 7800),
(7690, 'Environment Management Team', 7700),
(7680, 'Safety Management Team', 7700),

-- International Business Division second-level
(8890, 'East Asia Team', 8900),
(8880, 'Southeast Asia Team', 8900),
(8790, 'Western Europe Team', 8800),
(8780, 'Eastern Europe Team', 8800),
(8690, 'North America Team', 8700),
(8680, 'South America Team', 8700),

-- Strategic Planning Division second-level
(9890, 'Strategic Planning 1 Team', 9900),
(9880, 'Business Development Team', 9900),
(9790, 'Investment Management Team', 9900);

-- Create third-level departments (2-3 for each second-level department)
INSERT INTO departments (id, name, parent_id) VALUES
-- Management Support Division third-level
(889, 'General Affairs 1-1 Team', 890),
(888, 'General Affairs 1-2 Team', 890),
(879, 'Office Management 1 Team', 880),
(878, 'Office Management 2 Team', 880),
(789, 'Contract Management 1 Team', 790),
(788, 'Contract Management 2 Team', 790),
(779, 'Legal Support 1 Team', 780),
(778, 'Legal Support 2 Team', 780),
(689, 'Management Planning 1 Team', 690),
(688, 'Management Planning 2 Team', 690),
(679, 'Performance Management 1 Team', 680),
(678, 'Performance Management 2 Team', 680),

-- Sales Division third-level
(1889, 'Domestic Sales 1-1 Team', 1890),
(1888, 'Domestic Sales 1-2 Team', 1890),
(1879, 'Domestic Sales 2-1 Team', 1880),
(1878, 'Domestic Sales 2-2 Team', 1880),
(1789, 'Asia Sales 1 Team', 1790),
(1788, 'Asia Sales 2 Team', 1790),
(1779, 'Europe Sales 1 Team', 1780),
(1778, 'Europe Sales 2 Team', 1780),
(1689, 'Sales Support 1-1 Team', 1690),
(1688, 'Sales Support 1-2 Team', 1690),
(1679, 'Sales Support 2-1 Team', 1680),
(1678, 'Sales Support 2-2 Team', 1680),

-- Production Division third-level
(2889, 'Production Planning 1 Team', 2890),
(2888, 'Production Planning 2 Team', 2890),
(2879, 'Production Control 1 Team', 2880),
(2878, 'Production Control 2 Team', 2880),
(2789, 'Facility Maintenance 1 Team', 2790),
(2788, 'Facility Maintenance 2 Team', 2790),
(2779, 'Facility Improvement 1 Team', 2780),
(2778, 'Facility Improvement 2 Team', 2780),
(2689, 'Material Procurement 1 Team', 2690),
(2688, 'Material Procurement 2 Team', 2690),
(2679, 'Material Control 1 Team', 2680),
(2678, 'Material Control 2 Team', 2680),

-- R&D Division third-level
(3889, 'Product Development 1-1 Team', 3890),
(3888, 'Product Development 1-2 Team', 3890),
(3879, 'Product Development 2-1 Team', 3880),
(3878, 'Product Development 2-2 Team', 3880),
(3789, 'Technical Research 1-1 Team', 3790),
(3788, 'Technical Research 1-2 Team', 3790),
(3779, 'Technical Research 2-1 Team', 3780),
(3778, 'Technical Research 2-2 Team', 3780),
(3689, 'Quality Improvement 1-1 Team', 3690),
(3688, 'Quality Improvement 1-2 Team', 3690),
(3679, 'Quality Improvement 2-1 Team', 3680),
(3678, 'Quality Improvement 2-2 Team', 3680),

-- IT Division third-level
(4889, 'Web Development 1 Team', 4890),
(4888, 'Web Development 2 Team', 4890),
(4879, 'Mobile Development 1 Team', 4880),
(4878, 'Mobile Development 2 Team', 4880),
(4789, 'Server Management 1 Team', 4790),
(4788, 'Server Management 2 Team', 4790),
(4779, 'Network 1 Team', 4780),
(4778, 'Network 2 Team', 4780),
(4689, 'Security Management 1 Team', 4690),
(4688, 'Security Management 2 Team', 4690),
(4679, 'Security Audit 1 Team', 4680),
(4678, 'Security Audit 2 Team', 4680),

-- HR & General Affairs Division third-level
(5889, 'HR Management 1 Team', 5890),
(5888, 'HR Management 2 Team', 5890),
(5879, 'Recruitment 1 Team', 5880),
(5878, 'Recruitment 2 Team', 5880),
(5789, 'Training Planning 1 Team', 5790),
(5788, 'Training Planning 2 Team', 5790),
(5779, 'Training Operations 1 Team', 5780),
(5778, 'Training Operations 2 Team', 5780),
(5689, 'Benefits 1-1 Team', 5690),
(5688, 'Benefits 1-2 Team', 5690),
(5679, 'Benefits 2-1 Team', 5680),
(5678, 'Benefits 2-2 Team', 5680),

-- Finance & Accounting Division third-level
(6889, 'Financial Planning 1 Team', 6890),
(6888, 'Financial Planning 2 Team', 6890),
(6879, 'Treasury 1 Team', 6880),
(6878, 'Treasury 2 Team', 6880),
(6789, 'Accounting Management 1 Team', 6790),
(6788, 'Accounting Management 2 Team', 6790),
(6779, 'Financial Accounting 1 Team', 6780),
(6778, 'Financial Accounting 2 Team', 6780),
(6689, 'Tax Management 1 Team', 6690),
(6688, 'Tax Management 2 Team', 6690),
(6679, 'Tax Planning 1 Team', 6680),
(6678, 'Tax Planning 2 Team', 6680),


-- Add more employees to reach 100000 (using a stored procedure)
DELIMITER //
CREATE PROCEDURE AddMoreEmployees()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE dept_id INT;
    DECLARE dept_name VARCHAR(100);
    DECLARE dept_prefix VARCHAR(10);
    DECLARE i INT;
    
    -- Set prefix for each department
    DECLARE dept_prefixes CURSOR FOR 
        SELECT id, name,
            CASE 
                WHEN id = 1000 THEN 'MS'  -- Management Support
                WHEN id = 2000 THEN 'SD'  -- Sales Division
                WHEN id = 3000 THEN 'PD'  -- Production Division
                WHEN id = 4000 THEN 'RD'  -- R&D Division
                WHEN id = 5000 THEN 'IT'  -- IT Division
                WHEN id = 6000 THEN 'HR'  -- HR Division
                WHEN id = 7000 THEN 'FA'  -- Finance & Accounting
                ELSE CONCAT('D', LPAD(SUBSTRING(id, 1, 2), 2, '0'))
            END as prefix
        FROM departments;
    
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
    OPEN dept_prefixes;
    
    read_loop: LOOP
        FETCH dept_prefixes INTO dept_id, dept_name, dept_prefix;
        IF done THEN
            LEAVE read_loop;
        END IF;
        
        -- Add 100 employees to each department
        SET i = 1;
        WHILE i <= 100 DO
            INSERT INTO employees (
                name, 
                department_id, 
                position, 
                hire_date, 
                employee_number,
                large_text
            ) VALUES (
                CONCAT(dept_name, ' Employee ', i),
                dept_id,
                CASE 
                    WHEN i <= 5 THEN 'Manager'
                    WHEN i <= 15 THEN 'Deputy Manager'
                    WHEN i <= 30 THEN 'Team Leader'
                    WHEN i <= 50 THEN 'Senior Staff'
                    ELSE 'Staff'
                END,
                DATE_ADD('2020-01-01', INTERVAL FLOOR(RAND() * 1000) DAY),
                CONCAT(dept_prefix, LPAD(SUBSTRING(dept_id, -2), 2, '0'), LPAD(i, 3, '0')),
                REPEAT(CONCAT('Employee Number: ', CONCAT(dept_prefix, LPAD(SUBSTRING(dept_id, -2), 2, '0'), LPAD(i, 3, '0')), 
                       ', Department: ', dept_name, 
                       ', Position: ', 
                       CASE 
                           WHEN i <= 5 THEN 'Manager'
                           WHEN i <= 15 THEN 'Deputy Manager'
                           WHEN i <= 30 THEN 'Team Leader'
                           WHEN i <= 50 THEN 'Senior Staff'
                           ELSE 'Staff'
                       END, 
                       ', Hire Date: ', 
                       DATE_FORMAT(DATE_ADD('2020-01-01', INTERVAL FLOOR(RAND() * 1000) DAY), '%Y-%m-%d')), 1000)
            );
            SET i = i + 1;
        END WHILE;
    END LOOP;
    
    CLOSE dept_prefixes;
END //
DELIMITER ;

-- Execute stored procedure
CALL AddMoreEmployees();
DROP PROCEDURE IF EXISTS AddMoreEmployees;

-- Verify the data
SELECT 
    d.id as dept_id,
    d.name as dept_name,
    COUNT(e.id) as employee_count
FROM departments d
LEFT JOIN employees e ON d.id = e.department_id
GROUP BY d.id, d.name
ORDER BY d.id; 