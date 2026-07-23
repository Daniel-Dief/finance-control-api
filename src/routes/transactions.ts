import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { eq, and, gte, lte } from "drizzle-orm";
import { createDb, type Env } from "../db";
import { transactions } from "../../database/schema";
import { ErrorSchema } from "./schemas";

const app = new OpenAPIHono<{ Bindings: Env }>();

const TransactionSchema = z
  .object({
    id: z.number().int().openapi({ example: 1 }),
    date: z.string().openapi({ example: "2026-07-15" }),
    amount: z.number().openapi({ example: 250.0 }),
    categoryId: z.number().int().nullable().openapi({ example: 1 }),
    areaId: z.number().int().openapi({ example: 1 }),
    type: z.enum(["income", "expense"]).openapi({ example: "expense" }),
  })
  .openapi("Transaction");

const CreateTransactionSchema = z
  .object({
    date: z
      .string()
      .regex(/^\d{4}-\d{2}-\d{2}$/)
      .openapi({ example: "2026-07-15" }),
    amount: z.number().openapi({ example: 250.0 }),
    categoryId: z.number().int().optional().openapi({ example: 1 }),
    areaId: z.number().int().openapi({ example: 1 }),
    type: z.enum(["income", "expense"]).openapi({ example: "expense" }),
  })
  .openapi("CreateTransaction");

const UpdateTransactionSchema = z
  .object({
    date: z
      .string()
      .regex(/^\d{4}-\d{2}-\d{2}$/)
      .optional()
      .openapi({ example: "2026-07-15" }),
    amount: z.number().optional().openapi({ example: 250.0 }),
    categoryId: z.number().int().nullable().optional().openapi({ example: 1 }),
    areaId: z.number().int().optional().openapi({ example: 1 }),
    type: z.enum(["income", "expense"]).optional().openapi({ example: "expense" }),
  })
  .openapi("UpdateTransaction");

const listTransactionsRoute = createRoute({
  method: "get",
  path: "/",
  tags: ["Transactions"],
  summary: "Listar transações",
  description:
    "Retorna todas as transações. Filtre por ?type=income|expense, ?categoryId=, ?areaId=, ?from= (data mínima), ?to= (data máxima).",
  request: {
    query: z.object({
      type: z.enum(["income", "expense"]).optional().openapi({ example: "expense" }),
      categoryId: z.coerce.number().int().optional().openapi({ example: 1 }),
      areaId: z.coerce.number().int().optional().openapi({ example: 1 }),
      from: z.string().optional().openapi({ example: "2026-01-01" }),
      to: z.string().optional().openapi({ example: "2026-12-31" }),
    }),
  },
  responses: {
    200: {
      description: "Lista de transações",
      content: { "application/json": { schema: TransactionSchema.array() } },
    },
  },
});

const getTransactionRoute = createRoute({
  method: "get",
  path: "/{id}",
  tags: ["Transactions"],
  summary: "Buscar transação por ID",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Transação encontrada",
      content: { "application/json": { schema: TransactionSchema } },
    },
    404: {
      description: "Não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const createTransactionRoute = createRoute({
  method: "post",
  path: "/",
  tags: ["Transactions"],
  summary: "Criar transação",
  description: "Cria uma transação. O campo areaId é obrigatório.",
  request: {
    body: {
      content: { "application/json": { schema: CreateTransactionSchema } },
    },
  },
  responses: {
    201: {
      description: "Transação criada",
      content: { "application/json": { schema: TransactionSchema } },
    },
    400: {
      description: "Erro de validação ou categoria não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const updateTransactionRoute = createRoute({
  method: "put",
  path: "/{id}",
  tags: ["Transactions"],
  summary: "Atualizar transação",
  description: "Use categoryId: null para desvincular uma categoria.",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
    body: {
      content: { "application/json": { schema: UpdateTransactionSchema } },
    },
  },
  responses: {
    200: {
      description: "Transação atualizada",
      content: { "application/json": { schema: TransactionSchema } },
    },
    400: {
      description: "Nenhum campo para atualizar ou categoria não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
    404: {
      description: "Não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const deleteTransactionRoute = createRoute({
  method: "delete",
  path: "/{id}",
  tags: ["Transactions"],
  summary: "Remover transação",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Transação removida",
      content: {
        "application/json": {
          schema: z.object({ message: z.string() }).openapi("DeleteResponse"),
        },
      },
    },
    404: {
      description: "Não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

app.openapi(listTransactionsRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { type, categoryId, areaId, from, to } = c.req.valid("query");

  const conditions = [];
  if (type) conditions.push(eq(transactions.type, type));
  if (categoryId !== undefined) conditions.push(eq(transactions.categoryId, categoryId));
  if (areaId !== undefined) conditions.push(eq(transactions.areaId, areaId));
  if (from) conditions.push(gte(transactions.date, from));
  if (to) conditions.push(lte(transactions.date, to));

  const where = conditions.length > 0 ? and(...conditions) : undefined;

  const rows = await db.select().from(transactions).where(where);
  return c.json(rows, 200);
});

app.openapi(getTransactionRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { id } = c.req.valid("param");

  const rows = await db
    .select()
    .from(transactions)
    .where(eq(transactions.id, id));

  if (rows.length === 0) {
    return c.json({ error: "Transação não encontrada" }, 404);
  }

  return c.json(rows[0], 200);
});

app.openapi(createTransactionRoute, async (c) => {
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  try {
    const rows = await db.insert(transactions).values(body).returning();
    return c.json(rows[0], 201);
  } catch (e: any) {
    if (e.message?.includes("FOREIGN")) {
      return c.json({ error: "Categoria ou área não encontrada" }, 400);
    }
    throw e;
  }
});

app.openapi(updateTransactionRoute, async (c) => {
  const { id } = c.req.valid("param");
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  if (Object.keys(body).length === 0) {
    return c.json({ error: "Nenhum campo para atualizar" }, 400);
  }

  try {
    const rows = await db
      .update(transactions)
      .set(body)
      .where(eq(transactions.id, id))
      .returning();

    if (rows.length === 0) {
      return c.json({ error: "Transação não encontrada" }, 404);
    }

    return c.json(rows[0], 200);
  } catch (e: any) {
    if (e.message?.includes("FOREIGN")) {
      return c.json({ error: "Categoria ou área não encontrada" }, 400);
    }
    throw e;
  }
});

app.openapi(deleteTransactionRoute, async (c) => {
  const { id } = c.req.valid("param");
  const db = createDb(c.env.finance);

  const rows = await db
    .delete(transactions)
    .where(eq(transactions.id, id))
    .returning();

  if (rows.length === 0) {
    return c.json({ error: "Transação não encontrada" }, 404);
  }

  return c.json({ message: "Transação removida com sucesso" }, 200);
});

export default app;
