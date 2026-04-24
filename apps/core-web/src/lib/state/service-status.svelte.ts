// service-status.svelte.ts — class-with-$state pattern (per docs/05-frontend-conventions/05-state-sharing.md).
// `createSubscriber` owns the 5s poll loop (per .claude/skills/svelte-core-bestpractices/references/svelte-reactivity.md):
// the interval starts when the first reactive consumer reads `services` (or `lastPolledAt`),
// and stops when the last consumer goes away. No manual `$effect` + `setInterval` boilerplate;
// lifecycle is automatic.
//
// Polling backs onto a single aggregate gateway endpoint:
//   GET /v1/system/backends  →  { data: { backends: [{name, status, error?}] } }
// The gateway fans `grpc.health.v1.Health.Check` out to every registered backend
// (BL-IAM-019 / S1-E-14). One request per tick regardless of how many backends
// the platform has, and one uniform wire protocol across every service.

import { createSubscriber } from 'svelte/reactivity';
import { gateway } from '$lib/api/gateway/client';

export type ServiceStatusKind = 'pending' | 'ok' | 'fail';

export interface ServiceState {
  name: string;       // "iam-svc"
  shortName: string;  // "iam" — UI-side short label; derived at seed time.
  status: ServiceStatusKind;
  error: string | null;
}

// Seed list — order here fixes the grid layout. The server's response is
// sorted alphabetically by name; this list is mapped in server-order for the
// UI render, but a server-side entry without a seed entry still appears
// (future backends just need a seed row added).
export const BACKEND_SERVICES: ReadonlyArray<{ name: string; shortName: string }> = [
  { name: 'iam-svc',       shortName: 'iam' },
  { name: 'catalog-svc',   shortName: 'catalog' },
  { name: 'booking-svc',   shortName: 'booking' },
  { name: 'jamaah-svc',    shortName: 'jamaah' },
  { name: 'payment-svc',   shortName: 'payment' },
  { name: 'visa-svc',      shortName: 'visa' },
  { name: 'ops-svc',       shortName: 'ops' },
  { name: 'logistics-svc', shortName: 'logistics' },
  { name: 'finance-svc',   shortName: 'finance' },
  { name: 'crm-svc',       shortName: 'crm' }
];

const POLL_INTERVAL_MS = 5_000;

// Wire shape for /v1/system/backends. Mirror of the openapi component
// BackendStatus schema; kept local to avoid coupling the state class to the
// generated api types file (which is regenerated on every oapi change).
interface BackendStatusWire {
  name: string;
  status: 'SERVING' | 'NOT_SERVING' | 'UNKNOWN';
  error?: string;
}

interface SystemBackendsWire {
  data: { backends: BackendStatusWire[] };
}

export class ServiceStatus {
  // Backing state — reactive but private. The public getters call #subscribe()
  // first so that any reactive consumer (template, $derived, $effect) activates
  // the poll lifecycle.
  #services = $state<ServiceState[]>([]);
  #lastPolledAt = $state<Date | null>(null);

  #subscribe: () => void;
  #shortNameByName: Map<string, string>;

  constructor(initial: ReadonlyArray<{ name: string; shortName: string }>) {
    this.#shortNameByName = new Map(initial.map((s) => [s.name, s.shortName]));
    this.#services = initial.map((s) => ({
      name: s.name,
      shortName: s.shortName,
      status: 'pending',
      error: null
    }));

    // The poller only runs while something reactive is reading the getters
    // (i.e. while the page is mounted and `services` is in the `{#each}` block).
    // When the last reader goes away, createSubscriber calls our cleanup fn.
    this.#subscribe = createSubscriber((update) => {
      const tick = () => {
        this.#pollAll().finally(() => {
          this.#lastPolledAt = new Date();
          update();
        });
      };
      tick(); // first poll immediately
      const id = setInterval(tick, POLL_INTERVAL_MS);
      return () => clearInterval(id);
    });
  }

  /** Reactive list, rendered directly in the grid. Reading it activates the poll loop. */
  get services(): ServiceState[] {
    this.#subscribe();
    return this.#services;
  }

  /** Wall-clock of the most recent poll completion (any outcome). `null` before the first poll. */
  get lastPolledAt(): Date | null {
    this.#subscribe();
    return this.#lastPolledAt;
  }

  async #pollAll(): Promise<void> {
    try {
      // openapi-fetch's path argument is type-checked against the generated paths union.
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const { data, error, response } = await (gateway as any).GET('/v1/system/backends');
      if (error || !response.ok) {
        this.#markAllFail(`HTTP ${response.status}`);
        return;
      }
      const body = data as SystemBackendsWire | undefined;
      const backends = body?.data?.backends ?? [];
      this.#applyBackends(backends);
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : String(e);
      this.#markAllFail(message);
    }
  }

  #applyBackends(wire: BackendStatusWire[]): void {
    const byName = new Map(wire.map((b) => [b.name, b]));
    this.#services = this.#services.map((svc) => {
      const w = byName.get(svc.name);
      if (!w) {
        return { ...svc, status: 'fail', error: 'backend not reported by gateway' };
      }
      const ok = w.status === 'SERVING';
      return {
        ...svc,
        status: ok ? 'ok' : 'fail',
        error: ok ? null : w.error || w.status
      };
    });
  }

  #markAllFail(message: string): void {
    this.#services = this.#services.map((svc) => ({
      ...svc,
      status: 'fail',
      error: message
    }));
  }
}
