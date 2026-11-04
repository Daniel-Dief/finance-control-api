-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW."UpdatedAt" = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_status_updated_at
    BEFORE UPDATE ON "Status"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON "Users"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trg_categories_updated_at
    BEFORE UPDATE ON "Categories"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trg_areas_updated_at
    BEFORE UPDATE ON "Areas"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trg_budget_updated_at
    BEFORE UPDATE ON "MonthlyBudget"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trg_transactions_updated_at
    BEFORE UPDATE ON "Transactions"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- +goose Down

DROP TRIGGER IF EXISTS trg_transactions_updated_at ON "Transactions";
DROP TRIGGER IF EXISTS trg_budget_updated_at ON "MonthlyBudget";
DROP TRIGGER IF EXISTS trg_areas_updated_at ON "Areas";
DROP TRIGGER IF EXISTS trg_categories_updated_at ON "Categories";
DROP TRIGGER IF EXISTS trg_users_updated_at ON "Users";
DROP TRIGGER IF EXISTS trg_status_updated_at ON "Status";

DROP FUNCTION IF EXISTS update_updated_at();
