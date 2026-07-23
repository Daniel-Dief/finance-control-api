import { OpenAPIHono } from "@hono/zod-openapi";
import { swaggerUI } from "@hono/swagger-ui";
import { cors } from "hono/cors";
import categoriesRoutes from "./routes/categories";
import areasRoutes from "./routes/areas";
import budgetsRoutes from "./routes/budgets";
import transactionsRoutes from "./routes/transactions";

type Bindings = {
  finance: D1Database;
};

const app = new OpenAPIHono<{ Bindings: Bindings }>({
  defaultHook: (result, c) => {
    if (!result.success) {
      return c.json(
        { error: result.error.issues.map((i) => i.message).join(", ") },
        400
      );
    }
  },
});

app.use("/*", cors());

app.get("/", (c) => {
  return c.json({
    message: "Finance Control API",
    docs: "/ui",
    endpoints: {
      categories: "/categories",
      areas: "/areas",
      budgets: "/budgets",
      transactions: "/transactions",
    },
  });
});

app.route("/categories", categoriesRoutes);
app.route("/areas", areasRoutes);
app.route("/budgets", budgetsRoutes);
app.route("/transactions", transactionsRoutes);

app.doc("/doc", {
  openapi: "3.0.0",
  info: {
    version: "1.0.0",
    title: "Finance Control API",
    description:
      "API para controle de finanças pessoais. Gerencie categorias, áreas, orçamentos mensais e transações.",
  },
  servers: [
    { url: "http://127.0.0.1:8787", description: "Desenvolvimento local" },
  ],
});

app.get("/ui", swaggerUI({ url: "/doc" }));

app.notFound((c) => {
  return c.json({ error: "Rota não encontrada" }, 404);
});

app.onError((err, c) => {
  console.error("Unhandled error:", err);
  return c.json({ error: err.message || "Erro interno do servidor" }, 500);
});

export default app;
