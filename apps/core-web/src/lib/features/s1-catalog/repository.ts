import { fetchCatalogDepartureDetail, fetchCatalogPackageDetail, fetchCatalogPackages } from './api';
import { getMockDepartureDetail, getMockPackageDetail, listMockPackages } from './mock';
import type { DepartureDetail, PackageCard, PackageDetail } from './types';

const useMockCatalog = (import.meta.env.VITE_USE_CATALOG_MOCK ?? 'true') === 'true';

// Adapter for S1-L-02 mock-first flow. In S1-L-03 this can switch to
// gateway-backed requests without changing route/page consumers.
export async function getCatalogPackages(): Promise<PackageCard[]> {
  if (useMockCatalog) {
    return listMockPackages();
  }

  try {
    return await fetchCatalogPackages();
  } catch {
    return listMockPackages();
  }
}

export async function getCatalogPackageDetail(packageId: string): Promise<PackageDetail | null> {
  if (useMockCatalog) {
    return getMockPackageDetail(packageId);
  }

  try {
    return await fetchCatalogPackageDetail(packageId);
  } catch {
    return getMockPackageDetail(packageId);
  }
}

export async function getCatalogDepartureDetail(departureId: string): Promise<DepartureDetail | null> {
  if (useMockCatalog) {
    return getMockDepartureDetail(departureId);
  }

  try {
    return await fetchCatalogDepartureDetail(departureId);
  } catch {
    return getMockDepartureDetail(departureId);
  }
}
