// Thin helper that opens a pg client against the dev-compose Postgres.
//
// Used by specs that need to assert DB-side effects directly (e.g. 02d
// confirms an iam.audit_logs row lands after a state-changing REST call).
// Connection string defaults mirror the top-level Makefile's LOCAL_DB_URL so
// the same credentials work on the host during Playwright runs.

import { Client } from "pg";

const DEFAULT_URL =
  "postgres://postgres:changeme@127.0.0.1:5432/umrohos_dev?sslmode=disable";

export async function withPg<T>(
  fn: (client: Client) => Promise<T>,
): Promise<T> {
  const client = new Client({
    connectionString: process.env.LOCAL_DB_URL || DEFAULT_URL,
  });
  await client.connect();
  try {
    return await fn(client);
  } finally {
    await client.end();
  }
}
