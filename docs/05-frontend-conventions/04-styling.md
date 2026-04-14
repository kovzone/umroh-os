# Styling

Component styles are **scoped by default**. For parent→child styling, prefer CSS custom properties over `:global`. For class manipulation, prefer clsx-style arrays over the `class:` directive.

## JS → CSS via `style:`

To pass a JavaScript value into CSS, set a CSS custom property with the `style:` directive:

```svelte
<div style:--columns={columns}>...</div>

<style>
    div {
        grid-template-columns: repeat(var(--columns), 1fr);
    }
</style>
```

## Parent controls child — custom props first

The preferred way for a parent to influence child styling is to pass CSS custom properties:

```svelte
<!-- Parent.svelte -->
<Child --color="red" />
```

```svelte
<!-- Child.svelte -->
<h1>Hello</h1>

<style>
    h1 {
        color: var(--color, black);  /* sensible fallback */
    }
</style>
```

## Parent controls child — `:global` escape hatch

Only when the child is from a library you can't modify (or the styling surface is too large to expose via custom props), use `:global`:

```svelte
<div>
    <ThirdPartyChart />
</div>

<style>
    div :global {
        .chart-tooltip {
            background: var(--brand-color);
        }
    }
</style>
```

## Class attributes — clsx-style, not `class:`

Use arrays/objects directly in the `class` attribute. Do **not** use the `class:` directive.

```svelte
<!-- do this -->
<button class={['btn', size, { primary: isPrimary, disabled: !enabled }]}>
    Save
</button>

<!-- not this -->
<button
    class="btn {size}"
    class:primary={isPrimary}
    class:disabled={!enabled}
>
    Save
</button>
```

## See also

- `.claude/skills/svelte-core-bestpractices/SKILL.md` sections _"Using JavaScript variables in CSS"_, _"Styling child components"_
