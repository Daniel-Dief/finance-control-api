-- +goose Up

-- Nomes eram únicos globalmente (UNIQUE na coluna). Agora devem ser únicos
-- por usuário. Usamos índices parciais (só ativos) para que nomes de registros
-- soft-deleted (StatusId = 3) possam ser reutilizados.

ALTER TABLE "Areas" DROP CONSTRAINT IF EXISTS "Areas_Name_key";
ALTER TABLE "Categories" DROP CONSTRAINT IF EXISTS "Categories_Name_key";
ALTER TABLE "MonthlyBudget" DROP CONSTRAINT IF EXISTS uq_budget;

CREATE UNIQUE INDEX uq_areas_name_user ON "Areas"("Name", "UserId") WHERE "StatusId" = 1;
CREATE UNIQUE INDEX uq_categories_name_user ON "Categories"("Name", "UserId") WHERE "StatusId" = 1;
CREATE UNIQUE INDEX uq_budget_period_user ON "MonthlyBudget"("Year", "Month", "AreaId", "UserId") WHERE "StatusId" = 1;

-- +goose Down

DROP INDEX IF EXISTS uq_budget_period_user;
DROP INDEX IF EXISTS uq_categories_name_user;
DROP INDEX IF EXISTS uq_areas_name_user;

ALTER TABLE "Categories" ADD CONSTRAINT "Categories_Name_key" UNIQUE ("Name");
ALTER TABLE "Areas" ADD CONSTRAINT "Areas_Name_key" UNIQUE ("Name");
ALTER TABLE "MonthlyBudget" ADD CONSTRAINT uq_budget UNIQUE ("Year", "Month", "AreaId");