import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { createDeparture, CatalogWriteError } from '$lib/features/s1-catalog/write-api';

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

  // Fetch package name for breadcrumb
  try {
    const res = await fetch(`${baseUrl}/v1/packages/${params.id}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    });
    if (res.ok) {
      const body = await res.json() as { package: { id: string; name: string } };
      return { packageId: params.id, packageName: body.package.name };
    }
  } catch {
    // non-fatal
  }

  return { packageId: params.id, packageName: null };
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
      await createDeparture(params.id, { departure_date, return_date, total_seats, status }, accessToken);
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
        serverError: 'Gagal menyimpan keberangkatan. Silakan coba lagi.',
        values: { departure_date, return_date, total_seats: String(total_seats), status }
      });
    }

    throw redirect(303, `/console/packages/${params.id}/departures`);
  }
};
