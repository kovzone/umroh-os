import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

type CatalogDeparture = {
  id: string;
  package_id: string;
  departure_date: string;
  return_date: string;
  total_seats: number;
  remaining_seats: number;
  status: 'open' | 'closed' | 'draft';
};

type CatalogPackageDetailResponse = {
  package: {
    id: string;
    kind: string;
    name: string;
    description: string;
    status?: string;
    cover_photo_url: string;
    starting_price?: {
      list_amount: number;
      list_currency: string;
    };
    highlights: string[];
    departures: CatalogDeparture[];
  };
};

const baseUrl =
  process.env.VITE_CATALOG_API_BASE_URL ??
  process.env.VITE_GATEWAY_URL ??
  process.env.GATEWAY_URL ??
  'http://localhost:4000';

export const load: PageServerLoad = async ({ cookies, params, fetch }) => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  const response = await fetch(`${baseUrl}/v1/packages/${params.id}`, {
    headers: { Authorization: `Bearer ${accessToken}` }
  });

  if (response.status === 404) {
    throw error(404, 'Paket tidak ditemukan.');
  }

  if (!response.ok) {
    throw error(502, 'Gagal memuat data paket.');
  }

  const body = (await response.json()) as CatalogPackageDetailResponse;
  return {
    package: body.package,
    departures: body.package.departures ?? []
  };
};
