-- +goose Up

-- =========================
-- Table: Status
-- =========================
CREATE TABLE "Status" (
    "Id" SERIAL PRIMARY KEY,
    "Name" TEXT NOT NULL UNIQUE,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed default statuses
INSERT INTO "Status" ("Name") VALUES ('active'), ('inactive'), ('deleted');

-- =========================
-- Add StatusId to entity tables
-- =========================
ALTER TABLE "Areas" ADD COLUMN "StatusId" INTEGER NOT NULL DEFAULT 1;
ALTER TABLE "Areas" ADD CONSTRAINT fk_areas_status
    FOREIGN KEY ("StatusId") REFERENCES "Status"("Id");

ALTER TABLE "MonthlyBudget" ADD COLUMN "StatusId" INTEGER NOT NULL DEFAULT 1;
ALTER TABLE "MonthlyBudget" ADD CONSTRAINT fk_budget_status
    FOREIGN KEY ("StatusId") REFERENCES "Status"("Id");

ALTER TABLE "Transactions" ADD COLUMN "StatusId" INTEGER NOT NULL DEFAULT 1;
ALTER TABLE "Transactions" ADD CONSTRAINT fk_tx_status
    FOREIGN KEY ("StatusId") REFERENCES "Status"("Id");

-- +goose Down

ALTER TABLE "Transactions" DROP CONSTRAINT IF EXISTS fk_tx_status;
ALTER TABLE "Transactions" DROP COLUMN IF EXISTS "StatusId";

ALTER TABLE "MonthlyBudget" DROP CONSTRAINT IF EXISTS fk_budget_status;
ALTER TABLE "MonthlyBudget" DROP COLUMN IF EXISTS "StatusId";

ALTER TABLE "Areas" DROP CONSTRAINT IF EXISTS fk_areas_status;
ALTER TABLE "Areas" DROP COLUMN IF EXISTS "StatusId";

DELETE FROM "Status" WHERE "Name" IN ('active', 'inactive', 'deleted');
DROP TABLE IF EXISTS "Status";