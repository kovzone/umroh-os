export type RoomType = 'quad' | 'triple' | 'double';

export type PackageCard = {
  id: string;
  kind: string;
  name: string;
  blurb: string;
  coverPhotoUrl: string;
  startingPriceLabel: string;
  nextDepartureLabel: string;
  remainingSeats: number;
};

export type DepartureSummary = {
  id: string;
  departureDate: string;
  returnDate: string;
  status: 'open' | 'closed';
  remainingSeats: number;
};

export type PackageDetail = {
  id: string;
  kind: string;
  name: string;
  description: string;
  highlights: string[];
  coverPhotoUrl: string;
  startingPriceLabel: string;
  departures: DepartureSummary[];
};

export type DeparturePricing = {
  roomType: RoomType;
  amountLabel: string;
};

export type DepartureDetail = {
  id: string;
  packageId: string;
  departureDate: string;
  returnDate: string;
  totalSeats: number;
  remainingSeats: number;
  status: 'open' | 'closed';
  pricing: DeparturePricing[];
};
