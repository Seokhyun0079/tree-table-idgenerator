-- 데이터베이스 선택
USE mydatabase;

-- 부서 트리 테이블 생성
CREATE TABLE IF NOT EXISTS departments (
    id INT PRIMARY KEY,
    parent_id INT NULL,
    name VARCHAR(50) NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES departments(id)
);

-- 부서 트리 데이터 삽입
INSERT INTO departments (id, parent_id, name) VALUES
    (1, NULL, 'Company Headquarters'),
    (2, 1, 'Business Division'),
    (3, 1, 'Technology Division'),
    (4, 2, 'Sales'),
    (5, 2, 'Marketing'),
    (6, 2, 'Finance'),
    (7, 3, 'IT Development'),
    (8, 3, 'IT Infrastructure'),
    (9, 3, 'IT Security'),
    (10, 4, 'Domestic Sales'),
    (11, 4, 'International Sales'),
    (12, 5, 'Digital Marketing'),
    (13, 5, 'Brand Marketing'),
    (14, 7, 'Frontend Development'),
    (15, 7, 'Backend Development'),
    (16, NULL, 'Research Headquarters'),
    (17, 16, 'Bio Research'),
    (18, 16, 'Chemical Research'),
    (19, 16, 'Physics Research'),
    (20, 17, 'Medical Research'),
    (21, 17, 'Agricultural Research'),
    (22, 18, 'Materials Development'),
    (23, 18, 'Chemical Analysis'),
    (24, 19, 'Quantum Research'),
    (25, 19, 'Space Research'),
    (26, 20, 'Clinical Trials'),
    (27, 20, 'Drug Development'),
    (28, 22, 'New Materials Lab'),
    (29, 24, 'Quantum Computing'),
    (30, 25, 'Satellite Development');

-- 직원 테이블 생성
CREATE TABLE IF NOT EXISTS employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    department_id INT NOT NULL,
    FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- 직원 데이터 삽입 (균등 분포)
INSERT INTO employees (name, department_id) VALUES
    -- Company Headquarters (5명)
    ('Robert Johnson', 1),     -- CEO
    ('Mary Williams', 1),      -- CFO
    ('David Brown', 1),        -- COO
    ('Patricia Davis', 1),     -- CTO
    ('Michael Wilson', 1),     -- CHRO

    -- Research Headquarters (5명)
    ('James Anderson', 16),    -- Research Director
    ('Elizabeth Taylor', 16),  -- Deputy Research Director
    ('Richard Miller', 16),    -- Research Operations Manager
    ('Susan Martinez', 16),    -- Research Strategy Manager
    ('Thomas Moore', 16),      -- Research Administration Manager

    -- Business Division (30명)
    ('John Smith', 10),        -- Domestic Sales
    ('Emma Wilson', 10),
    ('Michael Brown', 10),
    ('Sarah Davis', 11),       -- International Sales
    ('James Johnson', 11),
    ('Lisa Anderson', 11),
    ('Robert Taylor', 12),     -- Digital Marketing
    ('Jennifer Martin', 12),
    ('William White', 12),
    ('Elizabeth Clark', 13),   -- Brand Marketing
    ('David Miller', 13),
    ('Mary Wilson', 13),
    ('Richard Moore', 6),      -- Finance
    ('Patricia Lee', 6),
    ('Joseph Allen', 6),

    -- Technology Division (30명)
    ('Susan Wright', 14),      -- Frontend Development
    ('Thomas Young', 14),
    ('Margaret Hall', 14),
    ('Charles King', 15),      -- Backend Development
    ('Linda Scott', 15),
    ('Daniel Lewis', 15),
    ('Barbara Hill', 8),       -- IT Infrastructure
    ('Paul Adams', 8),
    ('Karen Baker', 8),
    ('Mark Evans', 9),         -- IT Security
    ('Sandra Morris', 9),
    ('Steven Nelson', 9),

    -- Research Division - Bio Research (15명)
    ('Betty Carter', 20),      -- Medical Research
    ('Kevin Phillips', 20),
    ('Nancy Turner', 20),
    ('Gary Campbell', 21),     -- Agricultural Research
    ('Helen Rogers', 21),
    ('Frank Butler', 21),
    ('Carol Reed', 26),        -- Clinical Trials
    ('Edward Cox', 26),
    ('Sharon Murphy', 26),
    ('George Richardson', 27), -- Drug Development
    ('Ruth Cooper', 27),
    ('Gerald Morgan', 27),

    -- Research Division - Chemical Research (10명)
    ('Dennis Ward', 22),       -- Materials Development
    ('Michelle Foster', 22),
    ('Raymond Torres', 22),
    ('Christine Long', 23),    -- Chemical Analysis
    ('Peter Gray', 23),
    ('Janet Collins', 23),
    ('Jerry Stewart', 28),     -- New Materials Lab
    ('Alice Barnes', 28),
    ('Ralph Hughes', 28),

    -- Research Division - Physics Research (15명)
    ('Julie Price', 24),       -- Quantum Research
    ('Terry Sanders', 24),
    ('Gloria Ross', 24),
    ('Carl Wood', 25),         -- Space Research
    ('Teresa Jenkins', 25),
    ('Harry Bennett', 25),
    ('Ann Coleman', 29),       -- Quantum Computing
    ('Fred Perry', 29),
    ('Joyce Powell', 29),
    ('Victor Patterson', 30),  -- Satellite Development
    ('Jean Russell', 30),
    ('Roy Howard', 30); 