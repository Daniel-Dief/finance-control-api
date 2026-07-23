import { drizzle } from "drizzle-orm/d1";
import * as schema from "../database/schema";

export type Env = {
  finance: D1Database;
};

export function createDb(d1: D1Database) {
  return drizzle(d1, { schema });
}
