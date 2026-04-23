import { redirect } from '@sveltejs/kit';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

type LoadEvent = { cookies: Cookies; fetch: RequestEvent['fetch'] };

const baseUrl =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export interface KpiData {
  total_bookings_month: number;
  total_revenue_month: number;
  seats_sold: number;
  leads_new_month: number;
}

export interface BookingStatusCount {
  status: string;
  label: string;
  count: number;
  color: string;
}

export interface MonthlyRevenue {
  month: string; // "Jan", "Feb", etc.
  revenue: number;
}

export interface UpcomingDeparture {
  id: string;
  package_name: string;
  departure_date: string;
  seats_remaining: number;
  total_seats: number;
}

export interface RecentLead {
  id: string;
  name: string;
  status: string;
  cs_name: string;
  created_at: string;
}

export interface PageData {
  kpi: KpiData;
  bookingByStatus: BookingStatusCount[];
  revenueChart: MonthlyRevenue[];
  upcomingDepartures: UpcomingDeparture[];
  recentLeads: RecentLead[];
  error: string | null;
}

function mockData(): PageData {
  return {
    kpi: {
      total_bookings_month: 47,
      total_revenue_month: 1_342_500_000,
      seats_sold: 183,
      leads_new_month: 89
    },
    bookingByStatus: [
      { status: 'registered', label: 'Terdaftar', count: 12, color: '#93c5fd' },
      { status: 'dp_paid', label: 'Sudah DP', count: 18, color: '#fde68a' },
      { status: 'paid', label: 'Lunas', count: 14, color: '#6ee7b7' },
      { status: 'cancelled', label: 'Batal', count: 3, color: '#fca5a5' }
    ],
    revenueChart: [
      { month: 'Nov', revenue: 820_000_000 },
      { month: 'Des', revenue: 1_150_000_000 },
      { month: 'Jan', revenue: 980_000_000 },
      { month: 'Feb', revenue: 1_050_000_000 },
      { month: 'Mar', revenue: 1_280_000_000 },
      { month: 'Apr', revenue: 1_342_500_000 }
    ],
    upcomingDepartures: [
      { id: 'dep1', package_name: 'Umroh Ramadhan Premium', departure_date: '2026-03-01', seats_remaining: 7, total_seats: 45 },
      { id: 'dep2', package_name: 'Umroh Reguler April', departure_date: '2026-04-10', seats_remaining: 8, total_seats: 50 },
      { id: 'dep3', package_name: 'Umroh Plus Turki', departure_date: '2026-05-05', seats_remaining: 25, total_seats: 40 },
      { id: 'dep4', package_name: 'Umroh Hemat Juni', departure_date: '2026-06-15', seats_remaining: 40, total_seats: 50 },
      { id: 'dep5', package_name: 'Umroh VIP Juli', departure_date: '2026-07-01', seats_remaining: 20, total_seats: 25 }
    ],
    recentLeads: [
      { id: 'l1', name: 'Hendra Wijaya', status: 'new', cs_name: 'Siti Rahayu', created_at: new Date(Date.now() - 1800000).toISOString() },
      { id: 'l2', name: 'Nur Aini', status: 'contacted', cs_name: 'Dewi Kusuma', created_at: new Date(Date.now() - 3600000).toISOString() },
      { id: 'l3', name: 'Agus Susanto', status: 'qualified', cs_name: 'Siti Rahayu', created_at: new Date(Date.now() - 7200000).toISOString() },
      { id: 'l4', name: 'Rina Marlina', status: 'new', cs_name: 'Dewi Kusuma', created_at: new Date(Date.now() - 10800000).toISOString() },
      { id: 'l5', name: 'Dedi Kurniawan', status: 'contacted', cs_name: 'Siti Rahayu', created_at: new Date(Date.now() - 14400000).toISOString() }
    ],
    error: null
  };
}

export const load = async ({ cookies, fetch }: LoadEvent): Promise<PageData> => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  // Dashboard always uses mock/aggregated data for now
  // In production, aggregate from multiple APIs
  void fetch;
  void baseUrl;

  return mockData();
};
