import type { DepartureDetail, PackageCard, PackageDetail } from './types';

const packageDetails: PackageDetail[] = [
  {
    id: 'demo-pkg-umrah-12d',
    kind: 'umrah_reguler',
    name: 'Umrah Executive 12 Hari Ramadan',
    description:
      'Paket umrah 12 hari dengan fokus kenyamanan, hotel dekat Masjidil Haram, dan pendampingan ibadah intensif.',
    highlights: ['Direct Jakarta-Jeddah', 'Hotel bintang 4 jarak jalan kaki', 'Pendamping muthawwif berpengalaman'],
    coverPhotoUrl: 'https://images.unsplash.com/photo-1591604466107-ec97de577aff?auto=format&fit=crop&w=1200&q=80',
    startingPriceLabel: 'Mulai dari Rp 38.500.000',
    departures: [
      {
        id: 'dep-ramadan-12h-a',
        departureDate: '2026-11-12',
        returnDate: '2026-11-23',
        status: 'open',
        remainingSeats: 12
      },
      {
        id: 'dep-ramadan-12h-b',
        departureDate: '2026-12-03',
        returnDate: '2026-12-14',
        status: 'open',
        remainingSeats: 7
      }
    ]
  },
  {
    id: 'demo-pkg-umrah-9d',
    kind: 'umrah_plus',
    name: 'Umrah Plus Turki 9 Hari',
    description:
      'Paket umrah singkat untuk jamaah yang mengutamakan efisiensi waktu dengan itinerary ringkas dan terstruktur.',
    highlights: ['Transit ringan', 'Jadwal padat namun nyaman', 'Cocok untuk first timer'],
    coverPhotoUrl: 'https://images.unsplash.com/photo-1544739313-3b4e5f9f2f7b?auto=format&fit=crop&w=1200&q=80',
    startingPriceLabel: 'Mulai dari Rp 31.900.000',
    departures: [
      {
        id: 'dep-plus-9h-a',
        departureDate: '2026-10-21',
        returnDate: '2026-10-29',
        status: 'open',
        remainingSeats: 18
      }
    ]
  }
];

const departureDetails: DepartureDetail[] = [
  {
    id: 'dep-ramadan-12h-a',
    packageId: 'demo-pkg-umrah-12d',
    departureDate: '2026-11-12',
    returnDate: '2026-11-23',
    totalSeats: 45,
    remainingSeats: 12,
    status: 'open',
    pricing: [
      { roomType: 'quad', amountLabel: 'Rp 38.500.000' },
      { roomType: 'triple', amountLabel: 'Rp 41.500.000' },
      { roomType: 'double', amountLabel: 'Rp 45.500.000' }
    ]
  },
  {
    id: 'dep-ramadan-12h-b',
    packageId: 'demo-pkg-umrah-12d',
    departureDate: '2026-12-03',
    returnDate: '2026-12-14',
    totalSeats: 45,
    remainingSeats: 7,
    status: 'open',
    pricing: [
      { roomType: 'quad', amountLabel: 'Rp 39.000.000' },
      { roomType: 'triple', amountLabel: 'Rp 42.000.000' },
      { roomType: 'double', amountLabel: 'Rp 46.000.000' }
    ]
  },
  {
    id: 'dep-plus-9h-a',
    packageId: 'demo-pkg-umrah-9d',
    departureDate: '2026-10-21',
    returnDate: '2026-10-29',
    totalSeats: 60,
    remainingSeats: 18,
    status: 'open',
    pricing: [
      { roomType: 'quad', amountLabel: 'Rp 31.900.000' },
      { roomType: 'triple', amountLabel: 'Rp 34.500.000' },
      { roomType: 'double', amountLabel: 'Rp 37.500.000' }
    ]
  }
];

export function listMockPackages(): PackageCard[] {
  return packageDetails.map((pkg) => ({
    id: pkg.id,
    kind: pkg.kind,
    name: pkg.name,
    blurb: pkg.description,
    coverPhotoUrl: pkg.coverPhotoUrl,
    startingPriceLabel: pkg.startingPriceLabel,
    nextDepartureLabel: `${pkg.departures[0]?.departureDate ?? '-'} s.d. ${pkg.departures[0]?.returnDate ?? '-'}`,
    remainingSeats: pkg.departures[0]?.remainingSeats ?? 0
  }));
}

export function getMockPackageDetail(packageId: string): PackageDetail | null {
  return packageDetails.find((pkg) => pkg.id === packageId) ?? null;
}

export function getMockDepartureDetail(departureId: string): DepartureDetail | null {
  return departureDetails.find((dep) => dep.id === departureId) ?? null;
}
