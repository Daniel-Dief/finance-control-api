import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { eq, and } from "drizzle-orm";
import { createDb, type Env } from "../db";
import { monthlyBudget } from "../../database/schema";
import { ErrorSchema } from "./schemas";

const app = new OpenAPIHono<{ Bindings: Env }>();

const BudgetSchema = z
  .object({
    id: z.number().int().openapi({ example: 1 }),
    year: z.number().int().openapi({ example: 2026 }),
    month: z.number().int().openapi({ example: 7 }),
    areaId: z.number().int().openapi({ example: 1 }),
    amount: z.number().openapi({ example: 1500.0 }),
  })
  .openapi("Budget");

const CreateBudgetSchema = z
  .object({
    year: z.number().int().openapi({ example: 2026 }),
    month: z.number().int().min(1).max(12).openapi({ example: 7 }),
    areaId: z.number().int().openapi({ example: 1 }),
    amount: z.number().min(0).openapi({ example: 1500.0 }),
  })
  .openapi("CreateBudget");

const UpdateBudgetSchema = CreateBudgetSchema.partial();

const listBudgetsRoute = createRoute({
  method: "get",
  path: "/",
  tags: ["Budgets"],
  summary: "Listar orçamentos",
  description: "Retorna todos os orçamentos. Filtre por ?year=, ?month= e/ou ?areaId=.",
  request: {
    query: z.object({
      year: z.coerce.number().int().optional().openapi({ example: 2026 }),
      month: z.coerce.number().int().optional().openapi({ example: 7 }),
      areaId: z.coerce.number().int().optional().openapi({ example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Lista de orçamentos",
      content: { "application/json": { schema: BudgetSchema.array() } },
    },
  },
});

const getBudgetRoute = createRoute({
  method: "get",
  path: "/{id}",
  tags: ["Budgets"],
  summary: "Buscar orçamento por ID",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Orçamento encontrado",
      content: { "application/json": { schema: BudgetSchema } },
    },
    404: {
      description: "Não encontrado",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const createBudgetRoute = createRoute({
  method: "post",
  path: "/",
  tags: ["Budgets"],
  summary: "Criar orçamento",
  description: "Cria um orçamento mensal para uma área. Não pode duplicar (year+month+areaId).",
  request: {
    body: {
      content: { "application/json": { schema: CreateBudgetSchema } },
    },
  },
  responses: {
    201: {
      description: "Orçamento criado",
      content: { "application/json": { schema: BudgetSchema } },
    },
    400: {
      description: "Erro de validação ou área não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
    409: {
      description: "Já existe orçamento para esta área no período",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const updateBudgetRoute = createRoute({
  method: "put",
  path: "/{id}",
  tags: ["Budgets"],
  summary: "Atualizar orçamento",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
    body: {
      content: { "application/json": { schema: UpdateBudgetSchema } },
    },
  },
  responses: {
    200: {
      description: "Orçamento atualizado",
      content: { "application/json": { schema: BudgetSchema } },
    },
    400: {
      description: "Nenhum campo para atualizar",
      content: { "application/json": { schema: ErrorSchema } },
    },
    404: {
      description: "Não encontrado",
      content: { "application/json": { schema: ErrorSchema } },
    },
    409: {
      description: "Conflito de período/área",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const deleteBudgetRoute = createRoute({
  method: "delete",
  path: "/{id}",
  tags: ["Budgets"],
  summary: "Remover orçamento",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Orçamento removido",
      content: {
        "application/json": {
          schema: z.object({ message: z.string() }).openapi("DeleteResponse"),
        },
      },
    },
    404: {
      description: "Não encontrado",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

app.openapi(listBudgetsRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { year, month, areaId } = c.req.valid("query");

  const conditions = [];
  if (year !== undefined) conditions.push(eq(monthlyBudget.year, year));
  if (month !== undefined) conditions.push(eq(monthlyBudget.month, month));
  if (areaId !== undefined) conditions.push(eq(monthlyBudget.areaId, areaId));

  const where = conditions.length > 0 ? and(...conditions) : undefined;

  const rows = await db.select().from(monthlyBudget).where(where);
  return c.json(rows, 200);
});

app.openapi(getBudgetRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { id } = c.req.valid("param");

  const rows = await db
    .select()
    .from(monthlyBudget)
    .where(eq(monthlyBudget.id, id));

  if (rows.length === 0) {
    return c.json({ error: "Orçamento não encontrado" }, 404);
  }

  return c.json(rows[0], 200);
});

app.openapi(createBudgetRoute, async (c) => {
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  try {
    const rows = await db.insert(monthlyBudget).values(body).returning();
    return c.json(rows[0], 201);
  } catch (e: any) {
    if (e.message?.includes("UNIQUE")) {
      return c.json(
        { error: "Já existe um orçamento para esta área no período informado" },
        409
      );
    }
    if (e.message?.includes("FOREIGN")) {
      return c.json({ error: "Área não encontrada" }, 400);
    }
    throw e;
  }
});

app.openapi(updateBudgetRoute, async (c) => {
  const { id } = c.req.valid("param");
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  if (Object.keys(body).length === 0) {
    return c.json({ error: "Nenhum campo para atualizar" }, 400);
  }

  try {
    const rows = await db
      .update(monthlyBudget)
      .set(body)
      .where(eq(monthlyBudget.id, id))
      .returning();

    if (rows.length === 0) {
      return c.json({ error: "Orçamento não encontrado" }, 404);
    }

    return c.json(rows[0], 200);
  } catch (e: any) {
    if (e.message?.includes("UNIQUE")) {
      return c.json(
        { error: "Já existe um orçamento para esta área no período informado" },
        409
      );
    }
    throw e;
  }
});

app.openapi(deleteBudgetRoute, async (c) => {
  const { id } = c.req.valid("param");
  const db = createDb(c.env.finance);

  const rows = await db
    .delete(monthlyBudget)
    .where(eq(monthlyBudget.id, id))
    .returning();

  if (rows.length === 0) {
    return c.json({ error: "Orçamento não encontrado" }, 404);
  }

  return c.json({ message: "Orçamento removido com sucesso" }, 200);
});

export default app;
