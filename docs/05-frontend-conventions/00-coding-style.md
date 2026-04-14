# Coding Style

Svelte code in UmrohOS targets **Svelte 5 in runes mode** exclusively. No legacy Svelte 3/4 features are allowed in new code. The authoritative guide loaded into every session is `.claude/skills/svelte-core-bestpractices/SKILL.md`; these convention files are the human-readable mirror.

## Runes mode only

Every `.svelte` file uses runes (`$state`, `$derived`, `$effect`, `$props`). No implicit reactivity (`let count = 0; count += 1`), no `$:` labels, no `export let`.

## Banned legacy features

Never use these — each has a modern replacement:

| Legacy | Modern |
|---|---|
| Implicit reactivity (`let count = 0`) | `let count = $state(0)` |
| `$:` assignments and statements | `$derived(...)` or `$effect(...)` (prefer `$derived`) |
| `export let` / `$$props` / `$$restProps` | `$props()` |
| `on:click={...}` | `onclick={...}` |
| `<slot>` / `$$slots` / `<svelte:fragment>` | `{#snippet ...}` + `{@render ...}` |
| `<svelte:component this={X}>` | `<X>` directly (components are first-class values) |
| `<svelte:self>` | `import Self from './ThisComponent.svelte'`; `<Self>` |
| Stores (shared-module `writable`/`readable`) | Classes with `$state` fields |
| `use:action` | `{@attach ...}` |
| `class:` directive | clsx-style arrays/objects in `class` attribute |

## Naming _(Inferred — not specified by the skill)_

- **Component files:** `PascalCase.svelte` (e.g. `BookingCard.svelte`, `JamaahList.svelte`).
- **Non-component modules:** `kebab-case.ts` (e.g. `fx-rate.ts`).
- **Route files (if SvelteKit):** framework-dictated (`+page.svelte`, `+layout.svelte`, `+page.server.ts`).
- **Component props:** `camelCase`.
- **CSS custom props:** `--kebab-case` (e.g. `--columns`, `--brand-color`).

See `docs/07-open-questions/` for the frontend naming/layout question awaiting stakeholder input.

## Class attributes

Use clsx-style arrays/objects directly in the `class` attribute rather than the `class:` directive:

```svelte
<!-- do this -->
<button class={['btn', { active: isActive, disabled: !enabled }]}>...</button>

<!-- not this -->
<button class="btn" class:active={isActive} class:disabled={!enabled}>...</button>
```

## See also

- `.claude/skills/svelte-core-bestpractices/SKILL.md` — the canonical skill (loaded automatically in Svelte sessions)
- `01-reactivity.md` — runes reference
- `02-components-and-snippets.md` — snippets vs. slots, keyed `{#each}`
