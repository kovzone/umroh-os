import { createDraftBookingApi } from './api';
import { createDraftBookingMock } from './mock';
import type { DraftBookingPayload, DraftBookingResult } from './types';

// Default is false — mock must be explicitly enabled via VITE_USE_BOOKING_MOCK=true.
// A silent fallback to mock when the real API fails would hide booking-svc outages
// from users, which caused ISSUE-022 (user received a fake booking code).
const useMockBooking = (import.meta.env.VITE_USE_BOOKING_MOCK ?? 'false') === 'true';

export function makeIdempotencyKey(packageId: string): string {
  const randomPart = Math.random().toString(36).slice(2, 10);
  return `draft-${packageId}-${Date.now()}-${randomPart}`;
}

export async function createDraftBooking(payload: DraftBookingPayload): Promise<DraftBookingResult> {
  const key = makeIdempotencyKey(payload.packageId);

  if (useMockBooking) {
    return createDraftBookingMock(payload, key);
  }

  // No silent fallback — let the error propagate so the UI can show it.
  return createDraftBookingApi(payload, key);
}
