// -----------------------------------------------------------------------
// S4-CRM — Mock Data (enabled when VITE_MOCK_CRM=true)
// -----------------------------------------------------------------------
import type {
  CSUser,
  CreateLeadRequest,
  Lead,
  LeadFilters,
  LeadList,
  PackageOption,
  UpdateLeadRequest
} from './types';

// --------------- static reference data --------------------------------

export const MOCK_CS_USERS: CSUser[] = [
  { id: 'cs_001', name: 'Rina Kartini' },
  { id: 'cs_002', name: 'Dani Saputra' },
  { id: 'cs_003', name: 'Siti Rahayu' }
];

export const MOCK_PACKAGES: PackageOption[] = [
  { id: 'pkg_silver', name: 'Paket Silver Ramadhan 2026' },
  { id: 'pkg_gold', name: 'Paket Gold Syawal 2026' },
  { id: 'pkg_platinum', name: 'Paket Platinum Dzulhijjah 2026' },
  { id: 'pkg_turki', name: 'Umroh Plus Turki' }
];

// --------------- leads ------------------------------------------------

const _leads: Lead[] = [
  {
    id: 'lead_001',
    name: 'Ahmad Fauzi',
    phone: '081234567890',
    email: 'ahmad@example.com',
    source: 'instagram',
    utm_source: 'instagram',
    utm_medium: 'social',
    utm_campaign: 'ramadhan2026',
    interest_package_id: 'pkg_gold',
    interest_package_name: 'Paket Gold Syawal 2026',
    notes: 'Tertarik paket gold, tanya soal cicilan',
    status: 'new',
    assigned_cs_id: 'cs_001',
    assigned_cs_name: 'Rina Kartini',
    created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'lead_002',
    name: 'Siti Nurhaliza',
    phone: '085298765432',
    source: 'whatsapp',
    utm_source: 'wa',
    utm_medium: 'direct',
    interest_package_id: 'pkg_silver',
    interest_package_name: 'Paket Silver Ramadhan 2026',
    status: 'contacted',
    assigned_cs_id: 'cs_002',
    assigned_cs_name: 'Dani Saputra',
    created_at: new Date(Date.now() - 26 * 60 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'lead_003',
    name: 'Budi Santoso',
    phone: '08112233445',
    email: 'budi@gmail.com',
    source: 'facebook',
    utm_source: 'facebook',
    utm_medium: 'ads',
    utm_campaign: 'syawal2026',
    status: 'qualified',
    assigned_cs_id: 'cs_001',
    assigned_cs_name: 'Rina Kartini',
    notes: 'Sudah kirim brosur, tunggu konfirmasi keluarga',
    created_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'lead_004',
    name: 'Dewi Rahmawati',
    phone: '082345678901',
    source: 'organic',
    status: 'converted',
    assigned_cs_id: 'cs_003',
    assigned_cs_name: 'Siti Rahayu',
    interest_package_id: 'pkg_platinum',
    interest_package_name: 'Paket Platinum Dzulhijjah 2026',
    notes: 'Sudah booking — UMR-XYZ12345',
    created_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'lead_005',
    name: 'Hendra Wijaya',
    phone: '087612345678',
    source: 'tiktok',
    utm_source: 'tiktok',
    utm_medium: 'organic',
    status: 'lost',
    assigned_cs_id: 'cs_002',
    assigned_cs_name: 'Dani Saputra',
    notes: 'Tidak ada respons setelah 3x follow up',
    created_at: new Date(Date.now() - 14 * 24 * 60 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'lead_006',
    name: 'Nurul Hidayah',
    phone: '089876543210',
    email: 'nurul@example.com',
    source: 'landing_page',
    utm_source: 'google',
    utm_medium: 'cpc',
    utm_campaign: 'umroh-hemat',
    interest_package_id: 'pkg_silver',
    interest_package_name: 'Paket Silver Ramadhan 2026',
    status: 'new',
    created_at: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 30 * 60 * 1000).toISOString()
  }
];

// in-memory mutable copy for update within session
let leads = structuredClone(_leads);
let nextId = 7;

// --------------- helpers ----------------------------------------------

function applyFilters(all: Lead[], filters: LeadFilters): Lead[] {
  return all.filter((l) => {
    if (filters.status && filters.status !== 'all' && l.status !== filters.status) return false;
    if (filters.assigned_cs_id && l.assigned_cs_id !== filters.assigned_cs_id) return false;
    if (filters.search) {
      const q = filters.search.toLowerCase();
      const matchName = l.name.toLowerCase().includes(q);
      const matchPhone = l.phone.includes(q);
      if (!matchName && !matchPhone) return false;
    }
    return true;
  });
}

// --------------- mock API functions -----------------------------------

export async function listLeadsMock(filters: LeadFilters): Promise<LeadList> {
  await new Promise((r) => setTimeout(r, 200));
  const filtered = applyFilters(leads, filters);
  return { leads: filtered, total: filtered.length };
}

export async function createLeadMock(req: CreateLeadRequest): Promise<{ id: string }> {
  await new Promise((r) => setTimeout(r, 300));
  const id = `lead_${String(nextId++).padStart(3, '0')}`;
  const now = new Date().toISOString();
  const sourceVal = (req.source ?? 'other') as Lead['source'];
  const newLead: Lead = {
    id,
    name: req.name,
    phone: req.phone,
    email: req.email,
    source: sourceVal,
    utm_source: req.utm_source,
    utm_medium: req.utm_medium,
    utm_campaign: req.utm_campaign,
    interest_package_id: req.interest_package_id,
    notes: req.notes,
    status: 'new',
    created_at: now,
    updated_at: now
  };
  leads = [newLead, ...leads];
  return { id };
}

export async function updateLeadMock(id: string, req: UpdateLeadRequest): Promise<Lead> {
  await new Promise((r) => setTimeout(r, 250));
  const idx = leads.findIndex((l) => l.id === id);
  if (idx === -1) throw new Error(`Lead ${id} tidak ditemukan`);
  const cs = req.assigned_cs_id
    ? MOCK_CS_USERS.find((c) => c.id === req.assigned_cs_id)
    : undefined;
  leads[idx] = {
    ...leads[idx],
    status: req.status ?? leads[idx].status,
    notes: req.notes ?? leads[idx].notes,
    assigned_cs_id: req.assigned_cs_id ?? leads[idx].assigned_cs_id,
    assigned_cs_name: cs?.name ?? leads[idx].assigned_cs_name,
    updated_at: new Date().toISOString()
  };
  return leads[idx];
}

export async function listCsUsersMock(): Promise<CSUser[]> {
  await new Promise((r) => setTimeout(r, 100));
  return MOCK_CS_USERS;
}

export async function listPackagesMock(): Promise<PackageOption[]> {
  await new Promise((r) => setTimeout(r, 100));
  return MOCK_PACKAGES;
}
