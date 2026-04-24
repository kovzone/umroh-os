import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

type CatalogPackageItem = {
  id: string;
  kind: string;
  name: string;
  description: string;
  cover_photo_url: string;
  status?: string;
  starting_price: {
    list_amount: number;
    list_currency: string;
    settlement_currency: string;
  };
  next_departure?: {
    id: string;
    departure_date: string;
    return_date: string;
    remaining_seats: number;
  };
  departures_count?: number;
};

type CatalogListResponse = {
  packages: CatalogPackageItem[];
};

// GATEWAY_URL must come first: it's set at runtime in docker-compose (server-side).
// VITE_* vars are baked into the client bundle at build time and NOT available
// to the SSR Node.js process at runtime in production containers.
const baseUrl =
  process.env.GATEWAY_URL ??
  process.env.VITE_CATALOG_API_BASE_URL ??
  process.env.VITE_GATEWAY_URL ??
  'http://localhost:4000';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  try {
    const response = await fetch(`${baseUrl}/v1/packages`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    });

    if (!response.ok) {
      return { packages: [], error: `Gagal memuat daftar paket (${response.status})` };
    }

    const body = (await response.json()) as CatalogListResponse;

    return {
      packages: body.packages ?? [],
      error: null
    };
  } catch {
    return { packages: [], error: 'Tidak dapat terhubung ke layanan katalog.' };
  }
};
