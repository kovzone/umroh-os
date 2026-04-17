// service-status.svelte.ts — class-with-$state pattern (per docs/05-frontend-conventions/05-state-sharing.md).
// `createSubscriber` owns the 5s poll loop (per .claude/skills/svelte-core-bestpractices/references/svelte-reactivity.md):
// the interval starts when the first reactive consumer reads `services` (or `lastPolledAt`),
// and stops when the last consumer goes away. No manual `$effect` + `setInterval` boilerplate;
// lifecycle is automatic.

import { createSubscriber } from 'svelte/reactivity';
import { gateway } from '$lib/api/gateway/client';

export type ServiceStatusKind = 'pending' | 'ok' | 'fail';

export interface ServiceState {
  name: string;       // "iam-svc"
  shortName: string;  // "iam" — used as the gateway proxy path segment /v1/<shortName>/system/live
  status: ServiceStatusKind;
  error: string | null;
}

// One row per backend the gateway fronts. Keep this list in lock-step with
// services/gateway-svc/api/rest_oapi/openapi.yaml /v1/<svc>/system/live entries.
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

export class ServiceStatus {
  // Backing state — reactive but private. The public getters call #subscribe()
  // first so that any reactive consumer (template, $derived, $effect) activates
  // the poll lifecycle.
  #services = $state<ServiceState[]>([]);
  #lastPolledAt = $state<Date | null>(null);

  #subscribe: () => void;

  constructor(initial: ReadonlyArray<{ name: string; shortName: string }>) {
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
    await Promise.all(this.#services.map((svc, idx) => this.#pollOne(svc, idx)));
  }

  async #pollOne(svc: ServiceState, idx: number): Promise<void> {
    try {
      // openapi-fetch's path argument is type-checked against the generated paths union.
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const { data, error, response } = await (gateway as any).GET(`/v1/${svc.shortName}/system/live`);
      if (error || !response.ok) {
        this.#services[idx] = { ...svc, status: 'fail', error: `HTTP ${response.status}` };
        return;
      }
      const ok = (data as { data?: { ok?: boolean } } | undefined)?.data?.ok === true;
      this.#services[idx] = {
        ...svc,
        status: ok ? 'ok' : 'fail',
        error: ok ? null : 'envelope.data.ok was not true'
      };
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : String(e);
      this.#services[idx] = { ...svc, status: 'fail', error: message };
    }
  }
}
