import type { DraftBookingError, DraftBookingPayload, DraftBookingResult } from './types';

const baseUrl = import.meta.env.VITE_BOOKING_API_BASE_URL ?? import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

function mapPayload(payload: DraftBookingPayload) {
  return {
    channel: payload.channel,
    package_id: payload.packageId,
    departure_id: payload.departureId,
    room_type: payload.roomType,
    lead: {
      full_name: payload.lead.fullName,
      email: payload.lead.email,
      whatsapp: payload.lead.whatsapp,
      domicile: payload.lead.domicile
    },
    jamaah: Array.from({ length: payload.jamaahCount }, (_, idx) => ({
      full_name: idx === 0 ? payload.lead.fullName : `${payload.lead.fullName} #${idx + 1}`,
      email: idx === 0 ? payload.lead.email : undefined,
      whatsapp: idx === 0 ? payload.lead.whatsapp : undefined,
      domicile: payload.lead.domicile,
      is_lead: idx === 0
    })),
    add_on_ids: [],
    agent_id: null,
    notes: null
  };
}

export async function createDraftBookingApi(
  payload: DraftBookingPayload,
  idempotencyKey: string
): Promise<DraftBookingResult> {
  const response = await fetch(`${baseUrl}/v1/bookings`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Idempotency-Key': idempotencyKey
    },
    body: JSON.stringify(mapPayload(payload))
  });

  const body = (await response.json()) as {
    booking?: { id: string; status: 'draft'; created_at: string };
    error?: { code?: string; message?: string };
  };

  if (!response.ok || !body.booking) {
    const err: DraftBookingError = {
      code: body.error?.code ?? 'request_failed',
      message: body.error?.message ?? 'Draft booking request failed'
    };
    throw err;
  }

  return {
    bookingId: body.booking.id,
    status: body.booking.status,
    createdAt: body.booking.created_at,
    replayed: response.headers.get('Idempotency-Replayed') === 'true'
  };
}
