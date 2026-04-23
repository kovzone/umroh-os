import { redirect } from '@sveltejs/kit';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

type LoadEvent = { cookies: Cookies; fetch: RequestEvent['fetch']; params: Record<string, string> };

const MOCK = process.env.MOCK_OPS === 'true' || process.env.VITE_MOCK_OPS === 'true';
const baseUrl =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export type BookingStatus = 'registered' | 'dp_paid' | 'paid' | 'cancelled';
export type DocStatus = 'incomplete' | 'complete' | 'verified';
export type RoomType = 'double' | 'triple' | 'quad';

export interface ManifestEntry {
  no: number;
  id: string;
  name: string;
  nik: string;
  phone: string;
  room_type: RoomType;
  booking_status: BookingStatus;
  doc_status: DocStatus;
}

export interface DepartureInfo {
  id: string;
  package_name: string;
  departure_date: string;
  return_date: string;
  total_seats: number;
  booked_seats: number;
}

export interface PageData {
  departure: DepartureInfo;
  manifest: ManifestEntry[];
  error: string | null;
}

function mockData(departureId: string): PageData {
  const manifest: ManifestEntry[] = [
    { no: 1, id: 'j1', name: 'Ahmad Fauzi', nik: '3271010101800001', phone: '08123456789', room_type: 'double', booking_status: 'paid', doc_status: 'verified' },
    { no: 2, id: 'j2', name: 'Siti Aminah', nik: '3271010101850002', phone: '08234567890', room_type: 'double', booking_status: 'paid', doc_status: 'complete' },
    { no: 3, id: 'j3', name: 'Budi Santoso', nik: '3271010101780003', phone: '08345678901', room_type: 'triple', booking_status: 'dp_paid', doc_status: 'incomplete' },
    { no: 4, id: 'j4', name: 'Dewi Kurniawati', nik: '3271010101900004', phone: '08456789012', room_type: 'triple', booking_status: 'dp_paid', doc_status: 'incomplete' },
    { no: 5, id: 'j5', name: 'Rizky Maulana', nik: '3271010101950005', phone: '08567890123', room_type: 'quad', booking_status: 'registered', doc_status: 'incomplete' },
    { no: 6, id: 'j6', name: 'Nur Hidayah', nik: '3271010101880006', phone: '08678901234', room_type: 'double', booking_status: 'paid', doc_status: 'verified' },
    { no: 7, id: 'j7', name: 'Hasan Basri', nik: '3271010101750007', phone: '08789012345', room_type: 'triple', booking_status: 'paid', doc_status: 'complete' },
    { no: 8, id: 'j8', name: 'Fatimah Zahra', nik: '3271010101920008', phone: '08890123456', room_type: 'quad', booking_status: 'dp_paid', doc_status: 'incomplete' }
  ];

  return {
    departure: {
      id: departureId,
      package_name: 'Umroh Ramadhan Premium',
      departure_date: '2026-03-01',
      return_date: '2026-03-14',
      total_seats: 45,
      booked_seats: manifest.length
    },
    manifest,
    error: null
  };
}

export const load = async ({ cookies, fetch, params }: LoadEvent): Promise<PageData> => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  const { departure_id } = params;

  if (MOCK) {
    return mockData(departure_id);
  }

  try {
    const headers = { Authorization: `Bearer ${accessToken}` };
    const [depRes, manifestRes] = await Promise.all([
      fetch(`${baseUrl}/v1/departures/${departure_id}`, { headers }),
      fetch(`${baseUrl}/v1/departures/${departure_id}/manifest`, { headers })
    ]);

    const departure = depRes.ok
      ? ((await depRes.json()) as DepartureInfo)
      : { id: departure_id, package_name: '-', departure_date: '', return_date: '', total_seats: 0, booked_seats: 0 };

    const manifest = manifestRes.ok ? ((await manifestRes.json()) as ManifestEntry[]) : [];

    return { departure, manifest, error: depRes.ok ? null : `Gagal memuat data (${depRes.status})` };
  } catch {
    return { ...mockData(departure_id), error: 'Tidak dapat terhubung ke gateway.' };
  }
};
