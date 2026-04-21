import { createDraftBookingApi } from './api';
import { createDraftBookingMock } from './mock';
import type { DraftBookingPayload, DraftBookingResult } from './types';

const useMockBooking = (import.meta.env.VITE_USE_BOOKING_MOCK ?? 'true') === 'true';

export function makeIdempotencyKey(packageId: string): string {
  const randomPart = Math.random().toString(36).slice(2, 10);
  return `draft-${packageId}-${Date.now()}-${randomPart}`;
}

export async function createDraftBooking(payload: DraftBookingPayload): Promise<DraftBookingResult> {
  const key = makeIdempotencyKey(payload.packageId);

  if (useMockBooking) {
    return createDraftBookingMock(payload, key);
  }

  try {
    return await createDraftBookingApi(payload, key);
  } catch {
    return createDraftBookingMock(payload, key);
  }
}
