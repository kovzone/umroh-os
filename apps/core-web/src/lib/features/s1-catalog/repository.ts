import { getMockDepartureDetail, getMockPackageDetail, listMockPackages } from './mock';
import type { DepartureDetail, PackageCard, PackageDetail } from './types';

const useMockCatalog = true;

// Adapter for S1-L-02 mock-first flow. In S1-L-03 this can switch to
// gateway-backed requests without changing route/page consumers.
export async function getCatalogPackages(): Promise<PackageCard[]> {
  if (useMockCatalog) {
    return listMockPackages();
  }

  return [];
}

export async function getCatalogPackageDetail(packageId: string): Promise<PackageDetail | null> {
  if (useMockCatalog) {
    return getMockPackageDetail(packageId);
  }

  return null;
}

export async function getCatalogDepartureDetail(departureId: string): Promise<DepartureDetail | null> {
  if (useMockCatalog) {
    return getMockDepartureDetail(departureId);
  }

  return null;
}
