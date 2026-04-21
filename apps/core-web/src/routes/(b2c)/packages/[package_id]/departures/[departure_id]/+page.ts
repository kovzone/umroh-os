import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
  const { package_id, departure_id } = params;
  return {
    packageId: package_id,
    departureId: departure_id,
    stubNote: 'Replace with GET /v1/package-departures/{id} in S1-L-03.'
  };
};
