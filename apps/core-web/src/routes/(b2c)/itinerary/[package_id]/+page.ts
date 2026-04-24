import type { PageLoad } from './$types';
import { getCatalogPackageDetail } from '$lib/features/s1-catalog/repository';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params }) => {
  const pkg = await getCatalogPackageDetail(params.package_id);
  if (!pkg) {
    throw error(404, 'Paket tidak ditemukan');
  }
  return { package: pkg };
};
