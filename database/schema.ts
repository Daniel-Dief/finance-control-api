import {
  sqliteTable,
  integer,
  text,
  real,
  uniqueIndex,
  index,
  check,
} from "drizzle-orm/sqlite-core";
import { sql } from "drizzle-orm";

// =========================
// Table: Categories
// =========================
export const categories = sqliteTable(
  "Categories",
  {
    id: integer("Id").primaryKey({ autoIncrement: true }),
    name: text("Name").notNull().unique(),
  }
);

// =========================
// Table: Areas
// =========================
export const areas = sqliteTable(
  "Areas",
  {
    id: integer("Id").primaryKey({ autoIncrement: true }),
    name: text("Name").notNull().unique(),
  }
);

// =========================
// Table: MonthlyBudget
// =========================
export const monthlyBudget = sqliteTable(
  "MonthlyBudget",
  {
    id: integer("Id").primaryKey({ autoIncrement: true }),
    year: integer("Year").notNull(),
    month: integer("Month").notNull(),
    areaId: integer("AreaId")
      .notNull()
      .references(() => areas.id, { onDelete: "cascade" }),
    amount: real("Amount").notNull(),
  },
  (table) => ({
    // UNIQUE (Year, Month, AreaId)
    uniquePeriodArea: uniqueIndex("uq_budget_year_month_area").on(
      table.year,
      table.month,
      table.areaId
    ),

    // CHECK (Month BETWEEN 1 AND 12)
    checkMonth: check(
      "chk_month_range",
      sql`${table.month} BETWEEN 1 AND 12`
    ),

    // CHECK (Amount >= 0)
    checkAmount: check(
      "chk_amount_positive",
      sql`${table.amount} >= 0`
    ),

    // Indexes
    idxArea: index("idx_budget_area").on(table.areaId),
    idxPeriod: index("idx_budget_period").on(table.year, table.month),
  })
);

// =========================
// Table: Transactions
// =========================
export const transactions = sqliteTable(
  "Transactions",
  {
    id: integer("Id").primaryKey({ autoIncrement: true }),
    date: text("Date").notNull(), // ISO8601
    amount: real("Amount").notNull(),
    categoryId: integer("CategoryId").references(() => categories.id, {
      onDelete: "set null",
    }),
    areaId: integer("AreaId")
      .notNull()
      .references(() => areas.id, { onDelete: "cascade" }),
    type: text("Type").notNull(), // 'income' | 'expense'
  },
  (table) => ({
    // CHECK (Type IN ('income', 'expense'))
    checkType: check(
      "chk_transaction_type",
      sql`${table.type} IN ('income', 'expense')`
    ),

    // Indexes
    idxDate: index("idx_transactions_date").on(table.date),
    idxType: index("idx_transactions_type").on(table.type),
    idxCategory: index("idx_transactions_category").on(table.categoryId),
    idxArea: index("idx_transactions_area").on(table.areaId),
  })
);