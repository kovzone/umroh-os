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

// ---- BL-DASH-001: Vendor Execution Readiness ----
export interface VendorReadiness {
  vendor_id: string;
  vendor_name: string;
  departure: string;
  departure_date: string;
  checklist_done: number;
  checklist_total: number;
}

// ---- BL-DASH-002: Seat Availability ----
export interface SeatAvailability {
  package_id: string;
  package_name: string;
  seats_sold: number;
  seats_available: number;
  seats_total: number;
}

// ---- BL-DASH-003: Cash Flow ----
export interface CashFlowDay {
  label: string; // "Sen", "Sel", etc.
  cash_in: number;
  cash_out: number;
}

export interface CashFlowSummary {
  period: 'daily' | 'weekly';
  data: CashFlowDay[];
  net_today: number;
  net_week: number;
}

// ---- BL-DASH-004: Executive Financial Report ----
export interface FinancialReport {
  revenue_month: number;
  cogs_month: number;
  gross_profit: number;
  gross_margin_pct: number;
  opex_month: number;
  net_profit: number;
  net_margin_pct: number;
  total_assets: number;
  total_liabilities: number;
  equity: number;
}

// ---- BL-DASH-006: Ad Budget Monitor ----
export interface AdCampaign {
  campaign: string;
  channel: string;
  spend: number;
  budget: number;
  closings: number;
  cpl: number;
  cpa: number;
}

// ---- BL-DASH-007: CS Performance Board ----
export interface CsPerformance {
  cs_name: string;
  leads_handled: number;
  converted: number;
  conversion_pct: number;
  avg_response_min: number;
  sla_breach: number;
}

// ---- BL-DASH-008: Live Bus Radar ----
export interface BusStatus {
  bus_id: string;
  plate: string;
  route: string;
  status: 'on_route' | 'idle' | 'maintenance' | 'breakdown';
  last_location: string;
  pax: number;
  driver: string;
}

// ---- BL-DASH-009: Raudhah Status ----
export interface RaudhahStatus {
  departure_id: string;
  departure_name: string;
  total_jemaah: number;
  entered_raudhah: number;
}

// ---- BL-DASH-010: Luggage Tracking ----
export interface LuggageStatus {
  departure_id: string;
  departure_name: string;
  total_bags: number;
  at_origin: number;
  in_transit: number;
  delivered: number;
}

// ---- BL-DASH-011: Incident Report Feed ----
export interface IncidentReport {
  id: string;
  title: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  category: string;
  reported_by: string;
  reported_at: string;
  status: 'open' | 'investigating' | 'resolved';
}

// ---- BL-DASH-012: Warehouse Health ----
export interface WarehouseHealth {
  total_stock_value: number;
  items_critical: number;
  items_reorder: number;
  items_ok: number;
  categories: { name: string; value: number; status: 'ok' | 'reorder' | 'critical' }[];
}

// ---- BL-DASH-013: Logistics Execution Monitor ----
export interface LogisticsExecution {
  paid_unshipped_count: number;
  paid_unshipped_aging_avg_days: number;
  grn_backlog: number;
  po_backlog: number;
  aging_buckets: { label: string; count: number; color: string }[];
}

// ---- BL-DASH-014: Liquidity AR/AP Aging ----
export interface ArApAging {
  ar_current: number;
  ar_30: number;
  ar_60: number;
  ar_90plus: number;
  ap_current: number;
  ap_30: number;
  ap_60: number;
  ap_90plus: number;
  alerts: { type: 'ar' | 'ap'; message: string; amount: number }[];
}

// ---- BL-DASH-015: Inventory Health ----
export interface InventoryHealth {
  total_skus: number;
  healthy_pct: number;
  low_stock_pct: number;
  out_of_stock_pct: number;
  top_items: { name: string; qty: number; status: 'healthy' | 'low' | 'out' }[];
}

// ---- BL-DASH-016: Fulfillment & PO Monitor ----
export interface FulfillmentMonitor {
  open_pos: number;
  overdue_pos: number;
  pending_fulfillments: number;
  overdue_fulfillments: number;
  backlog_items: { po_number: string; vendor: string; due_date: string; status: string; amount: number }[];
}

// ---- BL-DASH-017: Damage Report ----
export interface DamageReport {
  id: string;
  item_name: string;
  qty_damaged: number;
  location: string;
  reported_by: string;
  reported_at: string;
  estimated_loss: number;
  status: 'pending' | 'reviewed' | 'written_off';
}

export interface PageData {
  kpi: KpiData;
  bookingByStatus: BookingStatusCount[];
  revenueChart: MonthlyRevenue[];
  upcomingDepartures: UpcomingDeparture[];
  recentLeads: RecentLead[];
  // Wave 6 widgets
  vendorReadiness: VendorReadiness[];          // BL-DASH-001
  seatAvailability: SeatAvailability[];        // BL-DASH-002
  cashFlow: CashFlowSummary;                   // BL-DASH-003
  financialReport: FinancialReport;            // BL-DASH-004
  adCampaigns: AdCampaign[];                   // BL-DASH-006
  csPerformance: CsPerformance[];              // BL-DASH-007
  busRadar: BusStatus[];                       // BL-DASH-008
  raudhahStatus: RaudhahStatus[];              // BL-DASH-009
  luggageTracking: LuggageStatus[];            // BL-DASH-010
  incidentFeed: IncidentReport[];              // BL-DASH-011
  warehouseHealth: WarehouseHealth;            // BL-DASH-012
  logisticsExecution: LogisticsExecution;      // BL-DASH-013
  arApAging: ArApAging;                        // BL-DASH-014
  inventoryHealth: InventoryHealth;            // BL-DASH-015
  fulfillmentMonitor: FulfillmentMonitor;      // BL-DASH-016
  damageReports: DamageReport[];               // BL-DASH-017
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

    // ---- BL-DASH-001: Vendor Readiness ----
    vendorReadiness: [
      { vendor_id: 'v1', vendor_name: 'PT Makkah Tours', departure: 'Umroh Ramadhan Premium', departure_date: '2026-03-01', checklist_done: 14, checklist_total: 16 },
      { vendor_id: 'v2', vendor_name: 'PT Madinah Express', departure: 'Umroh Reguler April', departure_date: '2026-04-10', checklist_done: 9, checklist_total: 16 },
      { vendor_id: 'v3', vendor_name: 'PT Haramain Travel', departure: 'Umroh Plus Turki', departure_date: '2026-05-05', checklist_done: 5, checklist_total: 16 },
      { vendor_id: 'v4', vendor_name: 'PT Nusantara Haji', departure: 'Umroh Hemat Juni', departure_date: '2026-06-15', checklist_done: 2, checklist_total: 16 }
    ],

    // ---- BL-DASH-002: Seat Availability ----
    seatAvailability: [
      { package_id: 'p1', package_name: 'Umroh Ramadhan Premium', seats_sold: 38, seats_available: 7, seats_total: 45 },
      { package_id: 'p2', package_name: 'Umroh Reguler April', seats_sold: 42, seats_available: 8, seats_total: 50 },
      { package_id: 'p3', package_name: 'Umroh Plus Turki', seats_sold: 15, seats_available: 25, seats_total: 40 },
      { package_id: 'p4', package_name: 'Umroh Hemat Juni', seats_sold: 10, seats_available: 40, seats_total: 50 },
      { package_id: 'p5', package_name: 'Umroh VIP Juli', seats_sold: 5, seats_available: 20, seats_total: 25 }
    ],

    // ---- BL-DASH-003: Cash Flow ----
    cashFlow: {
      period: 'weekly',
      data: [
        { label: 'Sen', cash_in: 185_000_000, cash_out: 72_000_000 },
        { label: 'Sel', cash_in: 210_000_000, cash_out: 95_000_000 },
        { label: 'Rab', cash_in: 148_000_000, cash_out: 125_000_000 },
        { label: 'Kam', cash_in: 230_000_000, cash_out: 88_000_000 },
        { label: 'Jum', cash_in: 175_000_000, cash_out: 140_000_000 },
        { label: 'Sab', cash_in: 90_000_000, cash_out: 45_000_000 },
        { label: 'Min', cash_in: 55_000_000, cash_out: 20_000_000 }
      ],
      net_today: 35_000_000,
      net_week: 508_000_000
    },

    // ---- BL-DASH-004: Executive Financial Report ----
    financialReport: {
      revenue_month: 1_342_500_000,
      cogs_month: 804_000_000,
      gross_profit: 538_500_000,
      gross_margin_pct: 40.1,
      opex_month: 215_000_000,
      net_profit: 323_500_000,
      net_margin_pct: 24.1,
      total_assets: 8_750_000_000,
      total_liabilities: 3_120_000_000,
      equity: 5_630_000_000
    },

    // ---- BL-DASH-006: Ad Budget Monitor ----
    adCampaigns: [
      { campaign: 'Ramadhan Promo Q1', channel: 'Meta Ads', spend: 12_500_000, budget: 20_000_000, closings: 8, cpl: 312_500, cpa: 1_562_500 },
      { campaign: 'Google Search April', channel: 'Google Ads', spend: 8_200_000, budget: 15_000_000, closings: 5, cpl: 410_000, cpa: 1_640_000 },
      { campaign: 'TikTok Awareness', channel: 'TikTok Ads', spend: 5_100_000, budget: 8_000_000, closings: 3, cpl: 510_000, cpa: 1_700_000 },
      { campaign: 'Retargeting Juni', channel: 'Meta Ads', spend: 3_800_000, budget: 5_000_000, closings: 4, cpl: 237_500, cpa: 950_000 }
    ],

    // ---- BL-DASH-007: CS Performance Board ----
    csPerformance: [
      { cs_name: 'Siti Rahayu', leads_handled: 38, converted: 14, conversion_pct: 36.8, avg_response_min: 8, sla_breach: 2 },
      { cs_name: 'Dewi Kusuma', leads_handled: 31, converted: 12, conversion_pct: 38.7, avg_response_min: 11, sla_breach: 3 },
      { cs_name: 'Ahmad Fauzi', leads_handled: 28, converted: 9, conversion_pct: 32.1, avg_response_min: 15, sla_breach: 5 },
      { cs_name: 'Rizka Aulia', leads_handled: 22, converted: 10, conversion_pct: 45.5, avg_response_min: 7, sla_breach: 1 }
    ],

    // ---- BL-DASH-008: Live Bus Radar ----
    busRadar: [
      { bus_id: 'b1', plate: 'B 1234 TRV', route: 'Jakarta → Bandara Soetta', status: 'on_route', last_location: 'Tol Sedyatmo KM 22', pax: 42, driver: 'Pak Hendra' },
      { bus_id: 'b2', plate: 'B 5678 TRV', route: 'Bekasi → Bandara Soetta', status: 'on_route', last_location: 'Tol Jakarta-Cikampek KM 8', pax: 38, driver: 'Pak Darmawan' },
      { bus_id: 'b3', plate: 'D 9012 TRV', route: 'Bandung → Bandara Husein', status: 'idle', last_location: 'Pool Bandung', pax: 0, driver: 'Pak Rudi' },
      { bus_id: 'b4', plate: 'B 3456 TRV', route: 'Tangerang → Bandara Soetta', status: 'maintenance', last_location: 'Bengkel Ciputat', pax: 0, driver: '-' },
      { bus_id: 'b5', plate: 'B 7890 TRV', route: 'Depok → Bandara Soetta', status: 'breakdown', last_location: 'Tol Depok-Antasari KM 3', pax: 0, driver: 'Pak Santoso' }
    ],

    // ---- BL-DASH-009: Raudhah Status ----
    raudhahStatus: [
      { departure_id: 'dep1', departure_name: 'Umroh Ramadhan Premium (Mar 2026)', total_jemaah: 45, entered_raudhah: 42 },
      { departure_id: 'dep2', departure_name: 'Umroh Reguler April', total_jemaah: 50, entered_raudhah: 31 },
      { departure_id: 'dep3', departure_name: 'Umroh Plus Turki (Mei)', total_jemaah: 40, entered_raudhah: 12 }
    ],

    // ---- BL-DASH-010: Luggage Tracking ----
    luggageTracking: [
      { departure_id: 'dep1', departure_name: 'Umroh Ramadhan Premium', total_bags: 90, at_origin: 0, in_transit: 5, delivered: 85 },
      { departure_id: 'dep2', departure_name: 'Umroh Reguler April', total_bags: 100, at_origin: 15, in_transit: 45, delivered: 40 },
      { departure_id: 'dep3', departure_name: 'Umroh Plus Turki', total_bags: 80, at_origin: 60, in_transit: 18, delivered: 2 }
    ],

    // ---- BL-DASH-011: Incident Feed ----
    incidentFeed: [
      { id: 'inc1', title: 'Bus B 7890 TRV mogok di tol Depok', severity: 'critical', category: 'Transportasi', reported_by: 'Pak Santoso', reported_at: new Date(Date.now() - 3600000).toISOString(), status: 'investigating' },
      { id: 'inc2', title: 'Koper jemaah hilang di baggage claim', severity: 'high', category: 'Koper', reported_by: 'Pembimbing Dep2', reported_at: new Date(Date.now() - 7200000).toISOString(), status: 'open' },
      { id: 'inc3', title: 'Jemaah sakit ringan di hotel Madinah', severity: 'medium', category: 'Kesehatan', reported_by: 'Muthowif Dep1', reported_at: new Date(Date.now() - 18000000).toISOString(), status: 'resolved' },
      { id: 'inc4', title: 'Dokumen visa terlambat 1 jemaah', severity: 'medium', category: 'Dokumen', reported_by: 'Tim Visa', reported_at: new Date(Date.now() - 86400000).toISOString(), status: 'resolved' },
      { id: 'inc5', title: 'Keterlambatan penerbangan 2 jam', severity: 'low', category: 'Penerbangan', reported_by: 'CS Ahmad', reported_at: new Date(Date.now() - 172800000).toISOString(), status: 'resolved' }
    ],

    // ---- BL-DASH-012: Warehouse Health ----
    warehouseHealth: {
      total_stock_value: 485_000_000,
      items_critical: 4,
      items_reorder: 11,
      items_ok: 87,
      categories: [
        { name: 'Koper & Tas', value: 185_000_000, status: 'ok' },
        { name: 'Perlengkapan Ihram', value: 92_000_000, status: 'reorder' },
        { name: 'Kit Jemaah', value: 78_000_000, status: 'critical' },
        { name: 'Konsumsi Perjalanan', value: 65_000_000, status: 'reorder' },
        { name: 'Souvenir & Oleh-oleh', value: 65_000_000, status: 'ok' }
      ]
    },

    // ---- BL-DASH-013: Logistics Execution Monitor ----
    logisticsExecution: {
      paid_unshipped_count: 18,
      paid_unshipped_aging_avg_days: 4.2,
      grn_backlog: 7,
      po_backlog: 12,
      aging_buckets: [
        { label: '1-3 hari', count: 8, color: '#6ee7b7' },
        { label: '4-7 hari', count: 6, color: '#fde68a' },
        { label: '8-14 hari', count: 3, color: '#fca5a5' },
        { label: '>14 hari', count: 1, color: '#f87171' }
      ]
    },

    // ---- BL-DASH-014: AR/AP Aging ----
    arApAging: {
      ar_current: 320_000_000,
      ar_30: 185_000_000,
      ar_60: 72_000_000,
      ar_90plus: 28_000_000,
      ap_current: 210_000_000,
      ap_30: 95_000_000,
      ap_60: 42_000_000,
      ap_90plus: 15_000_000,
      alerts: [
        { type: 'ar', message: 'PT Wisata Nusantara belum bayar 90+ hari', amount: 18_500_000 },
        { type: 'ap', message: 'Tagihan vendor Hotel Makkah jatuh tempo 3 hari', amount: 42_000_000 }
      ]
    },

    // ---- BL-DASH-015: Inventory Health ----
    inventoryHealth: {
      total_skus: 102,
      healthy_pct: 68,
      low_stock_pct: 22,
      out_of_stock_pct: 10,
      top_items: [
        { name: 'Koper 24 inch (Hijau)', qty: 45, status: 'healthy' },
        { name: 'Mukena Putih Premium', qty: 12, status: 'low' },
        { name: 'Buku Doa Umroh', qty: 3, status: 'low' },
        { name: 'Kain Ihram Standar', qty: 0, status: 'out' },
        { name: 'Kantong Koper Hitam', qty: 28, status: 'healthy' },
        { name: 'Gelang Jemaah RFID', qty: 7, status: 'low' }
      ]
    },

    // ---- BL-DASH-016: Fulfillment & PO Monitor ----
    fulfillmentMonitor: {
      open_pos: 14,
      overdue_pos: 3,
      pending_fulfillments: 9,
      overdue_fulfillments: 2,
      backlog_items: [
        { po_number: 'PO-2026-0412', vendor: 'PT Koper Indo', due_date: '2026-04-20', status: 'overdue', amount: 28_000_000 },
        { po_number: 'PO-2026-0418', vendor: 'PT Ihram Jaya', due_date: '2026-04-25', status: 'overdue', amount: 15_500_000 },
        { po_number: 'PO-2026-0421', vendor: 'PT Souvenir Haji', due_date: '2026-04-28', status: 'pending', amount: 9_800_000 },
        { po_number: 'PO-2026-0422', vendor: 'PT Katering Perjalanan', due_date: '2026-04-30', status: 'pending', amount: 22_300_000 },
        { po_number: 'PO-2026-0425', vendor: 'PT Mukena Premium', due_date: '2026-05-03', status: 'pending', amount: 11_200_000 }
      ]
    },

    // ---- BL-DASH-017: Damage Report ----
    damageReports: [
      { id: 'dmg1', item_name: 'Koper 24 inch (Merah)', qty_damaged: 3, location: 'Gudang Utama', reported_by: 'Tim Gudang', reported_at: new Date(Date.now() - 86400000).toISOString(), estimated_loss: 2_400_000, status: 'reviewed' },
      { id: 'dmg2', item_name: 'Buku Doa Umroh (basah)', qty_damaged: 25, location: 'Gudang B', reported_by: 'Supervisor Gudang', reported_at: new Date(Date.now() - 172800000).toISOString(), estimated_loss: 625_000, status: 'written_off' },
      { id: 'dmg3', item_name: 'Gelang RFID (rusak sensor)', qty_damaged: 8, location: 'Ruang Picking', reported_by: 'Tim IT', reported_at: new Date(Date.now() - 259200000).toISOString(), estimated_loss: 3_200_000, status: 'pending' },
      { id: 'dmg4', item_name: 'Mukena Putih (noda permanen)', qty_damaged: 6, location: 'Gudang Utama', reported_by: 'Tim Gudang', reported_at: new Date(Date.now() - 345600000).toISOString(), estimated_loss: 900_000, status: 'pending' }
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
