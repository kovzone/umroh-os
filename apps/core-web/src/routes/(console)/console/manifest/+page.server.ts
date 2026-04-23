import { redirect } from '@sveltejs/kit';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

type LoadEvent = { cookies: Cookies; fetch: RequestEvent['fetch'] };

const MOCK = process.env.MOCK_OPS === 'true' || process.env.VITE_MOCK_OPS === 'true';
const baseUrl =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export type DepartureStatus = 'upcoming' | 'ongoing' | 'completed' | 'cancelled';

export interface DepartureSummary {
  id: string;
  package_name: string;
  departure_date: string;
  return_date: string;
  total_seats: number;
  booked_seats: number;
  status: DepartureStatus;
}

export interface PageData {
  departures: DepartureSummary[];
  error: string | null;
}

function mockData(): PageData {
  return {
    departures: [
      {
        id: 'dep1',
        package_name: 'Umroh Ramadhan Premium',
        departure_date: '2026-03-01',
        return_date: '2026-03-14',
        total_seats: 45,
        booked_seats: 38,
        status: 'upcoming'
      },
      {
        id: 'dep2',
        package_name: 'Umroh Reguler April',
        departure_date: '2026-04-10',
        return_date: '2026-04-24',
        total_seats: 50,
        booked_seats: 42,
        status: 'upcoming'
      },
      {
        id: 'dep3',
        package_name: 'Umroh Plus Turki',
        departure_date: '2026-05-05',
        return_date: '2026-05-19',
        total_seats: 40,
        booked_seats: 15,
        status: 'upcoming'
      },
      {
        id: 'dep4',
        package_name: 'Umroh Reguler Februari',
        departure_date: '2026-02-10',
        return_date: '2026-02-24',
        total_seats: 50,
        booked_seats: 50,
        status: 'completed'
      }
    ],
    error: null
  };
}

export const load = async ({ cookies, fetch }: LoadEvent): Promise<PageData> => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  if (MOCK) {
    return mockData();
  }

  try {
    const res = await fetch(`${baseUrl}/v1/departures`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    });

    if (res.ok) {
      const body = (await res.json()) as DepartureSummary[] | { departures: DepartureSummary[] };
      const departures = Array.isArray(body) ? body : body.departures ?? [];
      return { departures, error: null };
    }

    return { departures: [], error: `Gagal memuat keberangkatan (${res.status})` };
  } catch {
    return { departures: [], error: 'Tidak dapat terhubung ke gateway.' };
  }
};
