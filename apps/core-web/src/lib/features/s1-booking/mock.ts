import type { DraftBookingPayload, DraftBookingResult } from './types';

function randomId(prefix: string): string {
  return `${prefix}_${Math.random().toString(36).slice(2, 10)}`;
}

export async function createDraftBookingMock(
  _payload: DraftBookingPayload,
  _idempotencyKey: string
): Promise<DraftBookingResult> {
  await new Promise((resolve) => setTimeout(resolve, 350));

  return {
    bookingId: randomId('bkg_demo'),
    status: 'draft',
    createdAt: new Date().toISOString(),
    replayed: false
  };
}
