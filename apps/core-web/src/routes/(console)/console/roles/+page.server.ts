import { redirect } from '@sveltejs/kit';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

type LoadEvent = { cookies: Cookies; fetch: RequestEvent['fetch'] };

const MOCK = process.env.MOCK_IAM === 'true' || process.env.VITE_MOCK_IAM === 'true';
const baseUrl =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export type PermAction = 'read' | 'write' | 'delete' | 'export';

export interface Permission {
  resource: string;
  actions: PermAction[];
}

export interface Role {
  id: string;
  name: string;
  description: string;
  permissions: Permission[];
}

export interface PageData {
  roles: Role[];
  allResources: string[];
  error: string | null;
}

const ALL_RESOURCES = [
  'packages',
  'departures',
  'bookings',
  'jamaah',
  'payments',
  'visa_docs',
  'ops',
  'finance',
  'leads',
  'catalog',
  'users',
  'roles',
  'audit_log'
];

function mockData(): PageData {
  return {
    roles: [
      {
        id: 'r1',
        name: 'admin',
        description: 'Akses penuh ke seluruh sistem',
        permissions: ALL_RESOURCES.map((r) => ({
          resource: r,
          actions: ['read', 'write', 'delete', 'export']
        }))
      },
      {
        id: 'r2',
        name: 'finance',
        description: 'Akses ke modul keuangan dan laporan',
        permissions: [
          { resource: 'finance', actions: ['read', 'write', 'export'] },
          { resource: 'payments', actions: ['read', 'write'] },
          { resource: 'bookings', actions: ['read'] }
        ]
      },
      {
        id: 'r3',
        name: 'cs',
        description: 'Customer Service — akses leads dan booking',
        permissions: [
          { resource: 'leads', actions: ['read', 'write'] },
          { resource: 'bookings', actions: ['read', 'write'] },
          { resource: 'jamaah', actions: ['read'] },
          { resource: 'packages', actions: ['read'] }
        ]
      },
      {
        id: 'r4',
        name: 'ops',
        description: 'Operasional — manifest, logistik, visa',
        permissions: [
          { resource: 'ops', actions: ['read', 'write'] },
          { resource: 'visa_docs', actions: ['read', 'write'] },
          { resource: 'jamaah', actions: ['read', 'write'] },
          { resource: 'departures', actions: ['read'] }
        ]
      }
    ],
    allResources: ALL_RESOURCES,
    error: null
  };
}

export const load = async ({ cookies, fetch }: LoadEvent): Promise<PageData> => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  if (MOCK) {
    return mockData();
  }

  try {
    const headers = { Authorization: `Bearer ${accessToken}` };
    const [rolesRes] = await Promise.all([fetch(`${baseUrl}/v1/iam/roles`, { headers })]);

    const roles = rolesRes.ok ? ((await rolesRes.json()) as Role[]) : [];
    return { roles, allResources: ALL_RESOURCES, error: rolesRes.ok ? null : `Gagal memuat roles (${rolesRes.status})` };
  } catch {
    return { ...mockData(), error: 'Tidak dapat terhubung ke gateway.' };
  }
};
