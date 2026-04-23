import { redirect } from '@sveltejs/kit';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

type LoadEvent = { cookies: Cookies; fetch: RequestEvent['fetch'] };

const MOCK = process.env.MOCK_IAM === 'true' || process.env.VITE_MOCK_IAM === 'true';
const baseUrl =
  process.env.GATEWAY_URL ?? process.env.VITE_GATEWAY_URL ?? 'http://localhost:4000';

export type UserStatus = 'active' | 'suspended' | 'pending';

export interface AppUser {
  id: string;
  name: string;
  email: string;
  roles: string[];
  status: UserStatus;
  last_login: string | null;
  created_at: string;
}

export interface PageData {
  users: AppUser[];
  error: string | null;
}

function mockData(): PageData {
  return {
    users: [
      {
        id: 'u1',
        name: 'Budi Santoso',
        email: 'budi@umrohos.id',
        roles: ['admin', 'finance'],
        status: 'active',
        last_login: new Date(Date.now() - 3600000).toISOString(),
        created_at: '2024-01-15T08:00:00Z'
      },
      {
        id: 'u2',
        name: 'Siti Rahayu',
        email: 'siti@umrohos.id',
        roles: ['cs'],
        status: 'active',
        last_login: new Date(Date.now() - 86400000).toISOString(),
        created_at: '2024-02-10T09:30:00Z'
      },
      {
        id: 'u3',
        name: 'Ahmad Fauzan',
        email: 'ahmad@umrohos.id',
        roles: ['ops'],
        status: 'suspended',
        last_login: new Date(Date.now() - 7 * 86400000).toISOString(),
        created_at: '2024-03-05T11:00:00Z'
      },
      {
        id: 'u4',
        name: 'Dewi Kusuma',
        email: 'dewi@umrohos.id',
        roles: ['cs', 'ops'],
        status: 'pending',
        last_login: null,
        created_at: '2024-04-20T14:00:00Z'
      },
      {
        id: 'u5',
        name: 'Rizky Pratama',
        email: 'rizky@umrohos.id',
        roles: ['finance'],
        status: 'active',
        last_login: new Date(Date.now() - 2 * 3600000).toISOString(),
        created_at: '2024-01-28T10:15:00Z'
      }
    ],
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
    const res = await fetch(`${baseUrl}/v1/iam/users`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    });

    if (res.ok) {
      const body = (await res.json()) as AppUser[] | { users: AppUser[] };
      const users = Array.isArray(body) ? body : body.users ?? [];
      return { users, error: null };
    }

    return { users: [], error: `Gagal memuat data pengguna (${res.status})` };
  } catch {
    return { users: [], error: 'Tidak dapat terhubung ke gateway.' };
  }
};
