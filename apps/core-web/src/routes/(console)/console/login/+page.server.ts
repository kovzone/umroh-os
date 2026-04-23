import { env } from '$env/dynamic/private';
import { dev } from '$app/environment';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type SessionLoginSuccess = {
  data?: {
    access_token?: string;
    refresh_token?: string;
    access_expires_at?: string;
    refresh_expires_at?: string;
  };
};

type SessionLoginError = {
  error?: {
    code?: string;
    message?: string;
  };
};

const iamBaseUrl = env.IAM_API_BASE_URL ?? env.VITE_IAM_API_BASE_URL ?? 'http://localhost:4001';

async function hasValidSession(accessToken: string, fetchFn: typeof fetch): Promise<boolean> {
  try {
    const response = await fetchFn(`${iamBaseUrl}/v1/me`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${accessToken}`
      }
    });
    return response.ok;
  } catch {
    return false;
  }
}

function toValidDate(value: string | undefined): Date | undefined {
  if (!value) {
    return undefined;
  }
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return undefined;
  }
  return parsed;
}

function loginErrorMessage(status: number, fallback?: string): string {
  if (status === 401) {
    return 'Email atau password tidak sesuai.';
  }
  if (status === 403) {
    return 'Akun Anda belum aktif atau sedang ditangguhkan.';
  }
  if (status === 400) {
    return 'Permintaan login tidak valid.';
  }
  return fallback?.trim() || 'Autentikasi gagal. Silakan coba lagi.';
}

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    return {};
  }

  if (await hasValidSession(accessToken, fetch)) {
    throw redirect(303, '/console');
  }

  return {};
};

export const actions: Actions = {
  default: async ({ request, cookies, fetch }) => {
    const form = await request.formData();
    const email = String(form.get('email') ?? '').trim().toLowerCase();
    const password = String(form.get('password') ?? '');
    const rememberMe = form.get('remember-me') === 'on';

    if (!email || !password) {
      return fail(400, {
        error: 'Email dan password wajib diisi.',
        values: { email }
      });
    }

    let response: Response;
    try {
      response = await fetch(`${iamBaseUrl}/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
    } catch {
      return fail(503, {
        error: 'Layanan autentikasi tidak tersedia saat ini.',
        values: { email }
      });
    }

    const body = (await response.json().catch(() => ({}))) as SessionLoginSuccess & SessionLoginError;
    const data = body.data;

    if (!response.ok || !data?.access_token || !data.refresh_token) {
      return fail(response.status || 500, {
        error: loginErrorMessage(response.status, body.error?.message),
        values: { email }
      });
    }

    const refreshExpiresAt = toValidDate(data.refresh_expires_at);
    const accessExpiresAt = toValidDate(data.access_expires_at);
    const secureCookie = !dev;

    cookies.set('umrohos_refresh_token', data.refresh_token, {
      path: '/',
      httpOnly: true,
      sameSite: 'lax',
      secure: secureCookie,
      ...(rememberMe && refreshExpiresAt ? { expires: refreshExpiresAt } : {})
    });

    cookies.set('umrohos_access_token', data.access_token, {
      path: '/',
      httpOnly: true,
      sameSite: 'lax',
      secure: secureCookie,
      ...(accessExpiresAt ? { expires: accessExpiresAt } : {})
    });

    throw redirect(303, '/console');
  }
};

