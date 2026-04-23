import { redirect } from '@sveltejs/kit';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

type LoadEvent = { cookies: Cookies; fetch: RequestEvent['fetch']; params: Record<string, string> };

const MOCK = process.env.MOCK_CATALOG === 'true' || process.env.VITE_MOCK_CATALOG === 'true';
const baseUrl =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export interface RoomPricing {
  room_type: 'double' | 'triple' | 'quad';
  price_idr: number;
}

export interface DeparturePricing {
  departure_id: string;
  departure_date: string;
  airline: string;
  pricing: RoomPricing[];
}

export interface PageData {
  package_id: string;
  package_name: string;
  departures: DeparturePricing[];
  error: string | null;
}

function mockData(packageId: string): PageData {
  return {
    package_id: packageId,
    package_name: 'Umroh Reguler Premium',
    departures: [
      {
        departure_id: 'dep1',
        departure_date: '2026-03-01',
        airline: 'Garuda Indonesia',
        pricing: [
          { room_type: 'double', price_idr: 28500000 },
          { room_type: 'triple', price_idr: 25000000 },
          { room_type: 'quad', price_idr: 22000000 }
        ]
      },
      {
        departure_id: 'dep2',
        departure_date: '2026-04-10',
        airline: 'Saudi Arabian Airlines',
        pricing: [
          { room_type: 'double', price_idr: 26500000 },
          { room_type: 'triple', price_idr: 23500000 },
          { room_type: 'quad', price_idr: 21000000 }
        ]
      },
      {
        departure_id: 'dep3',
        departure_date: '2026-05-05',
        airline: 'Emirates',
        pricing: [
          { room_type: 'double', price_idr: 29000000 },
          { room_type: 'triple', price_idr: 25500000 },
          { room_type: 'quad', price_idr: 22500000 }
        ]
      }
    ],
    error: null
  };
}

export const load = async ({ cookies, fetch, params }: LoadEvent): Promise<PageData> => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  const { id } = params;

  if (MOCK) {
    return mockData(id);
  }

  try {
    const headers = { Authorization: `Bearer ${accessToken}` };
    const [pkgRes, depRes] = await Promise.all([
      fetch(`${baseUrl}/v1/catalog/packages/${id}`, { headers }),
      fetch(`${baseUrl}/v1/catalog/packages/${id}/departures/pricing`, { headers })
    ]);

    const pkg = pkgRes.ok ? ((await pkgRes.json()) as { id: string; name: string }) : null;
    const departures = depRes.ok ? ((await depRes.json()) as DeparturePricing[]) : [];

    return {
      package_id: id,
      package_name: pkg?.name ?? '-',
      departures,
      error: pkgRes.ok ? null : `Gagal memuat data paket (${pkgRes.status})`
    };
  } catch {
    return { ...mockData(id), error: 'Tidak dapat terhubung ke gateway.' };
  }
};
