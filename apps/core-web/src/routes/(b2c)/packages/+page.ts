import type { PageLoad } from './$types';
import { getCatalogPackages } from '$lib/features/s1-catalog/repository';

export const load: PageLoad = async () => {
  const packages = await getCatalogPackages();

  return {
    packages
  };
};
