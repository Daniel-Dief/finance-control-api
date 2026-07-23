import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { eq, like } from "drizzle-orm";
import { createDb, type Env } from "../db";
import { categories } from "../../database/schema";
import { ErrorSchema } from "./schemas";

const app = new OpenAPIHono<{ Bindings: Env }>();

const CategorySchema = z
  .object({
    id: z.number().int().openapi({ example: 1 }),
    name: z.string().openapi({ example: "Alimentação" }),
  })
  .openapi("Category");

const CreateCategorySchema = z
  .object({
    name: z.string().min(1).openapi({ example: "Alimentação" }),
  })
  .openapi("CreateCategory");

const UpdateCategorySchema = CreateCategorySchema.partial();

const listCategoriesRoute = createRoute({
  method: "get",
  path: "/",
  tags: ["Categories"],
  summary: "Listar categorias",
  description: " Retorna todas as categorias. Use ?name= para buscar por nome (parcial).",
  request: {
    query: z.object({
      name: z.string().optional().openapi({ example: "Alim" }),
    }),
  },
  responses: {
    200: {
      description: "Lista de categorias",
      content: { "application/json": { schema: CategorySchema.array() } },
    },
  },
});

const getCategoryRoute = createRoute({
  method: "get",
  path: "/{id}",
  tags: ["Categories"],
  summary: "Buscar categoria por ID",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Categoria encontrada",
      content: { "application/json": { schema: CategorySchema } },
    },
    404: {
      description: "Não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const createCategoryRoute = createRoute({
  method: "post",
  path: "/",
  tags: ["Categories"],
  summary: "Criar categoria",
  request: {
    body: {
      content: { "application/json": { schema: CreateCategorySchema } },
    },
  },
  responses: {
    201: {
      description: "Categoria criada",
      content: { "application/json": { schema: CategorySchema } },
    },
    400: {
      description: "Erro de validação",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const updateCategoryRoute = createRoute({
  method: "put",
  path: "/{id}",
  tags: ["Categories"],
  summary: "Atualizar categoria",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
    body: {
      content: { "application/json": { schema: UpdateCategorySchema } },
    },
  },
  responses: {
    200: {
      description: "Categoria atualizada",
      content: { "application/json": { schema: CategorySchema } },
    },
    400: {
      description: "Erro de validação",
      content: { "application/json": { schema: ErrorSchema } },
    },
    404: {
      description: "Não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const deleteCategoryRoute = createRoute({
  method: "delete",
  path: "/{id}",
  tags: ["Categories"],
  summary: "Remover categoria",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Categoria removida",
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

app.openapi(listCategoriesRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { name } = c.req.valid("query");

  if (name) {
    const rows = await db
      .select()
      .from(categories)
      .where(like(categories.name, `%${name}%`));
    return c.json(rows, 200);
  }

  const rows = await db.select().from(categories);
  return c.json(rows, 200);
});

app.openapi(getCategoryRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { id } = c.req.valid("param");

  const rows = await db
    .select()
    .from(categories)
    .where(eq(categories.id, id));

  if (rows.length === 0) {
    return c.json({ error: "Categoria não encontrada" }, 404);
  }

  return c.json(rows[0], 200);
});

app.openapi(createCategoryRoute, async (c) => {
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  const rows = await db.insert(categories).values(body).returning();
  return c.json(rows[0], 201);
});

app.openapi(updateCategoryRoute, async (c) => {
  const { id } = c.req.valid("param");
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  const rows = await db
    .update(categories)
    .set(body)
    .where(eq(categories.id, id))
    .returning();

  if (rows.length === 0) {
    return c.json({ error: "Categoria não encontrada" }, 404);
  }

  return c.json(rows[0], 200);
});

app.openapi(deleteCategoryRoute, async (c) => {
  const { id } = c.req.valid("param");
  const db = createDb(c.env.finance);

  const rows = await db
    .delete(categories)
    .where(eq(categories.id, id))
    .returning();

  if (rows.length === 0) {
    return c.json({ error: "Categoria não encontrada" }, 404);
  }

  return c.json({ message: "Categoria removida com sucesso" }, 200);
});

export default app;
