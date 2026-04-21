import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';
import { getCatalogPackageDetail } from '$lib/features/s1-catalog/repository';

export const load: PageLoad = async ({ params }) => {
  const pkg = await getCatalogPackageDetail(params.package_id);

  if (!pkg) {
    throw error(404, 'Package not found');
  }

  return {
    package: pkg
  };
};
