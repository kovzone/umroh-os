// Keep prerender for production builds, but disable it in local dev so
// landing-page edits always reflect immediately while iterating in Docker.
export const prerender = process.env.NODE_ENV === 'production';
