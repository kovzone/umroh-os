import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { updatePackage, CatalogWriteError } from '$lib/features/s1-catalog/write-api';

type CatalogPackageDetailResponse = {
  package: {
    id: string;
    kind: string;
    name: string;
    description: string;
    cover_photo_url: string;
    status?: string;
    starting_price?: {
      list_amount: number;
      list_currency: string;
    };
    highlights: string[];
    departures: unknown[];
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
    throw error(502, 'Gagal memuat detail paket.');
  }

  const body = (await response.json()) as CatalogPackageDetailResponse;
  return { package: body.package };
};

export const actions: Actions = {
  default: async ({ request, cookies, params }) => {
    const accessToken = cookies.get('umrohos_access_token');
    if (!accessToken) {
      throw redirect(303, '/console/login');
    }

    const form = await request.formData();
    const name = String(form.get('name') ?? '').trim();
    const description = String(form.get('description') ?? '').trim();
    const kind = String(form.get('kind') ?? 'umroh') as 'umroh' | 'hajj' | 'ziarah';
    const status = String(form.get('status') ?? 'draft') as 'draft' | 'active';
    const priceRaw = String(form.get('starting_price_idr') ?? '').replace(/\D/g, '');
    const starting_price_idr = Number(priceRaw);

    const errors: Record<string, string> = {};
    if (!name) errors.name = 'Nama paket wajib diisi.';
    if (!starting_price_idr || starting_price_idr <= 0)
      errors.starting_price_idr = 'Harga harus lebih dari 0.';

    if (Object.keys(errors).length > 0) {
      return fail(400, {
        errors,
        values: { name, description, kind, status, starting_price_idr: priceRaw }
      });
    }

    try {
      await updatePackage(params.id, { name, description, kind, status, starting_price_idr }, accessToken);
    } catch (e) {
      if (e instanceof CatalogWriteError) {
        return fail(e.status, {
          errors: {},
          serverError: e.message,
          values: { name, description, kind, status, starting_price_idr: priceRaw }
        });
      }
      return fail(500, {
        errors: {},
        serverError: 'Gagal menyimpan perubahan. Silakan coba lagi.',
        values: { name, description, kind, status, starting_price_idr: priceRaw }
      });
    }

    throw redirect(303, '/console/packages');
  }
};
