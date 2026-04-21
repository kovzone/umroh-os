import type { PageLoad } from './$types';

/** Stub catalog rows — replaced by `GET /v1/packages` in S1-L-03. */
export const load: PageLoad = async () => {
  return {
    packages: [
      {
        id: 'demo-pkg-umrah-12d',
        name: 'Umrah Executive — 12 days',
        blurb: 'Placeholder row until catalog API is wired.'
      },
      {
        id: 'demo-pkg-economy-9d',
        name: 'Umrah Economy — 9 days',
        blurb: 'Placeholder row until catalog API is wired.'
      }
    ]
  };
};
