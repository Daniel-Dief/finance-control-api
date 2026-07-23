import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { eq, like } from "drizzle-orm";
import { createDb, type Env } from "../db";
import { areas } from "../../database/schema";
import { ErrorSchema } from "./schemas";

const app = new OpenAPIHono<{ Bindings: Env }>();

const AreaSchema = z
  .object({
    id: z.number().int().openapi({ example: 1 }),
    name: z.string().openapi({ example: "Moradia" }),
  })
  .openapi("Area");

const CreateAreaSchema = z
  .object({
    name: z.string().min(1).openapi({ example: "Moradia" }),
  })
  .openapi("CreateArea");

const UpdateAreaSchema = CreateAreaSchema.partial();

const listAreasRoute = createRoute({
  method: "get",
  path: "/",
  tags: ["Areas"],
  summary: "Listar áreas",
  description: "Retorna todas as áreas. Use ?name= para buscar por nome (parcial).",
  request: {
    query: z.object({
      name: z.string().optional().openapi({ example: "Mora" }),
    }),
  },
  responses: {
    200: {
      description: "Lista de áreas",
      content: { "application/json": { schema: AreaSchema.array() } },
    },
  },
});

const getAreaRoute = createRoute({
  method: "get",
  path: "/{id}",
  tags: ["Areas"],
  summary: "Buscar área por ID",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Área encontrada",
      content: { "application/json": { schema: AreaSchema } },
    },
    404: {
      description: "Não encontrada",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const createAreaRoute = createRoute({
  method: "post",
  path: "/",
  tags: ["Areas"],
  summary: "Criar área",
  request: {
    body: {
      content: { "application/json": { schema: CreateAreaSchema } },
    },
  },
  responses: {
    201: {
      description: "Área criada",
      content: { "application/json": { schema: AreaSchema } },
    },
    400: {
      description: "Erro de validação",
      content: { "application/json": { schema: ErrorSchema } },
    },
  },
});

const updateAreaRoute = createRoute({
  method: "put",
  path: "/{id}",
  tags: ["Areas"],
  summary: "Atualizar área",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
    body: {
      content: { "application/json": { schema: UpdateAreaSchema } },
    },
  },
  responses: {
    200: {
      description: "Área atualizada",
      content: { "application/json": { schema: AreaSchema } },
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

const deleteAreaRoute = createRoute({
  method: "delete",
  path: "/{id}",
  tags: ["Areas"],
  summary: "Remover área",
  request: {
    params: z.object({
      id: z.coerce.number().int().openapi({ param: { name: "id", in: "path" }, example: 1 }),
    }),
  },
  responses: {
    200: {
      description: "Área removida",
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

app.openapi(listAreasRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { name } = c.req.valid("query");

  if (name) {
    const rows = await db
      .select()
      .from(areas)
      .where(like(areas.name, `%${name}%`));
    return c.json(rows, 200);
  }

  const rows = await db.select().from(areas);
  return c.json(rows, 200);
});

app.openapi(getAreaRoute, async (c) => {
  const db = createDb(c.env.finance);
  const { id } = c.req.valid("param");

  const rows = await db.select().from(areas).where(eq(areas.id, id));

  if (rows.length === 0) {
    return c.json({ error: "Área não encontrada" }, 404);
  }

  return c.json(rows[0], 200);
});

app.openapi(createAreaRoute, async (c) => {
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  const rows = await db.insert(areas).values(body).returning();
  return c.json(rows[0], 201);
});

app.openapi(updateAreaRoute, async (c) => {
  const { id } = c.req.valid("param");
  const body = c.req.valid("json");
  const db = createDb(c.env.finance);

  const rows = await db
    .update(areas)
    .set(body)
    .where(eq(areas.id, id))
    .returning();

  if (rows.length === 0) {
    return c.json({ error: "Área não encontrada" }, 404);
  }

  return c.json(rows[0], 200);
});

app.openapi(deleteAreaRoute, async (c) => {
  const { id } = c.req.valid("param");
  const db = createDb(c.env.finance);

  const rows = await db
    .delete(areas)
    .where(eq(areas.id, id))
    .returning();

  if (rows.length === 0) {
    return c.json({ error: "Área não encontrada" }, 404);
  }

  return c.json({ message: "Área removida com sucesso" }, 200);
});

export default app;
