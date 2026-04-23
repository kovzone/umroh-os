// -----------------------------------------------------------------------
// S4-CRM — API Wrapper
// All calls go through gateway-svc:4000
// Toggle VITE_MOCK_CRM=true to use mock data
// -----------------------------------------------------------------------
import type {
  CreateLeadRequest,
  CSUser,
  Lead,
  LeadFilters,
  LeadList,
  PackageOption,
  UpdateLeadRequest
} from './types';
import {
  createLeadMock,
  listCsUsersMock,
  listLeadsMock,
  listPackagesMock,
  updateLeadMock
} from './mock';

const MOCK = import.meta.env.VITE_MOCK_CRM === 'true';
const baseUrl = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

// --------------- helpers ----------------------------------------------

async function apiFetch<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(`${baseUrl}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error((err as { message?: string }).message ?? `HTTP ${res.status}`);
  }
  return res.json() as Promise<T>;
}

// --------------- public API -------------------------------------------

/**
 * POST /v1/leads
 * Submit a new lead (public lead capture form).
 */
export async function createLead(req: CreateLeadRequest): Promise<{ id: string }> {
  if (MOCK) return createLeadMock(req);
  return apiFetch<{ id: string }>('/v1/leads', {
    method: 'POST',
    body: JSON.stringify(req)
  });
}

/**
 * GET /v1/leads
 * List leads with optional filters (CS console).
 */
export async function listLeads(filters: LeadFilters = {}): Promise<LeadList> {
  if (MOCK) return listLeadsMock(filters);

  const params = new URLSearchParams();
  if (filters.status && filters.status !== 'all') params.set('status', filters.status);
  if (filters.assigned_cs_id) params.set('assigned_cs_id', filters.assigned_cs_id);
  if (filters.search) params.set('search', filters.search);

  return apiFetch<LeadList>(`/v1/leads?${params.toString()}`);
}

/**
 * PUT /v1/leads/{id}
 * Update lead status/notes/assignment.
 */
export async function updateLead(id: string, req: UpdateLeadRequest): Promise<Lead> {
  if (MOCK) return updateLeadMock(id, req);
  return apiFetch<Lead>(`/v1/leads/${id}`, {
    method: 'PUT',
    body: JSON.stringify(req)
  });
}

/**
 * GET /v1/cs-users
 * List active CS users for assignment dropdown.
 */
export async function listCsUsers(): Promise<CSUser[]> {
  if (MOCK) return listCsUsersMock();
  return apiFetch<CSUser[]>('/v1/cs-users');
}

/**
 * GET /v1/packages (public catalog) — used in lead capture form dropdown.
 */
export async function listPackages(): Promise<PackageOption[]> {
  if (MOCK) return listPackagesMock();
  return apiFetch<{ packages: PackageOption[] }>('/v1/packages').then(
    (r) => r.packages ?? []
  );
}
