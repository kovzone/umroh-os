/**
 * Typed wrapper for catalog write endpoints (POST/PUT).
 * These endpoints require Bearer token authorization.
 *
 * NOTE: openapi-fetch schema.d.ts does not yet include write endpoints
 * (pending S1-E-07 schema regeneration). This manual wrapper bridges
 * the gap until the schema is updated and npm run gen:api is re-run.
 *
 * All public functions accept an optional `baseUrl` so that server-side
 * callers (SvelteKit +page.server.ts) can pass process.env values
 * instead of relying on import.meta.env (which is Vite/browser-only).
 */

// ---- Request / Response shapes ----

export type PackageStatus = 'draft' | 'active';
export type PackageKind = 'umroh' | 'hajj' | 'ziarah';

export type CreatePackageRequest = {
  name: string;
  description: string;
  kind: PackageKind;
  status: PackageStatus;
  starting_price_idr: number;
};

export type UpdatePackageRequest = Partial<CreatePackageRequest>;

export type PackageWriteResponse = {
  package: {
    id: string;
    name: string;
    description: string;
    kind: string;
    status: PackageStatus;
    cover_photo_url: string;
    starting_price: {
      list_amount: number;
      list_currency: string;
      settlement_currency: string;
    };
  };
};

export type DepartureStatus = 'open' | 'closed' | 'draft';

export type CreateDepartureRequest = {
  departure_date: string;
  return_date: string;
  total_seats: number;
  status: DepartureStatus;
};

export type UpdateDepartureRequest = Partial<CreateDepartureRequest>;

export type DepartureWriteResponse = {
  departure: {
    id: string;
    package_id: string;
    departure_date: string;
    return_date: string;
    total_seats: number;
    remaining_seats: number;
    status: DepartureStatus;
  };
};

// ---- Error handling ----

export class CatalogWriteError extends Error {
  constructor(
    public readonly status: number,
    message: string
  ) {
    super(message);
    this.name = 'CatalogWriteError';
  }
}

// ---- Internal helpers ----

function resolveBase(override?: string): string {
  if (override) return override;
  // Server-side (Node): try process.env
  if (typeof process !== 'undefined' && process.env) {
    return (
      process.env['VITE_CATALOG_API_BASE_URL'] ??
      process.env['VITE_GATEWAY_URL'] ??
      process.env['GATEWAY_URL'] ??
      'http://localhost:4000'
    );
  }
  return 'http://localhost:4000';
}

async function writeJson<T>(
  path: string,
  method: 'POST' | 'PUT',
  body: unknown,
  bearerToken: string,
  baseUrl?: string
): Promise<T> {
  const url = `${resolveBase(baseUrl)}${path}`;
  const response = await fetch(url, {
    method,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${bearerToken}`
    },
    body: JSON.stringify(body)
  });

  if (!response.ok) {
    const err = await response.json().catch(() => ({ error: { message: 'Unknown error' } }));
    const message = (err as { error?: { message?: string } }).error?.message ?? `HTTP ${response.status}`;
    throw new CatalogWriteError(response.status, message);
  }

  return (await response.json()) as T;
}

// ---- Public API ----

export async function createPackage(
  data: CreatePackageRequest,
  bearerToken: string,
  baseUrl?: string
): Promise<PackageWriteResponse> {
  return writeJson<PackageWriteResponse>('/v1/packages', 'POST', data, bearerToken, baseUrl);
}

export async function updatePackage(
  packageId: string,
  data: UpdatePackageRequest,
  bearerToken: string,
  baseUrl?: string
): Promise<PackageWriteResponse> {
  return writeJson<PackageWriteResponse>(`/v1/packages/${packageId}`, 'PUT', data, bearerToken, baseUrl);
}

export async function createDeparture(
  packageId: string,
  data: CreateDepartureRequest,
  bearerToken: string,
  baseUrl?: string
): Promise<DepartureWriteResponse> {
  return writeJson<DepartureWriteResponse>(
    `/v1/packages/${packageId}/departures`,
    'POST',
    data,
    bearerToken,
    baseUrl
  );
}

export async function updateDeparture(
  departureId: string,
  data: UpdateDepartureRequest,
  bearerToken: string,
  baseUrl?: string
): Promise<DepartureWriteResponse> {
  return writeJson<DepartureWriteResponse>(
    `/v1/departures/${departureId}`,
    'PUT',
    data,
    bearerToken,
    baseUrl
  );
}
