import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
  const { package_id } = params;
  return {
    packageId: package_id,
    stubNote: 'Replace with GET /v1/packages/{id} payload in S1-L-03.'
  };
};
