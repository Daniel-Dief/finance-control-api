-- Enable foreign keys (supported in Cloudflare D1)
PRAGMA foreign_keys = ON;

-- =========================
-- Table: Categories
-- =========================
CREATE TABLE Categories (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
);

-- =========================
-- Table: Areas
-- =========================
CREATE TABLE Areas (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
);

-- =========================
-- Table: MonthlyBudget
-- =========================
CREATE TABLE MonthlyBudget (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Year INTEGER NOT NULL,
    Month INTEGER NOT NULL CHECK (Month BETWEEN 1 AND 12),
    AreaId INTEGER NOT NULL,
    Amount REAL NOT NULL CHECK (Amount >= 0),

    FOREIGN KEY (AreaId) REFERENCES Areas(Id) ON DELETE CASCADE,

    -- Prevent duplicate budget for same area/month/year
    UNIQUE (Year, Month, AreaId)
);

-- =========================
-- Table: Transactions
-- =========================
CREATE TABLE Transactions (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Date TEXT NOT NULL, -- ISO8601: YYYY-MM-DD
    Amount REAL NOT NULL,
    CategoryId INTEGER,
    Type TEXT NOT NULL CHECK (Type IN ('income', 'expense')),

    FOREIGN KEY (CategoryId) REFERENCES Categories(Id) ON DELETE SET NULL
);

-- =========================
-- Indexes
-- =========================

CREATE INDEX idx_transactions_date ON Transactions(Date);
CREATE INDEX idx_transactions_type ON Transactions(Type);
CREATE INDEX idx_transactions_category ON Transactions(CategoryId);

CREATE INDEX idx_budget_area ON MonthlyBudget(AreaId);
CREATE INDEX idx_budget_period ON MonthlyBudget(Year, Month);