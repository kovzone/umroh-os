import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { updateDeparture, CatalogWriteError } from '$lib/features/s1-catalog/write-api';

type DepartureDetailResponse = {
  departure: {
    id: string;
    package_id: string;
    departure_date: string;
    return_date: string;
    total_seats: number;
    remaining_seats: number;
    status: 'open' | 'closed' | 'draft';
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

  // Fetch departure detail
  const depRes = await fetch(`${baseUrl}/v1/package-departures/${params.dep_id}`, {
    headers: { Authorization: `Bearer ${accessToken}` }
  });

  if (depRes.status === 404) {
    throw error(404, 'Keberangkatan tidak ditemukan.');
  }

  if (!depRes.ok) {
    throw error(502, 'Gagal memuat data keberangkatan.');
  }

  const depBody = (await depRes.json()) as DepartureDetailResponse;

  // Fetch package name for breadcrumb (non-fatal)
  let packageName: string | null = null;
  try {
    const pkgRes = await fetch(`${baseUrl}/v1/packages/${params.id}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    });
    if (pkgRes.ok) {
      const pkgBody = await pkgRes.json() as { package: { name: string } };
      packageName = pkgBody.package.name;
    }
  } catch {
    // non-fatal
  }

  return {
    departure: depBody.departure,
    packageId: params.id,
    packageName
  };
};

export const actions: Actions = {
  default: async ({ request, cookies, params }) => {
    const accessToken = cookies.get('umrohos_access_token');
    if (!accessToken) {
      throw redirect(303, '/console/login');
    }

    const form = await request.formData();
    const departure_date = String(form.get('departure_date') ?? '').trim();
    const return_date = String(form.get('return_date') ?? '').trim();
    const total_seats = Number(form.get('total_seats') ?? 0);
    const status = String(form.get('status') ?? 'draft') as 'open' | 'closed' | 'draft';

    const errors: Record<string, string> = {};
    if (!departure_date) errors.departure_date = 'Tanggal keberangkatan wajib diisi.';
    if (!return_date) errors.return_date = 'Tanggal kembali wajib diisi.';
    if (departure_date && return_date && departure_date >= return_date) {
      errors.return_date = 'Tanggal kembali harus setelah tanggal keberangkatan.';
    }
    if (!total_seats || total_seats < 1) {
      errors.total_seats = 'Kapasitas minimal 1 kursi.';
    }

    if (Object.keys(errors).length > 0) {
      return fail(400, {
        errors,
        values: { departure_date, return_date, total_seats: String(total_seats), status }
      });
    }

    try {
      await updateDeparture(params.dep_id, { departure_date, return_date, total_seats, status }, accessToken);
    } catch (e) {
      if (e instanceof CatalogWriteError) {
        return fail(e.status, {
          errors: {},
          serverError: e.message,
          values: { departure_date, return_date, total_seats: String(total_seats), status }
        });
      }
      return fail(500, {
        errors: {},
        serverError: 'Gagal menyimpan perubahan. Silakan coba lagi.',
        values: { departure_date, return_date, total_seats: String(total_seats), status }
      });
    }

    throw redirect(303, `/console/packages/${params.id}/departures`);
  }
};
