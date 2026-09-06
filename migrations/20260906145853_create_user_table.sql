-- +goose Up
CREATE TABLE IF NOT EXISTS "Users" (
    "Id" SERIAL PRIMARY KEY,
    "StatusId" INTEGER NOT NULL,
    "Name" TEXT NOT NULL,
    "Login" TEXT NOT NULL UNIQUE,
    "PasswordHash" TEXT NOT NULL,
    "CreatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "UpdatedAt" TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_status
        FOREIGN KEY ("StatusId")
        REFERENCES "Status"("Id")
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS "Users";