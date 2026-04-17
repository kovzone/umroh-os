// Landing page is static content — prerender at build time so it ships as
// fast, cacheable HTML (no SSR roundtrip, no hydration of dynamic state).
// When future landing content becomes dynamic (e.g. status badge pulled live),
// revisit this.
export const prerender = true;
