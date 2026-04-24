import type { RequestEvent } from '@sveltejs/kit';

const GATEWAY_URL =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export interface ItineraryDay {
  day: number;
  title: string;
  date: string;
  activities: {
    time: string;
    icon: string;
    title: string;
    desc: string;
    type: 'flight' | 'hotel' | 'ibadah' | 'meal' | 'tour' | 'ziarah' | 'transport';
  }[];
  hotel?: {
    name: string;
    stars: number;
    location: string;
  };
}

export interface PackageItinerary {
  id: string;
  name: string;
  duration_days: number;
  departure_date: string;
  return_date: string;
  destinations: string[];
  highlights: string[];
  days: ItineraryDay[];
  guide_name: string;
  airline: string;
  hotel_makkah: string;
  hotel_madinah: string;
  quota: number;
  price: number;
}

export interface PageData {
  itinerary: PackageItinerary | null;
  error: string | null;
  shareUrl: string;
}

function mockItinerary(packageId: string): PackageItinerary {
  const pkg: PackageItinerary = {
    id: packageId,
    name: 'Umroh Ramadhan Premium 2026',
    duration_days: 15,
    departure_date: '2026-03-01',
    return_date: '2026-03-15',
    destinations: ['Jakarta', 'Madinah', 'Makkah'],
    highlights: ['Hotel Bintang 5 di Masjidil Haram', 'Muthawif berpengalaman', 'Visa garansi', 'Ziarah full', 'Makan 3x sehari'],
    guide_name: 'Ustadz Ahmad Fauzi, Lc.',
    airline: 'Garuda Indonesia',
    hotel_makkah: 'Pullman ZamZam Makkah (bintang 5)',
    hotel_madinah: 'Hilton Madinah (bintang 5)',
    quota: 45,
    price: 32_500_000,
    days: [
      {
        day: 1,
        title: 'Keberangkatan dari Jakarta',
        date: '2026-03-01',
        activities: [
          { time: '06:00', icon: 'groups', title: 'Kumpul di Bandara Soekarno-Hatta', desc: 'Check-in dan pengumpulan dokumen. Terminal 3 Domestik.', type: 'transport' },
          { time: '10:30', icon: 'flight_takeoff', title: 'Penerbangan Jakarta → Madinah', desc: 'Garuda Indonesia GA-980. Estimasi 9 jam penerbangan.', type: 'flight' },
          { time: '16:30', icon: 'flight_land', title: 'Tiba di Bandara Madinah (AMAA)', desc: 'Proses imigrasi dan pengambilan bagasi.', type: 'flight' },
          { time: '19:00', icon: 'hotel', title: 'Check-in Hotel Hilton Madinah', desc: 'Hilton Madinah — 800m dari Masjid Nabawi.', type: 'hotel' },
          { time: '20:00', icon: 'restaurant', title: 'Makan malam', desc: 'Buffet Indonesia & Middle Eastern di hotel.', type: 'meal' }
        ],
        hotel: { name: 'Hilton Madinah', stars: 5, location: '800m dari Masjid Nabawi' }
      },
      {
        day: 2,
        title: 'Madinah — Ziarah & Ibadah',
        date: '2026-03-02',
        activities: [
          { time: '04:00', icon: 'mosque', title: 'Shalat Subuh di Masjid Nabawi', desc: 'Shalat berjamaah dan zikir pagi bersama rombongan.', type: 'ibadah' },
          { time: '09:00', icon: 'explore', title: 'Ziarah Raudhah', desc: 'Berziarah ke makam Nabi SAW di Raudhah.', type: 'ziarah' },
          { time: '12:00', icon: 'restaurant', title: 'Makan siang', desc: 'Kembali ke hotel untuk makan siang.', type: 'meal' },
          { time: '14:00', icon: 'location_on', title: 'Ziarah Baqi & sekitar Madinah', desc: 'Mengunjungi pemakaman Baqi, Masjid Quba, dan Masjid Qiblatayn.', type: 'ziarah' },
          { time: '20:30', icon: 'restaurant', title: 'Makan malam', desc: 'Makan malam di hotel.', type: 'meal' }
        ],
        hotel: { name: 'Hilton Madinah', stars: 5, location: '800m dari Masjid Nabawi' }
      },
      {
        day: 3,
        title: 'Madinah — Hari Bebas',
        date: '2026-03-03',
        activities: [
          { time: '04:00', icon: 'mosque', title: 'Shalat Subuh di Masjid Nabawi', desc: '', type: 'ibadah' },
          { time: '09:00', icon: 'shopping_bag', title: 'Belanja di Pasar Madinah', desc: 'Kurma, minyak atar, perlengkapan ibadah, oleh-oleh.', type: 'tour' },
          { time: '14:00', icon: 'self_improvement', title: 'Istirahat & Ibadah Mandiri', desc: 'Waktu bebas untuk ibadah dan istirahat.', type: 'ibadah' },
          { time: '20:30', icon: 'restaurant', title: 'Makan malam bersama', desc: 'Makan malam spesial bersama rombongan.', type: 'meal' }
        ],
        hotel: { name: 'Hilton Madinah', stars: 5, location: '800m dari Masjid Nabawi' }
      },
      {
        day: 5,
        title: 'Perjalanan Madinah → Makkah',
        date: '2026-03-05',
        activities: [
          { time: '07:00', icon: 'restaurant', title: 'Sarapan & persiapan', desc: 'Sarapan di hotel, persiapan perjalanan.', type: 'meal' },
          { time: '09:00', icon: 'directions_bus', title: 'Berangkat ke Makkah', desc: 'Bus AC, perjalanan ±5 jam melewati Bir Ali (miqat).', type: 'transport' },
          { time: '10:30', icon: 'mosque', title: 'Miqat di Bir Ali', desc: 'Memakai ihram dan berniat Umroh.', type: 'ibadah' },
          { time: '15:00', icon: 'hotel', title: 'Tiba di Makkah — Check-in Hotel', desc: 'Pullman ZamZam Makkah, langsung menghadap Masjidil Haram.', type: 'hotel' },
          { time: '17:00', icon: 'star', title: 'Pelaksanaan Umroh Wajib', desc: 'Tawaf 7 putaran, Sai 7 putaran, Tahallul.', type: 'ibadah' }
        ],
        hotel: { name: 'Pullman ZamZam Makkah', stars: 5, location: 'Berhadapan langsung dengan Masjidil Haram' }
      },
      {
        day: 10,
        title: 'Makkah — Ziarah Dalam Kota',
        date: '2026-03-10',
        activities: [
          { time: '08:00', icon: 'explore', title: 'Ziarah Jabal Nur & Hira', desc: 'Mendaki Jabal Nur, melihat Gua Hira tempat turunnya wahyu pertama.', type: 'ziarah' },
          { time: '12:00', icon: 'restaurant', title: 'Makan siang', desc: 'Makan siang di restoran dekat hotel.', type: 'meal' },
          { time: '14:00', icon: 'location_on', title: 'Ziarah Jabal Tsur & Area Bersejarah', desc: 'Mengunjungi tempat-tempat bersejarah di sekitar Makkah.', type: 'ziarah' },
          { time: '21:00', icon: 'mosque', title: 'Shalat Tarawih di Masjidil Haram', desc: 'Shalat Tarawih berjamaah di Masjidil Haram.', type: 'ibadah' }
        ],
        hotel: { name: 'Pullman ZamZam Makkah', stars: 5, location: 'Berhadapan langsung dengan Masjidil Haram' }
      },
      {
        day: 15,
        title: 'Kepulangan ke Jakarta',
        date: '2026-03-15',
        activities: [
          { time: '06:00', icon: 'restaurant', title: 'Sarapan terakhir di Makkah', desc: 'Sarapan dan persiapan kepulangan.', type: 'meal' },
          { time: '09:00', icon: 'directions_bus', title: 'Transfer ke Bandara Jeddah', desc: 'Bus AC dari Makkah ke King Abdulaziz International Airport.', type: 'transport' },
          { time: '14:00', icon: 'flight_takeoff', title: 'Penerbangan Jeddah → Jakarta', desc: 'Garuda Indonesia GA-981. Estimasi 10 jam penerbangan.', type: 'flight' },
          { time: '23:00', icon: 'flight_land', title: 'Tiba di Jakarta', desc: 'Kembali ke tanah air. Alhamdulillah.', type: 'flight' }
        ]
      }
    ]
  };
  return pkg;
}

export const load = async ({ params, fetch }: RequestEvent): Promise<PageData> => {
  const packageId: string = params['package_id'] ?? 'unknown';
  const shareUrl = `${process.env.PUBLIC_APP_URL ?? 'http://localhost:3001'}/console/crm/itinerary/${packageId}`;

  try {
    const res = await fetch(`${GATEWAY_URL}/v1/packages/${packageId}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();
    return {
      itinerary: {
        ...mockItinerary(packageId),
        ...data,
        days: mockItinerary(packageId).days // always use mock days for now
      },
      error: null,
      shareUrl
    };
  } catch {
    // Fallback to mock
    return {
      itinerary: mockItinerary(packageId),
      error: null,
      shareUrl
    };
  }
};
