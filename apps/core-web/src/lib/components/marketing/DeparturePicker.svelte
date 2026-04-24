<script lang="ts">
  import type { DepartureSummary } from '$lib/features/s1-catalog/types';
  import { airlineOrDefault, departureSeatsLabel, formatDepartureRangeId } from '$lib/utils/format-departure';

  let {
    departures,
    selectedId = $bindable<string | null>(null),
    groupName = 'departure'
  }: {
    departures: DepartureSummary[];
    selectedId: string | null;
    groupName?: string;
  } = $props();
</script>

<div class="departure-list" role="radiogroup" aria-label="Pilih keberangkatan">
  {#each departures as dep (dep.id)}
    {@const seats = departureSeatsLabel(dep)}
    <label
      class="departure-row"
      class:urgent={seats.urgent}
      class:is-selected={selectedId === dep.id}
      class:is-closed={dep.status === 'closed'}
    >
      <input
        type="radio"
        name={groupName}
        value={dep.id}
        class="dep-radio"
        checked={selectedId === dep.id}
        disabled={dep.status === 'closed'}
        onchange={() => {
          selectedId = dep.id;
        }}
      />
      <div class="dep-fields" class:has-price={!!dep.pricePerPaxIdr}>
        <div class="dep-field">
          <span class="dep-label">Tanggal keberangkatan</span>
          <span class="dep-value">{formatDepartureRangeId(dep)}</span>
        </div>
        <div class="dep-field">
          <span class="dep-label">Maskapai</span>
          <span class="dep-value dep-airline">
            <span class="material-symbols-outlined gold">flight_takeoff</span>
            {airlineOrDefault(dep)}
          </span>
        </div>
        <div class="dep-field">
          <span class="dep-label">Status kuota</span>
          <span class="dep-seats" class:err={seats.urgent}>{seats.text}</span>
        </div>
        {#if dep.pricePerPaxIdr}
          <div class="dep-field dep-field-price">
            <span class="dep-label">Mulai dari</span>
            <span class="dep-value dep-price">
              Rp {new Intl.NumberFormat('id-ID').format(dep.pricePerPaxIdr)} / pax
            </span>
          </div>
        {/if}
      </div>
    </label>
  {/each}
</div>

<style>
  .departure-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .departure-row {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
    padding: 1.25rem;
    border-radius: 1rem;
    background: #fff;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
    cursor: pointer;
    border: 2px solid transparent;
  }
  .dep-radio {
    margin-top: 0.2rem;
    width: 1.15rem;
    height: 1.15rem;
    accent-color: #004d34;
    flex-shrink: 0;
  }
  .dep-fields {
    flex: 1;
    display: grid;
    gap: 1rem;
    grid-template-columns: 1fr;
  }
  @media (min-width: 768px) {
    .dep-fields {
      grid-template-columns: 1fr 1fr 1fr;
    }
    .dep-fields.has-price {
      grid-template-columns: 1fr 1fr 1fr 1fr;
    }
  }
  .departure-row:hover {
    border-color: rgba(0, 103, 71, 0.25);
  }
  .departure-row.is-selected {
    border-color: #006747;
    background: rgba(0, 103, 71, 0.04);
  }
  .departure-row.urgent {
    border-left: 4px solid #ba1a1a;
  }
  .departure-row.is-closed {
    opacity: 0.55;
    cursor: not-allowed;
  }
  .dep-field {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }
  .dep-label {
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #6f7a72;
  }
  .dep-value {
    font-weight: 700;
    color: #004d34;
    font-size: 0.95rem;
  }
  .dep-airline {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    color: #1b1c1c;
  }
  .material-symbols-outlined.gold {
    color: #775a19;
  }
  .dep-seats {
    font-weight: 700;
    color: #775a19;
    font-size: 0.95rem;
  }
  .dep-seats.err {
    color: #ba1a1a;
    font-style: italic;
  }
  .dep-price {
    color: #006747;
    font-size: 1rem;
  }
  .dep-field-price {
    border-left: 2px solid rgba(0, 103, 71, 0.25);
    padding-left: 0.75rem;
  }
  @media (max-width: 767px) {
    .dep-field-price {
      border-left: none;
      padding-left: 0;
      border-top: 2px solid rgba(0, 103, 71, 0.15);
      padding-top: 0.5rem;
    }
  }
</style>
