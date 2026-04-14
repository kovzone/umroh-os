# State Sharing

How state crosses component boundaries. The two allowed mechanisms are **context** (for scoped per-request state) and **classes with `$state` fields** (for reusable reactive state). Writable stores from shared modules are banned.

## Context over shared-module state

For state that is scoped to a subtree — a booking wizard, a modal dialog, a per-page session — use **context**, not a top-level `.ts` module holding `$state`.

Why: module-level `$state` leaks between requests during server-side rendering. Context is scoped to the component tree of the current render.

## Use `createContext` (Svelte 5.40+)

`createContext` returns a type-safe `[get, set]` pair — no stringly-typed keys, full TypeScript inference. Prefer it over `setContext` / `getContext`.

```ts
// bookingContext.ts
import { createContext } from 'svelte';

export const [getBookingCtx, setBookingCtx] = createContext<BookingContext>();
```

```svelte
<!-- BookingWizard.svelte (parent) -->
<script lang="ts">
    import { setBookingCtx } from './bookingContext';

    const ctx = {
        booking: $state(initialBooking),
        advance: () => { /* ... */ },
    };
    setBookingCtx(ctx);
</script>
```

```svelte
<!-- StepTwo.svelte (descendant) -->
<script lang="ts">
    import { getBookingCtx } from './bookingContext';
    const ctx = getBookingCtx();
</script>

<input bind:value={ctx.booking.leaderName} />
```

## Classes with `$state` — not stores

To share reactivity between components without context, export a class instance with `$state` fields. **Do not use writable stores.**

```ts
// fxRate.svelte.ts
class FxRate {
    idrPerSar = $state(4100);
    locked = $state(true);

    rotate(newRate: number) {
        this.idrPerSar = newRate;
    }
}

export const fxRate = new FxRate();
```

Any component importing `fxRate` gets reactive access:

```svelte
<script>
    import { fxRate } from './fxRate.svelte.ts';
</script>

<p>Current rate: {fxRate.idrPerSar}</p>
```

> Module-exported class instances still have the SSR leak caveat. Scope anything request-bound via context. Class-as-shared-state is safe for genuinely global, non-user-bound values (feature flags, i18n locale, etc.) and for state that only exists in the browser.

## See also

- `.claude/skills/svelte-core-bestpractices/SKILL.md` section _"Context"_
