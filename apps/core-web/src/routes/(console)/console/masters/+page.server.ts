import { redirect } from '@sveltejs/kit';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

type LoadEvent = { cookies: Cookies; fetch: RequestEvent['fetch'] };

const MOCK = process.env.MOCK_CATALOG === 'true' || process.env.VITE_MOCK_CATALOG === 'true';
const baseUrl =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export interface Hotel {
  id: string;
  name: string;
  city: string;
  stars: number;
  distance_m: number;
}

export interface Airline {
  code: string;
  name: string;
  type: 'airline' | 'rail' | 'bus';
}

export interface Muthawwif {
  id: string;
  name: string;
  photo_url: string | null;
}

export interface Addon {
  id: string;
  name: string;
  price_idr: number;
}

export interface PageData {
  hotels: Hotel[];
  airlines: Airline[];
  muthawwifs: Muthawwif[];
  addons: Addon[];
  error: string | null;
}

function mockData(): PageData {
  return {
    hotels: [
      { id: 'h1', name: 'Hotel Dar Al Iman Royal', city: 'Madinah', stars: 5, distance_m: 150 },
      { id: 'h2', name: 'Pullman ZamZam Makkah', city: 'Makkah', stars: 5, distance_m: 50 },
      { id: 'h3', name: 'Al Shohada Hotel', city: 'Makkah', stars: 4, distance_m: 600 },
      { id: 'h4', name: 'Crowne Plaza Madinah', city: 'Madinah', stars: 4, distance_m: 350 },
      { id: 'h5', name: 'Dallah Taibah Hotel', city: 'Madinah', stars: 3, distance_m: 900 }
    ],
    airlines: [
      { code: 'GA', name: 'Garuda Indonesia', type: 'airline' },
      { code: 'SV', name: 'Saudi Arabian Airlines', type: 'airline' },
      { code: 'EK', name: 'Emirates', type: 'airline' },
      { code: 'QZ', name: 'AirAsia', type: 'airline' }
    ],
    muthawwifs: [
      { id: 'm1', name: 'Ust. Ahmad Fauzi', photo_url: null },
      { id: 'm2', name: 'Ust. Hasan Basri', photo_url: null },
      { id: 'm3', name: 'Ust. Ridwan Kamil', photo_url: null }
    ],
    addons: [
      { id: 'a1', name: 'Koper Koper Premium 28"', price_idr: 350000 },
      { id: 'a2', name: 'Travel Kit Lengkap', price_idr: 250000 },
      { id: 'a3', name: 'Asuransi Perjalanan', price_idr: 175000 },
      { id: 'a4', name: 'Upgrade Kamar Single', price_idr: 3500000 }
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
    const headers = { Authorization: `Bearer ${accessToken}` };
    const [hotelsRes, airlinesRes, muthawwifRes, addonsRes] = await Promise.all([
      fetch(`${baseUrl}/v1/catalog/hotels`, { headers }),
      fetch(`${baseUrl}/v1/catalog/airlines`, { headers }),
      fetch(`${baseUrl}/v1/catalog/muthawwif`, { headers }),
      fetch(`${baseUrl}/v1/catalog/addons`, { headers })
    ]);

    return {
      hotels: hotelsRes.ok ? ((await hotelsRes.json()) as Hotel[]) : [],
      airlines: airlinesRes.ok ? ((await airlinesRes.json()) as Airline[]) : [],
      muthawwifs: muthawwifRes.ok ? ((await muthawwifRes.json()) as Muthawwif[]) : [],
      addons: addonsRes.ok ? ((await addonsRes.json()) as Addon[]) : [],
      error: !hotelsRes.ok ? `Gagal memuat data katalog (${hotelsRes.status})` : null
    };
  } catch {
    return { ...mockData(), error: 'Tidak dapat terhubung ke gateway.' };
  }
};
