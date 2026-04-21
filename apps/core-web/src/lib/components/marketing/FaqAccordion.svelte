<script lang="ts">
  type FaqItem = {
    question: string;
    answer: string;
  };

  let { items }: { items: FaqItem[] } = $props();
</script>

<div class="faq-list">
  {#each items as faq (faq.question)}
    <details class="faq-item">
      <summary>
        {faq.question}
        <span class="material-symbols-outlined chev">expand_more</span>
      </summary>
      <p>{faq.answer}</p>
    </details>
  {/each}
</div>

<style>
  .faq-list {
    max-width: 48rem;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  .faq-item {
    border-radius: 1rem;
    background: #efeded;
    padding: 0;
    border: none;
  }
  .faq-item summary {
    list-style: none;
    cursor: pointer;
    padding: 1.25rem 1.25rem;
    font-weight: 700;
    font-size: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
  }
  .faq-item summary::-webkit-details-marker {
    display: none;
  }
  .chev {
    transition: transform 0.2s ease;
    color: #6f7a72;
    font-size: 1.25rem;
  }
  .faq-item[open] .chev {
    transform: rotate(180deg);
  }
  .faq-item p {
    margin: 0;
    padding: 0 1.25rem 1.25rem;
    color: #3f4943;
    line-height: 1.6;
    font-size: 0.92rem;
  }
</style>
