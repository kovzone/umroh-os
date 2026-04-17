import { svelteTesting } from '@testing-library/svelte/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

// Vitest config is split from vite.config.ts because `vitest/config` pulls in
// Vitest's nested Vite copy, whose Plugin types don't structurally match the
// top-level Vite's Plugin types. Keeping the two configs separate lets each
// tool see its own consistent type universe.
export default defineConfig({
  plugins: [sveltekit(), svelteTesting()],
  test: {
    environment: 'jsdom',
    include: ['src/**/*.{test,spec}.{js,ts}'],
    globals: true,
    setupFiles: ['./vitest.setup.ts']
  }
});
