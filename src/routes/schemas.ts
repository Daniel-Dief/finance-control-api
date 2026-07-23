import { z } from "@hono/zod-openapi";

export const ErrorSchema = z
  .object({
    error: z.string().openapi({ example: "Mensagem de erro" }),
  })
  .openapi("Error");
