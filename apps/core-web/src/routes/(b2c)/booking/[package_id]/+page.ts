import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';
import {
  getCatalogDepartureDetail,
  getCatalogPackageDetail
} from '$lib/features/s1-catalog/repository';
import type { DepartureDetail } from '$lib/features/s1-catalog/types';

export const load: PageLoad = async ({ params, url }) => {
  const pkg = await getCatalogPackageDetail(params.package_id);

  if (!pkg) {
    throw error(404, 'Package not found');
  }

  const preferredDepartureId = url.searchParams.get('departure') ?? undefined;

  const pairs = await Promise.all(
    pkg.departures.map(async (d) => {
      const detail = await getCatalogDepartureDetail(d.id);
      return [d.id, detail] as const;
    })
  );

  const departureDetailById: Record<string, DepartureDetail> = {};
  for (const [id, detail] of pairs) {
    if (detail) {
      departureDetailById[id] = detail;
    }
  }

  return {
    package: pkg,
    preferredDepartureId,
    departureDetailById
  };
};
