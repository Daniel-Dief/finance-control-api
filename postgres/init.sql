-- =========================
-- Table: Categories
-- =========================
CREATE TABLE "Categories" (
    "Id" SERIAL PRIMARY KEY,
    "Name" TEXT NOT NULL UNIQUE
);

-- =========================
-- Table: Areas
-- =========================
CREATE TABLE "Areas" (
    "Id" SERIAL PRIMARY KEY,
    "Name" TEXT NOT NULL UNIQUE
);

-- =========================
-- Table: MonthlyBudget
-- =========================
CREATE TABLE "MonthlyBudget" (
    "Id" SERIAL PRIMARY KEY,
    "Year" INTEGER NOT NULL,
    "Month" INTEGER NOT NULL CHECK ("Month" BETWEEN 1 AND 12),
    "AreaId" INTEGER NOT NULL,
    "Amount" NUMERIC NOT NULL CHECK ("Amount" >= 0),

    CONSTRAINT fk_area
        FOREIGN KEY ("AreaId")
        REFERENCES "Areas"("Id")
        ON DELETE CASCADE,

    CONSTRAINT uq_budget UNIQUE ("Year", "Month", "AreaId")
);

-- =========================
-- Table: Transactions
-- =========================
CREATE TABLE "Transactions" (
    "Id" SERIAL PRIMARY KEY,
    "Date" DATE NOT NULL,
    "Amount" NUMERIC NOT NULL,
    "CategoryId" INTEGER,
    "AreaId" INTEGER NOT NULL,
    "Type" TEXT NOT NULL CHECK ("Type" IN ('income', 'expense')),

    CONSTRAINT fk_category
        FOREIGN KEY ("CategoryId")
        REFERENCES "Categories"("Id")
        ON DELETE SET NULL,

    CONSTRAINT fk_area_tx
        FOREIGN KEY ("AreaId")
        REFERENCES "Areas"("Id")
        ON DELETE CASCADE
);

-- =========================
-- Indexes
-- =========================

CREATE INDEX idx_transactions_date ON "Transactions"("Date");
CREATE INDEX idx_transactions_type ON "Transactions"("Type");
CREATE INDEX idx_transactions_category ON "Transactions"("CategoryId");
CREATE INDEX idx_transactions_area ON "Transactions"("AreaId");

CREATE INDEX idx_budget_area ON "MonthlyBudget"("AreaId");
CREATE INDEX idx_budget_period ON "MonthlyBudget"("Year", "Month");