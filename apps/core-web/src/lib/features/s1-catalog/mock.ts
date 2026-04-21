import type { DepartureDetail, PackageCard, PackageDetail } from './types';

const defaultWhatsappHref = 'https://wa.me/6281200000000';

const defaultImportantNotes = [
  'Paspor dengan masa berlaku minimal 6 bulan saat keberangkatan.',
  'Vaksinasi Meningitis & Influenza (disarankan).',
  'Pas foto terbaru ukuran 4x6 dengan latar putih.',
  'Kartu Keluarga (KK) & Buku Nikah bagi pasangan suami istri.'
];

const defaultFaqs = [
  {
    question: 'Apakah tersedia sistem cicilan?',
    answer:
      'Ya, kami bekerja sama dengan berbagai bank syariah dan lembaga pembiayaan resmi untuk program Tabungan Umrah dan Cicilan hingga 24 bulan tanpa bunga (syarat berlaku).'
  },
  {
    question: 'Bagaimana jika dokumen paspor saya belum siap?',
    answer:
      'Kami menyediakan surat rekomendasi resmi untuk pengurusan paspor di kantor imigrasi terdekat. Anda dapat membooking paket terlebih dahulu sambil mengurus dokumen.'
  },
  {
    question: 'Apa kebijakan pembatalan paket?',
    answer:
      'Pembatalan sebelum issued tiket akan dikenakan biaya administrasi minimal. Jika tiket sudah issued, biaya pembatalan mengikuti regulasi maskapai dan hotel. Kami menyarankan asuransi perjalanan.'
  }
];

const packageDetails: PackageDetail[] = [
  {
    id: 'demo-pkg-umrah-12d',
    kind: 'umrah_reguler',
    name: 'Umrah Executive 12 Hari Ramadan',
    description:
      'Perjalanan spiritual yang tenang dengan layanan premium: hotel dekat Masjidil Haram, pendampingan muthawwif, dan jadwal ibadah yang terstruktur.',
    highlights: ['Direct Jakarta-Jeddah', 'Hotel bintang 4 jarak jalan kaki', 'Pendamping muthawwif berpengalaman'],
    coverPhotoUrl: 'https://images.unsplash.com/photo-1591604466107-ec97de577aff?auto=format&fit=crop&w=1200&q=80',
    startingPriceLabel: 'Mulai dari Rp 38.500.000',
    displayPriceShort: 'Rp 38,5 jt',
    primaryBadge: 'Paling Populer',
    secondaryBadges: ['12 Hari', 'Ramadan'],
    priceFinePrint: '*Harga dapat berubah mengikuti kuota dan ketersediaan maskapai.',
    trustPpiu: 'Izin PPIU No. 123/2024',
    ratingScore: '4.9/5',
    ratingReviewsLabel: '(2.4k ulasan)',
    whatsappHref: defaultWhatsappHref,
    galleryPhotoUrls: [
      'https://images.unsplash.com/photo-1591604466107-ec97de577aff?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1564769625905-50e93615e769?auto=format&fit=crop&w=800&q=80',
      'https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?auto=format&fit=crop&w=800&q=80',
      'https://images.unsplash.com/photo-1544739313-3b4e5f9f2f7b?auto=format&fit=crop&w=800&q=80'
    ],
    inclusions: [
      {
        icon: 'apartment',
        title: 'Hotel bintang 4–5',
        description: 'Akses jalan kaki atau shuttle ke Masjidil Haram'
      },
      {
        icon: 'flight',
        title: 'Tiket pesawat PP',
        description: 'Maskapai full service dengan bagasi terjadwal'
      },
      {
        icon: 'restaurant',
        title: 'Makan 3x sehari',
        description: 'Menu halal bervariasi selama di Saudi'
      },
      {
        icon: 'directions_bus',
        title: 'Transportasi',
        description: 'Bus AC dan city tour sesuai itinerary'
      },
      {
        icon: 'mosque',
        title: 'Manasik & pembimbingan',
        description: 'Bimbingan intensif sebelum dan selama perjalanan'
      }
    ],
    importantNotes: defaultImportantNotes,
    itineraryDays: [
      {
        dayLabel: '01',
        title: 'Jakarta — Jeddah — Makkah',
        body: 'Keberangkatan dari Jakarta, proses imigrasi di Jeddah, lalu menuju Makkah untuk ibadah Umrah pertama.'
      },
      {
        dayLabel: '02–05',
        title: 'Ibadah di Makkah',
        body: 'Thawaf, sai, dan ibadah sunnah di Masjidil Haram dengan pendampingan tim lapangan.'
      },
      {
        dayLabel: '06',
        title: 'Ziarah Makkah',
        body: 'Kunjungan ke situs bersejarah di Makkah sesuai jadwal rombongan.'
      },
      {
        dayLabel: '07–12',
        title: 'Makkah — Madinah — Jakarta',
        body: 'Perjalanan ke Madinah, ziarah ke Masjid Nabawi, persiapan pulang ke Tanah Air.'
      }
    ],
    faqs: defaultFaqs,
    facilitiesBody:
      'Paket ini mencakup akomodasi hotel sesuai kontrak rombongan, transport bandara–hotel, tiket pesawat kelas ekonomi (kecuali dinyatakan lain), makan sesuai jadwal, manasik pra-keberangkatan, dan pendampingan muthawwif berizin. Biaya visa, asuransi perjalanan, dan pengeluaran pribadi di luar jadwal rombongan tidak termasuk kecuali disebutkan di kontrak pemesanan.',
    termsSummary:
      'Pemesanan mengikat setelah uang muka dan dokumen identitas valid diterima. Perubahan jadwal atau pembatalan mengikuti kebijakan maskapai, hotel, dan biaya administrasi travel. Syarat kesehatan dan visa mengikuti ketentuan pemerintah Saudi dan Indonesia yang berlaku pada saat keberangkatan.',
    departures: [
      {
        id: 'dep-ramadan-12h-a',
        departureDate: '2026-11-12',
        returnDate: '2026-11-23',
        status: 'open',
        remainingSeats: 12,
        airline: 'Saudia Airlines'
      },
      {
        id: 'dep-ramadan-12h-b',
        departureDate: '2026-12-03',
        returnDate: '2026-12-14',
        status: 'open',
        remainingSeats: 5,
        airline: 'Saudia Airlines'
      },
      {
        id: 'dep-ramadan-12h-c',
        departureDate: '2026-12-18',
        returnDate: '2026-12-29',
        status: 'open',
        remainingSeats: 28,
        airline: 'Garuda Indonesia'
      }
    ]
  },
  {
    id: 'demo-pkg-umrah-9d',
    kind: 'umrah_plus',
    name: 'Umrah Plus Turki 9 Hari',
    description:
      'Paket singkat untuk jamaah yang mengutamakan efisiensi waktu dengan itinerary ringkas dan terstruktur.',
    highlights: ['Transit ringan', 'Jadwal padat namun nyaman', 'Cocok untuk first timer'],
    coverPhotoUrl: 'https://images.unsplash.com/photo-1544739313-3b4e5f9f2f7b?auto=format&fit=crop&w=1200&q=80',
    startingPriceLabel: 'Mulai dari Rp 31.900.000',
    displayPriceShort: 'Rp 31,9 jt',
    primaryBadge: 'Efisien',
    secondaryBadges: ['9 Hari', 'Reguler'],
    priceFinePrint: '*Harga dapat berubah mengikuti kuota dan ketersediaan maskapai.',
    trustPpiu: 'Izin PPIU No. 123/2024',
    ratingScore: '4.8/5',
    ratingReviewsLabel: '(980 ulasan)',
    whatsappHref: defaultWhatsappHref,
    galleryPhotoUrls: [
      'https://images.unsplash.com/photo-1544739313-3b4e5f9f2f7b?auto=format&fit=crop&w=1200&q=80',
      'https://images.unsplash.com/photo-1524231757912-21f3fecef4b8?auto=format&fit=crop&w=800&q=80',
      'https://images.unsplash.com/photo-1591604466107-ec97de577aff?auto=format&fit=crop&w=800&q=80'
    ],
    inclusions: [
      {
        icon: 'flight',
        title: 'Tiket PP',
        description: 'Rute sesuai konfirmasi final dari travel'
      },
      {
        icon: 'hotel',
        title: 'Hotel',
        description: 'Akomodasi sesuai kelas paket'
      },
      {
        icon: 'restaurant',
        title: 'Makan',
        description: 'Sesuai jadwal rombongan'
      },
      {
        icon: 'groups',
        title: 'Pendamping',
        description: 'Koordinasi lapangan'
      }
    ],
    importantNotes: defaultImportantNotes,
    itineraryDays: [
      {
        dayLabel: '01',
        title: 'Keberangkatan — Makkah',
        body: 'Perjalanan menuju Saudi dan persiapan ibadah di Makkah.'
      },
      {
        dayLabel: '02–07',
        title: 'Ibadah & ziarah',
        body: 'Rangkaian ibadah utama dan ziarah sesuai program singkat.'
      },
      {
        dayLabel: '08–09',
        title: 'Kepulangan',
        body: 'Persiapan check-out dan penerbangan kembali ke Indonesia.'
      }
    ],
    faqs: defaultFaqs,
    facilitiesBody:
      'Paket mencakup komponen standar perjalanan umrah sesuai brosur: tiket, hotel, makan terjadwal, dan pendampingan. Detail kamar dan maskapai final mengikuti kontrak pada saat pemesanan.',
    termsSummary:
      'Setuju dengan syarat pembatalan dan perubahan tarif dari pihak ketiga (maskapai, hotel). Dokumen jamaah harus lengkap sebelum deadline yang diberikan travel.',
    departures: [
      {
        id: 'dep-plus-9h-a',
        departureDate: '2026-10-21',
        returnDate: '2026-10-29',
        status: 'open',
        remainingSeats: 18,
        airline: 'Qatar Airways'
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
      { roomType: 'quad', amountLabel: 'Rp 38.500.000', listAmountIdr: 38_500_000 },
      { roomType: 'triple', amountLabel: 'Rp 41.500.000', listAmountIdr: 41_500_000 },
      { roomType: 'double', amountLabel: 'Rp 45.500.000', listAmountIdr: 45_500_000 }
    ]
  },
  {
    id: 'dep-ramadan-12h-b',
    packageId: 'demo-pkg-umrah-12d',
    departureDate: '2026-12-03',
    returnDate: '2026-12-14',
    totalSeats: 45,
    remainingSeats: 5,
    status: 'open',
    pricing: [
      { roomType: 'quad', amountLabel: 'Rp 39.000.000', listAmountIdr: 39_000_000 },
      { roomType: 'triple', amountLabel: 'Rp 42.000.000', listAmountIdr: 42_000_000 },
      { roomType: 'double', amountLabel: 'Rp 46.000.000', listAmountIdr: 46_000_000 }
    ]
  },
  {
    id: 'dep-ramadan-12h-c',
    packageId: 'demo-pkg-umrah-12d',
    departureDate: '2026-12-18',
    returnDate: '2026-12-29',
    totalSeats: 40,
    remainingSeats: 28,
    status: 'open',
    pricing: [
      { roomType: 'quad', amountLabel: 'Rp 38.800.000', listAmountIdr: 38_800_000 },
      { roomType: 'triple', amountLabel: 'Rp 41.800.000', listAmountIdr: 41_800_000 },
      { roomType: 'double', amountLabel: 'Rp 45.800.000', listAmountIdr: 45_800_000 }
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
      { roomType: 'quad', amountLabel: 'Rp 31.900.000', listAmountIdr: 31_900_000 },
      { roomType: 'triple', amountLabel: 'Rp 34.500.000', listAmountIdr: 34_500_000 },
      { roomType: 'double', amountLabel: 'Rp 37.500.000', listAmountIdr: 37_500_000 }
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
