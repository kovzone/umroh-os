/** Parse labels like "Rp 38.500.000" or "Rp38.500.000" to whole rupiah. */
export function parseIdrAmountLabel(label: string | undefined | null): number | null {
  if (!label) {
    return null;
  }
  const digits = label.replace(/\D/g, '');
  if (!digits) {
    return null;
  }
  const n = Number(digits);
  return Number.isFinite(n) ? n : null;
}

/** Whole rupiah → "Rp 38.500.000" */
export function formatIdrAmountLabel(amount: number): string {
  const formatted = new Intl.NumberFormat('id-ID').format(Math.round(amount));
  return `Rp ${formatted}`;
}
