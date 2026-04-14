# Components and Snippets

Components are Svelte files. Snippets are reusable chunks of markup declared _inside_ templates. Slots are gone — everything is a snippet in Svelte 5.

## Snippets

Declare with `{#snippet name(params)}...{/snippet}` and render with `{@render name(args)}`:

```svelte
{#snippet greeting(name)}
    <p>hello {name}!</p>
{/snippet}

{@render greeting('world')}
```

Rules:
- Snippets must be declared **within a template** (not in `<script>`).
- A snippet at the top level of a component (not nested in blocks) can be referenced from `<script>`.
- A top-level snippet that references no component state is also available in `<script module>` and can be `export`ed for use in other components (Svelte 5.5+).
- Any content inside component tags that is **not** a snippet declaration becomes the implicit `children` snippet:

```svelte
<!-- Button.svelte -->
<script>
    let { children } = $props();
</script>
<button>{@render children()}</button>

<!-- usage -->
<Button>click me</Button>
```

## Typing snippets (TypeScript)

Import the `Snippet` type. The type argument is a tuple of the snippet's parameters:

```svelte
<script lang="ts" generics="T">
    import type { Snippet } from 'svelte';

    let { data, row }: { data: T[]; row: Snippet<[T]> } = $props();
</script>
```

## Each blocks — always keyed

Always provide a key that uniquely identifies each item. **Never use the index as the key** — it defeats the optimization.

```svelte
{#each bookings as booking (booking.id)}
    <BookingCard {booking} />
{/each}
```

Avoid destructuring inside `{#each ... as item}` if you need to mutate the item (e.g. `bind:value={item.count}`) — the destructured binding loses the reference back to the array.

Destructuring and rest patterns are fine for read-only rendering:

```svelte
{#each items as { id, name, qty }, i (id)}
    <li>{i + 1}: {name} × {qty}</li>
{/each}
```

## Dynamic components

Components are first-class values. Don't use `<svelte:component>` — that's legacy:

```svelte
<!-- do this -->
<DynamicComponent {...props} />

<!-- not this -->
<svelte:component this={DynamicComponent} {...props} />
```

For recursive self-reference, import the component as `Self` and use `<Self>`:

```svelte
<script>
    import Self from './TreeNode.svelte';
</script>
```

## See also

- `.claude/skills/svelte-core-bestpractices/references/snippet.md`
- `.claude/skills/svelte-core-bestpractices/references/@render.md`
- `.claude/skills/svelte-core-bestpractices/references/each.md`
