import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { createPackage, CatalogWriteError } from '$lib/features/s1-catalog/write-api';

export const load: PageServerLoad = async ({ cookies }) => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }
  return {};
};

export const actions: Actions = {
  default: async ({ request, cookies }) => {
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

    // Client-side validations on server
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
      await createPackage({ name, description, kind, status, starting_price_idr }, accessToken);
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
        serverError: 'Gagal menyimpan paket. Silakan coba lagi.',
        values: { name, description, kind, status, starting_price_idr: priceRaw }
      });
    }

    throw redirect(303, '/console/packages');
  }
};
