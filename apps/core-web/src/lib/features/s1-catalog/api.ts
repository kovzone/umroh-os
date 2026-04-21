import type {
  CatalogDepartureDetailResponse,
  CatalogListResponse,
  CatalogPackageDetailResponse,
  DepartureDetail,
  PackageCard,
  PackageDetail
} from './types';

const baseUrl = import.meta.env.VITE_CATALOG_API_BASE_URL ?? import.meta.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

function toPriceLabel(amount: number, currency: string): string {
  const symbol = currency === 'IDR' ? 'Rp' : currency;
  return `${symbol} ${new Intl.NumberFormat('id-ID').format(amount)}`;
}

async function fetchJson<T>(path: string): Promise<T> {
  const response = await fetch(`${baseUrl}${path}`);
  if (!response.ok) {
    throw new Error(`Catalog API request failed (${response.status}) for ${path}`);
  }
  return (await response.json()) as T;
}

export async function fetchCatalogPackages(): Promise<PackageCard[]> {
  const payload = await fetchJson<CatalogListResponse>('/v1/packages');

  return payload.packages.map((pkg) => ({
    id: pkg.id,
    kind: pkg.kind,
    name: pkg.name,
    blurb: pkg.description,
    coverPhotoUrl: pkg.cover_photo_url,
    startingPriceLabel: toPriceLabel(pkg.starting_price.list_amount, pkg.starting_price.list_currency),
    nextDepartureLabel: pkg.next_departure
      ? `${pkg.next_departure.departure_date} s.d. ${pkg.next_departure.return_date}`
      : 'Belum ada keberangkatan terjadwal',
    remainingSeats: pkg.next_departure?.remaining_seats ?? 0
  }));
}

export async function fetchCatalogPackageDetail(packageId: string): Promise<PackageDetail | null> {
  const payload = await fetchJson<CatalogPackageDetailResponse>(`/v1/packages/${packageId}`);
  const pkg = payload.package;

  return {
    id: pkg.id,
    kind: pkg.kind,
    name: pkg.name,
    description: pkg.description,
    highlights: pkg.highlights,
    coverPhotoUrl: pkg.cover_photo_url,
    startingPriceLabel: 'Harga akan ditampilkan dari departure',
    departures: pkg.departures.map((dep) => ({
      id: dep.id,
      departureDate: dep.departure_date,
      returnDate: dep.return_date,
      status: dep.status,
      remainingSeats: dep.remaining_seats
    }))
  };
}

export async function fetchCatalogDepartureDetail(departureId: string): Promise<DepartureDetail | null> {
  const payload = await fetchJson<CatalogDepartureDetailResponse>(`/v1/package-departures/${departureId}`);
  const dep = payload.departure;

  return {
    id: dep.id,
    packageId: dep.package_id,
    departureDate: dep.departure_date,
    returnDate: dep.return_date,
    totalSeats: dep.total_seats,
    remainingSeats: dep.remaining_seats,
    status: dep.status,
    pricing: dep.pricing.map((price) => ({
      roomType: price.room_type,
      amountLabel: toPriceLabel(price.list_amount, price.list_currency)
    }))
  };
}
