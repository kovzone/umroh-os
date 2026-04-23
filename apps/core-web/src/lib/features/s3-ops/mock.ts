import type { Departure, FulfillmentTask, FulfillmentTaskList, OpsFilters, StatusUpdate } from './types';

export const MOCK_DEPARTURES: Departure[] = [
  { id: 'dep_001', label: 'Paket Silver — Ramadhan 2026', departure_date: '2026-03-01' },
  { id: 'dep_002', label: 'Paket Gold — Syawal 2026', departure_date: '2026-04-15' },
  { id: 'dep_003', label: 'Paket Platinum — Dzulhijjah 2026', departure_date: '2026-06-20' }
];

const MOCK_TASKS: FulfillmentTask[] = [
  {
    id: 'ft_0001',
    booking_id: 'bkg_abc12345',
    booking_code: 'UMR-ABC12345',
    jamaah_name: 'Ahmad Fauzi',
    package_name: 'Paket Silver Ramadhan 2026',
    departure_date: '2026-03-01',
    payment_status: 'paid_in_full',
    fulfillment_status: 'queued'
  },
  {
    id: 'ft_0002',
    booking_id: 'bkg_def67890',
    booking_code: 'UMR-DEF67890',
    jamaah_name: 'Siti Rahmawati',
    package_name: 'Paket Silver Ramadhan 2026',
    departure_date: '2026-03-01',
    payment_status: 'paid_in_full',
    fulfillment_status: 'processing'
  },
  {
    id: 'ft_0003',
    booking_id: 'bkg_ghi11111',
    booking_code: 'UMR-GHI11111',
    jamaah_name: 'Budi Santoso',
    package_name: 'Paket Gold Syawal 2026',
    departure_date: '2026-04-15',
    payment_status: 'paid_in_full',
    fulfillment_status: 'shipped',
    tracking_number: 'JNE-2026-XXXYYY',
    shipped_at: '2026-04-10T08:30:00Z'
  },
  {
    id: 'ft_0004',
    booking_id: 'bkg_jkl22222',
    booking_code: 'UMR-JKL22222',
    jamaah_name: 'Dewi Kartika',
    package_name: 'Paket Gold Syawal 2026',
    departure_date: '2026-04-15',
    payment_status: 'paid_in_full',
    fulfillment_status: 'delivered',
    tracking_number: 'JNE-2026-AAABBB',
    shipped_at: '2026-04-08T10:00:00Z'
  },
  {
    id: 'ft_0005',
    booking_id: 'bkg_mno33333',
    booking_code: 'UMR-MNO33333',
    jamaah_name: 'Hendra Wijaya',
    package_name: 'Paket Platinum Dzulhijjah 2026',
    departure_date: '2026-06-20',
    payment_status: 'paid_in_full',
    fulfillment_status: 'queued'
  },
  // bkg_pqr44444 (Nurul Hidayah) is partially_paid — no fulfillment task should
  // exist per §S3-J-02 gate rule ("Kit dispatch only on paid_in_full").
  // Removed from the mock task list; partially_paid bookings never appear in the
  // fulfillment queue UI.
  {
    id: 'ft_0007',
    booking_id: 'bkg_stu55555',
    booking_code: 'UMR-STU55555',
    jamaah_name: 'Agus Pratama',
    package_name: 'Paket Silver Ramadhan 2026',
    departure_date: '2026-03-01',
    payment_status: 'paid_in_full',
    fulfillment_status: 'cancelled'
  }
];

// in-memory copy for updates within the session
let tasks = structuredClone(MOCK_TASKS);

function applyFilters(all: FulfillmentTask[], filters: OpsFilters): FulfillmentTask[] {
  return all.filter((t) => {
    if (filters.status && filters.status !== 'all' && t.fulfillment_status !== filters.status) {
      return false;
    }
    if (filters.departure_id) {
      const dep = MOCK_DEPARTURES.find((d) => d.id === filters.departure_id);
      if (dep && t.departure_date !== dep.departure_date) return false;
    }
    if (filters.search) {
      const q = filters.search.toLowerCase();
      const matchCode = t.booking_code.toLowerCase().includes(q);
      const matchName = t.jamaah_name.toLowerCase().includes(q);
      if (!matchCode && !matchName) return false;
    }
    return true;
  });
}

export async function listFulfillmentTasksMock(filters: OpsFilters): Promise<FulfillmentTaskList> {
  await new Promise((r) => setTimeout(r, 200));
  const filtered = applyFilters(tasks, filters);
  return { tasks: filtered, total: filtered.length };
}

export async function updateFulfillmentStatusMock(
  taskId: string,
  update: StatusUpdate
): Promise<FulfillmentTask> {
  await new Promise((r) => setTimeout(r, 250));
  const idx = tasks.findIndex((t) => t.id === taskId);
  if (idx === -1) throw new Error(`Task ${taskId} not found`);
  tasks[idx] = {
    ...tasks[idx],
    fulfillment_status: update.status,
    tracking_number: update.tracking_number ?? tasks[idx].tracking_number,
    shipped_at:
      update.status === 'shipped' && !tasks[idx].shipped_at
        ? new Date().toISOString()
        : tasks[idx].shipped_at
  };
  return tasks[idx];
}
