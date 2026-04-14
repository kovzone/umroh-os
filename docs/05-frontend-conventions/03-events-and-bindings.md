# Events, Bindings, and Attachments

Svelte 5 turns events into plain attributes and replaces actions with attachments. The old `on:click` / `use:action` syntax is banned.

## Event listeners

Any attribute starting with `on` is an event listener:

```svelte
<button onclick={handleClick}>click me</button>
<input oninput={(e) => value = e.currentTarget.value} />
```

Shorthand and spread both work:

```svelte
<button {onclick}>...</button>
<button {...props}>...</button>
```

## Window / document / body listeners

Use the special elements, not `onMount` + `addEventListener`:

```svelte
<svelte:window onkeydown={handleKey} />
<svelte:document onvisibilitychange={handleVisibility} />
```

Avoid `onMount` or `$effect` for global listeners — let the compiler wire and clean them up.

## Two-way bindings

`bind:value={x}` is the standard form. For validation or transformation, use a **function binding** (Svelte 5.9+), which takes `{get, set}`:

```svelte
<input bind:value={() => value, (v) => (value = v.toLowerCase())} />
```

For readonly bindings (like `clientWidth` / `clientHeight`), pass `null` as the getter:

```svelte
<div bind:clientWidth={null, redraw}>...</div>
```

## Attachments (`{@attach}`)

Attachments replace actions. They are functions that run when an element mounts to the DOM, and re-run when any `$state` they read changes. Available from Svelte 5.29.

```svelte
<script>
    /** @type {import('svelte/attachments').Attachment} */
    function autofocus(element) {
        element.focus();
    }
</script>

<input {@attach autofocus} />
```

Attachment factories — functions that _return_ attachments — enable reactive third-party integrations:

```svelte
<script>
    import tippy from 'tippy.js';

    let content = $state('Hello!');

    function tooltip(content) {
        return (element) => {
            const t = tippy(element, { content });
            return t.destroy;  // cleanup
        };
    }
</script>

<button {@attach tooltip(content)}>Hover me</button>
```

Because the factory runs inside an effect, the attachment is destroyed and recreated whenever `content` changes. If that's too aggressive, pass a getter function and read it inside a child `$effect`:

```js
function foo(getBar) {
    return (node) => {
        veryExpensiveSetup(node);
        $effect(() => {
            update(node, getBar());
        });
    };
}
```

Falsy values disable the attachment — handy for conditionals:

```svelte
<div {@attach enabled && myAttachment}>...</div>
```

## See also

- `.claude/skills/svelte-core-bestpractices/references/@attach.md`
- `.claude/skills/svelte-core-bestpractices/references/bind.md`
