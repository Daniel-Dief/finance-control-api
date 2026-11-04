-- +goose Up

ALTER TABLE "Areas" ADD COLUMN "UserId" INTEGER NOT NULL;
ALTER TABLE "Areas" ADD CONSTRAINT fk_areas_user
    FOREIGN KEY ("UserId") REFERENCES "Users"("Id");

ALTER TABLE "Categories" ADD COLUMN "UserId" INTEGER NOT NULL;
ALTER TABLE "Categories" ADD CONSTRAINT fk_categories_user
    FOREIGN KEY ("UserId") REFERENCES "Users"("Id");

ALTER TABLE "MonthlyBudget" ADD COLUMN "UserId" INTEGER NOT NULL;
ALTER TABLE "MonthlyBudget" ADD CONSTRAINT fk_budget_user
    FOREIGN KEY ("UserId") REFERENCES "Users"("Id");

ALTER TABLE "Transactions" ADD COLUMN "UserId" INTEGER NOT NULL;
ALTER TABLE "Transactions" ADD CONSTRAINT fk_tx_user
    FOREIGN KEY ("UserId") REFERENCES "Users"("Id");

CREATE INDEX idx_areas_user ON "Areas"("UserId");
CREATE INDEX idx_categories_user ON "Categories"("UserId");
CREATE INDEX idx_budget_user ON "MonthlyBudget"("UserId");
CREATE INDEX idx_transactions_user ON "Transactions"("UserId");

-- +goose Down

DROP INDEX IF EXISTS idx_transactions_user;
DROP INDEX IF EXISTS idx_budget_user;
DROP INDEX IF EXISTS idx_categories_user;
DROP INDEX IF EXISTS idx_areas_user;

ALTER TABLE "Transactions" DROP CONSTRAINT IF EXISTS fk_tx_user;
ALTER TABLE "Transactions" DROP COLUMN IF EXISTS "UserId";

ALTER TABLE "MonthlyBudget" DROP CONSTRAINT IF EXISTS fk_budget_user;
ALTER TABLE "MonthlyBudget" DROP COLUMN IF EXISTS "UserId";

ALTER TABLE "Categories" DROP CONSTRAINT IF EXISTS fk_categories_user;
ALTER TABLE "Categories" DROP COLUMN IF EXISTS "UserId";

ALTER TABLE "Areas" DROP CONSTRAINT IF EXISTS fk_areas_user;
ALTER TABLE "Areas" DROP COLUMN IF EXISTS "UserId";
