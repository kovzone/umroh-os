import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, url }) => {
  if (url.pathname === '/console/login') {
    return {
      pathname: url.pathname
    };
  }

  const accessToken = cookies.get('umrohos_access_token');
  if (!accessToken) {
    throw redirect(303, '/console/login');
  }

  return {
    pathname: url.pathname
  };
};

