import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
  return {
    packageId: params.package_id,
    stubNote:
      'Draft create via POST /v1/bookings lands in S1-L-04; form fields stay disabled until booking-svc is wired.'
  };
};
