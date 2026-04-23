import type {
  FulfillmentTask,
  FulfillmentTaskList,
  OpsFilters,
  StatusUpdate
} from './types';
import {
  listFulfillmentTasksMock,
  updateFulfillmentStatusMock
} from './mock';

const MOCK = import.meta.env.VITE_MOCK_OPS === 'true';
const baseUrl = import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

/**
 * GET /v1/fulfillment-tasks
 * Query params: status?, departure_id?, search?
 */
export async function listFulfillmentTasks(filters: OpsFilters): Promise<FulfillmentTaskList> {
  if (MOCK) return listFulfillmentTasksMock(filters);

  const params = new URLSearchParams();
  if (filters.status && filters.status !== 'all') params.set('status', filters.status);
  if (filters.departure_id) params.set('departure_id', filters.departure_id);
  if (filters.search) params.set('search', filters.search);

  const res = await fetch(`${baseUrl}/v1/fulfillment-tasks?${params.toString()}`, {
    headers: { 'Content-Type': 'application/json' }
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error((err as { message?: string }).message ?? `HTTP ${res.status}`);
  }

  return res.json() as Promise<FulfillmentTaskList>;
}

/**
 * PUT /v1/fulfillment-tasks/{id}/status
 * Request: { status, tracking_number? }
 */
export async function updateFulfillmentStatus(
  taskId: string,
  update: StatusUpdate
): Promise<FulfillmentTask> {
  if (MOCK) return updateFulfillmentStatusMock(taskId, update);

  const res = await fetch(`${baseUrl}/v1/fulfillment-tasks/${taskId}/status`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(update)
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error((err as { message?: string }).message ?? `HTTP ${res.status}`);
  }

  return res.json() as Promise<FulfillmentTask>;
}
