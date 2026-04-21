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

export type CatalogListResponse = {
  packages: Array<{
    id: string;
    kind: string;
    name: string;
    description: string;
    cover_photo_url: string;
    starting_price: {
      list_amount: number;
      list_currency: string;
      settlement_currency: string;
    };
    next_departure?: {
      id: string;
      departure_date: string;
      return_date: string;
      remaining_seats: number;
    };
  }>;
};

export type CatalogPackageDetailResponse = {
  package: {
    id: string;
    kind: string;
    name: string;
    description: string;
    cover_photo_url: string;
    highlights: string[];
    departures: Array<{
      id: string;
      departure_date: string;
      return_date: string;
      status: 'open' | 'closed';
      remaining_seats: number;
    }>;
  };
};

export type CatalogDepartureDetailResponse = {
  departure: {
    id: string;
    package_id: string;
    departure_date: string;
    return_date: string;
    total_seats: number;
    remaining_seats: number;
    status: 'open' | 'closed';
    pricing: Array<{
      room_type: RoomType;
      list_amount: number;
      list_currency: string;
    }>;
  };
};
