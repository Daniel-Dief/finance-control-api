-- +goose Up

ALTER TABLE "Categories" ADD COLUMN "StatusId" INTEGER NOT NULL DEFAULT 1;
ALTER TABLE "Categories" ADD CONSTRAINT fk_categories_status
    FOREIGN KEY ("StatusId") REFERENCES "Status"("Id");

-- +goose Down

ALTER TABLE "Categories" DROP CONSTRAINT IF EXISTS fk_categories_status;
ALTER TABLE "Categories" DROP COLUMN IF EXISTS "StatusId";
