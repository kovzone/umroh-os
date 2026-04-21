import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';
import { getCatalogDepartureDetail, getCatalogPackageDetail } from '$lib/features/s1-catalog/repository';

export const load: PageLoad = async ({ params }) => {
  const dep = await getCatalogDepartureDetail(params.departure_id);
  const pkg = await getCatalogPackageDetail(params.package_id);

  if (!pkg || !dep || dep.packageId !== pkg.id) {
    throw error(404, 'Departure not found');
  }

  return {
    package: pkg,
    departure: dep
  };
};
