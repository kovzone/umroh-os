// Typed REST client for gateway-svc.
//
// Types are generated from services/gateway-svc/api/rest_oapi/openapi.yaml via
//   npm run gen:api
// which writes `./schema.d.ts`. The runtime is openapi-fetch — a tiny typed
// wrapper around `fetch` that constrains the path argument to the
// generated `paths` union.
//
// Base URL comes from the `VITE_GATEWAY_URL` env var (browser-side). In dev
// the browser hits the host port directly (http://localhost:4000); in prod
// the same env var points at the deployed gateway.

import createClient from 'openapi-fetch';
import type { paths } from './schema';

const baseUrl = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export const gateway = createClient<paths>({ baseUrl });
