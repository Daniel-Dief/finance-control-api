-- =========================
-- Table: Status
-- =========================
CREATE TABLE "Status" (
    "Id" SERIAL PRIMARY KEY,
    "Name" TEXT NOT NULL UNIQUE,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO "Status" ("Name") VALUES ('active'), ('inactive'), ('deleted');

-- =========================
-- Table: Users
-- =========================
CREATE TABLE "Users" (
    "Id" SERIAL PRIMARY KEY,
    "StatusId" INTEGER NOT NULL,
    "Name" TEXT NOT NULL,
    "Login" TEXT NOT NULL UNIQUE,
    "PasswordHash" TEXT NOT NULL,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_users_status
        FOREIGN KEY ("StatusId")
        REFERENCES "Status"("Id")
        ON DELETE CASCADE
);

-- =========================
-- Table: Categories
-- =========================
CREATE TABLE "Categories" (
    "Id" SERIAL PRIMARY KEY,
    "Name" TEXT NOT NULL,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "StatusId" INTEGER NOT NULL DEFAULT 1,
    "UserId" INTEGER NOT NULL,
    CONSTRAINT fk_categories_status
        FOREIGN KEY ("StatusId")
        REFERENCES "Status"("Id"),
    CONSTRAINT fk_categories_user
        FOREIGN KEY ("UserId")
        REFERENCES "Users"("Id")
);

-- =========================
-- Table: Areas
-- =========================
CREATE TABLE "Areas" (
    "Id" SERIAL PRIMARY KEY,
    "Name" TEXT NOT NULL,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "StatusId" INTEGER NOT NULL DEFAULT 1,
    "UserId" INTEGER NOT NULL,
    CONSTRAINT fk_areas_status
        FOREIGN KEY ("StatusId")
        REFERENCES "Status"("Id"),
    CONSTRAINT fk_areas_user
        FOREIGN KEY ("UserId")
        REFERENCES "Users"("Id")
);

-- =========================
-- Table: MonthlyBudget
-- =========================
CREATE TABLE "MonthlyBudget" (
    "Id" SERIAL PRIMARY KEY,
    "Year" INTEGER NOT NULL,
    "Month" INTEGER NOT NULL CHECK ("Month" BETWEEN 1 AND 12),
    "AreaId" INTEGER NOT NULL,
    "Amount" INTEGER NOT NULL CHECK ("Amount" >= 0),
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "StatusId" INTEGER NOT NULL DEFAULT 1,
    "UserId" INTEGER NOT NULL,

    CONSTRAINT fk_area
        FOREIGN KEY ("AreaId")
        REFERENCES "Areas"("Id")
        ON DELETE CASCADE,

    CONSTRAINT fk_budget_status
        FOREIGN KEY ("StatusId")
        REFERENCES "Status"("Id"),

    CONSTRAINT fk_budget_user
        FOREIGN KEY ("UserId")
        REFERENCES "Users"("Id")
);

-- =========================
-- Table: Transactions
-- =========================
CREATE TABLE "Transactions" (
    "Id" SERIAL PRIMARY KEY,
    "Date" DATE NOT NULL,
    "Amount" INTEGER NOT NULL,
    "CategoryId" INTEGER,
    "AreaId" INTEGER NOT NULL,
    "Type" TEXT NOT NULL CHECK ("Type" IN ('income', 'expense')),
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "StatusId" INTEGER NOT NULL DEFAULT 1,
    "UserId" INTEGER NOT NULL,

    CONSTRAINT fk_category
        FOREIGN KEY ("CategoryId")
        REFERENCES "Categories"("Id")
        ON DELETE SET NULL,

    CONSTRAINT fk_area_tx
        FOREIGN KEY ("AreaId")
        REFERENCES "Areas"("Id")
        ON DELETE CASCADE,

    CONSTRAINT fk_tx_status
        FOREIGN KEY ("StatusId")
        REFERENCES "Status"("Id"),

    CONSTRAINT fk_tx_user
        FOREIGN KEY ("UserId")
        REFERENCES "Users"("Id")
);

-- =========================
-- Indexes
-- =========================

CREATE INDEX idx_transactions_date ON "Transactions"("Date");
CREATE INDEX idx_transactions_type ON "Transactions"("Type");
CREATE INDEX idx_transactions_category ON "Transactions"("CategoryId");
CREATE INDEX idx_transactions_area ON "Transactions"("AreaId");
CREATE INDEX idx_transactions_user ON "Transactions"("UserId");

CREATE INDEX idx_budget_area ON "MonthlyBudget"("AreaId");
CREATE INDEX idx_budget_period ON "MonthlyBudget"("Year", "Month");
CREATE INDEX idx_budget_user ON "MonthlyBudget"("UserId");

CREATE INDEX idx_areas_user ON "Areas"("UserId");
CREATE INDEX idx_categories_user ON "Categories"("UserId");

-- Unicidade por usuário, considerando apenas registros ativos:
-- nomes e orçamentos podem ser reutilizados após soft delete.
CREATE UNIQUE INDEX uq_areas_name_user ON "Areas"("Name", "UserId") WHERE "StatusId" = 1;
CREATE UNIQUE INDEX uq_categories_name_user ON "Categories"("Name", "UserId") WHERE "StatusId" = 1;
CREATE UNIQUE INDEX uq_budget_period_user ON "MonthlyBudget"("Year", "Month", "AreaId", "UserId") WHERE "StatusId" = 1;
