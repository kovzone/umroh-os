// Vitest unit test — proves the Vitest setup works and locks in two invariants
// of ServiceStatus:
//   1. After construction, every service is in the `pending` state with no error.
//   2. The exported BACKEND_SERVICES list contains exactly 10 entries (every
//      backend the gateway fronts) and uses kebab-case `<svc>-svc` names.

import { describe, it, expect } from 'vitest';
import { BACKEND_SERVICES, ServiceStatus } from './service-status.svelte';

describe('BACKEND_SERVICES registry', () => {
  it('lists exactly 10 backends', () => {
    expect(BACKEND_SERVICES).toHaveLength(10);
  });

  it('uses kebab-case <svc>-svc names paired with shortName', () => {
    for (const s of BACKEND_SERVICES) {
      expect(s.name).toMatch(/^[a-z]+-svc$/);
      expect(s.shortName).toMatch(/^[a-z]+$/);
      expect(s.name).toBe(`${s.shortName}-svc`);
    }
  });

  it('shortNames are unique', () => {
    const set = new Set(BACKEND_SERVICES.map((s) => s.shortName));
    expect(set.size).toBe(BACKEND_SERVICES.length);
  });
});

describe('ServiceStatus initial state', () => {
  it('starts every service in pending state with no error', () => {
    const s = new ServiceStatus(BACKEND_SERVICES);
    expect(s.services).toHaveLength(BACKEND_SERVICES.length);
    for (const svc of s.services) {
      expect(svc.status).toBe('pending');
      expect(svc.error).toBeNull();
    }
    expect(s.lastPolledAt).toBeNull();
  });
});
