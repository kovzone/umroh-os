// The status page is a CSR-only dashboard — it polls the gateway every 5s.
// SSR would render an empty grid before the first poll completes; opting out
// of SSR keeps the rendered HTML honest and avoids a flash of empty cards.
export const ssr = false;
export const prerender = false;
