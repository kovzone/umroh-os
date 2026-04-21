import type { DepartureSummary } from '$lib/features/s1-catalog/types';

const rangeOpts: Intl.DateTimeFormatOptions = {
  day: 'numeric',
  month: 'long',
  year: 'numeric'
};

const dayOpts: Intl.DateTimeFormatOptions = {
  day: 'numeric',
  month: 'long',
  year: 'numeric'
};

export function formatDepartureRangeId(dep: DepartureSummary): string {
  const a = new Date(dep.departureDate);
  const b = new Date(dep.returnDate);
  return `${a.toLocaleDateString('id-ID', rangeOpts)} — ${b.toLocaleDateString('id-ID', rangeOpts)}`;
}

export function formatDepartureDayId(dep: DepartureSummary): string {
  return new Date(dep.departureDate).toLocaleDateString('id-ID', dayOpts);
}

export function departureSeatsLabel(dep: DepartureSummary): { text: string; urgent: boolean } {
  if (dep.status === 'closed') {
    return { text: 'Kuota penuh', urgent: false };
  }
  if (dep.remainingSeats <= 0) {
    return { text: 'Kuota habis', urgent: false };
  }
  if (dep.remainingSeats <= 7) {
    return { text: `Sisa ${dep.remainingSeats} kursi`, urgent: true };
  }
  if (dep.remainingSeats < 20) {
    return { text: `${dep.remainingSeats} kursi tersedia`, urgent: false };
  }
  return { text: 'Tersedia banyak', urgent: false };
}

export function airlineOrDefault(dep: DepartureSummary): string {
  return dep.airline ?? 'Maskapai sesuai kontrak';
}
