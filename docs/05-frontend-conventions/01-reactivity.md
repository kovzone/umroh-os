# Reactivity

Svelte 5 exposes reactivity through a small set of runes. Use the right tool for the job — most reactive code should be `$state` + `$derived`, with `$effect` reserved for genuine side effects.

## `$state`

Use `$state` only for variables that drive an `$effect`, `$derived`, or template expression. Everything else is a plain variable.

```js
let count = $state(0);      // reactive
let label = 'submit';       // not reactive — fine
```

Objects and arrays (`$state({...})` / `$state([...])`) are made **deeply reactive** via proxies — mutation triggers updates. For large objects that are only ever _reassigned_ (e.g. API responses), use `$state.raw` to skip the proxy overhead:

```js
let user = $state.raw(await fetchUser());  // reassigned, not mutated
```

## `$derived`

Compute from state with `$derived`, never `$effect`:

```js
// do this
let square = $derived(num * num);

// don't do this
let square;
$effect(() => { square = num * num });
```

`$derived` takes an _expression_. For complex logic needing a function body, use `$derived.by(() => ...)`. Derived values are writable (assignable) but re-evaluate when their dependencies change.

If a derived returns an object/array, it is _not_ made deeply reactive — use `$state` inside `$derived.by` if you need that.

## `$effect` — escape hatch

Effects are a last resort. **Avoid updating state inside effects.** Reach for an alternative first:

| Need | Use instead |
|---|---|
| Sync with an external library (D3, Chart.js, tippy) | `{@attach ...}` — see `03-events-and-bindings.md` |
| Respond to user interaction | Event handler (`onclick={...}`) or function binding |
| Log values for debugging | `$inspect(...)` / `$inspect.trace()` |
| Observe something external to Svelte | `createSubscriber` from `svelte/reactivity` |

Never wrap effect contents in `if (browser) {...}` — effects don't run on the server.

## `$props`

Treat props as though they will change. Derive any computed value from them:

```js
let { type } = $props();

// do this
let color = $derived(type === 'danger' ? 'red' : 'green');

// don't do this — won't update when `type` changes
let color = type === 'danger' ? 'red' : 'green';
```

## `$inspect.trace`

When reactivity misbehaves (something re-runs too often, or fails to update), add `$inspect.trace(label)` as the **first line** of the `$effect`, `$derived.by`, or any function they call. The console will print which piece of state triggered the re-run. Development-only; becomes a noop in production builds.

```svelte
$effect(() => {
    $inspect.trace('booking-sync');
    syncToServer();
});
```

## See also

- `.claude/skills/svelte-core-bestpractices/references/$inspect.md`
- `.claude/skills/svelte-core-bestpractices/references/svelte-reactivity.md` — `createSubscriber`
