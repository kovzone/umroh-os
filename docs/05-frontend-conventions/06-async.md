# Async

Svelte 5.36 introduced real `await` expressions inside components. The feature is still experimental — opt in explicitly per project, and understand the tradeoffs before using.

## Enable the flag

`await` inside components requires `experimental.async: true` in `svelte.config.js`. The flag will be removed in Svelte 6.

```js
// svelte.config.js
export default {
    compilerOptions: {
        experimental: {
            async: true,
        },
    },
};
```

## Where you can `await`

Three places:

- Top of a component's `<script>`.
- Inside a `$derived(...)` declaration.
- Directly in markup.

```svelte
<script>
    let a = $state(1);
    let b = $state(2);

    async function add(a, b) {
        await delay(500);
        return a + b;
    }
</script>

<input type="number" bind:value={a} />
<p>{a} + {b} = {await add(a, b)}</p>
```

When `a` changes, the `<p>` text does **not** flip to an inconsistent intermediate state — it holds the prior rendered value until `add(a, b)` resolves.

## Concurrency

Independent `await` expressions in markup run **in parallel**:

```svelte
<p>{await one()}</p><p>{await two()}</p>   <!-- both kick off at once -->
```

Sequential `await`s inside `<script>` or an async function run sequentially (normal JS semantics). Independent `$derived` values are the exception — they update independently after the first render.

If Svelte detects a waterfall pattern, it emits an `await_waterfall` warning — don't ignore it.

## Loading states and boundaries

Use `<svelte:boundary>` with a `pending` snippet for initial-render placeholders. Use `$effect.pending()` to detect subsequent async work (e.g. "validating…" next to a form field).

For coordinating multiple state changes, use `tick` + `settled` from `svelte`:

```js
import { tick, settled } from 'svelte';

updating = true;
await tick();            // flush the `updating = true` update first
color = 'octarine';
answer = 42;
await settled();         // wait for downstream updates
updating = false;
```

## `hydratable` — don't double-fetch on hydration

When you `await` data on the server, the client would re-fetch during hydration by default. `hydratable` serializes the server's result and replays it on the client.

```svelte
<script>
    import { hydratable } from 'svelte';
    import { getUser } from 'my-database-library';

    const user = await hydratable('user', () => getUser());
</script>

<h1>{user.name}</h1>
```

In practice you rarely call `hydratable` directly — the data-fetching library (e.g. SvelteKit remote functions) uses it under the hood.

If you're using CSP, pass a `nonce` or `hash` to `render(...)` — `hydratable` injects an inline `<script>` into `head`.

## Error handling

Errors in `await` expressions bubble to the nearest `<svelte:boundary>` error boundary.

## See also

- `.claude/skills/svelte-core-bestpractices/references/await-expressions.md`
- `.claude/skills/svelte-core-bestpractices/references/hydratable.md`
