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
  /** Optional; shown on marketing detail when set */
  airline?: string;
  /** Minimum list price per pax (IDR) across all room types; nil when not yet set */
  pricePerPaxIdr?: number;
};

export type PackageInclusion = {
  icon: string;
  title: string;
  description: string;
};

export type PackageItineraryDay = {
  dayLabel: string;
  title: string;
  body: string;
};

export type PackageFaq = {
  question: string;
  answer: string;
};

/** BL-B2C-012: Accommodation spec per city */
export type PackageAccommodation = {
  city: string;
  name: string;
  distance: string;
  roomType: string;
  mealPlan: string;
};

/** BL-B2C-013: Guide/muthawwif profile */
export type PackageGuide = {
  name: string;
  credentials: string;
  experience: string;
  photo: string;
  specialty: string;
};

/** BL-B2C-021: Kitting item included with package */
export type PackageKittingItem = {
  icon: string;
  name: string;
  description: string;
};

export type PackageDetail = {
  id: string;
  kind: string;
  name: string;
  description: string;
  highlights: string[];
  coverPhotoUrl: string;
  startingPriceLabel: string;
  /** Large hero price line, e.g. "Rp 38,5 jt" (optional; falls back to startingPriceLabel) */
  displayPriceShort?: string;
  departures: DepartureSummary[];
  /** Rich marketing fields (mock-first; live API may omit until extended) */
  galleryPhotoUrls?: string[];
  primaryBadge?: string;
  secondaryBadges?: string[];
  priceFinePrint?: string;
  trustPpiu?: string;
  ratingScore?: string;
  ratingReviewsLabel?: string;
  whatsappHref?: string;
  inclusions?: PackageInclusion[];
  importantNotes?: string[];
  itineraryDays?: PackageItineraryDay[];
  faqs?: PackageFaq[];
  /** Short paragraph for Fasilitas tab / section */
  facilitiesBody?: string;
  /** Short bullets or paragraph for S&K tab */
  termsSummary?: string;
  /** BL-B2C-012: Hotel specs per city */
  accommodations?: PackageAccommodation[];
  /** BL-B2C-013: Muthawwif / guide profiles */
  guides?: PackageGuide[];
  /** BL-B2C-021: Kitting items included in package */
  kittingItems?: PackageKittingItem[];
};

export type DeparturePricing = {
  roomType: RoomType;
  amountLabel: string;
  /** Whole IDR per jamaah for this room type (optional; set from API `list_amount` when IDR). */
  listAmountIdr?: number;
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
    itinerary?: {
      id: string;
      public_url: string;
      days: Array<{
        day: number;
        title: string;
        description: string;
        photo_url?: string;
      }>;
    };
    hotels?: Array<{
      id: string;
      name: string;
      city: string;
      star_rating: number;
      walking_distance_m: number;
    }>;
    airline?: {
      id: string;
      code: string;
      name: string;
      operator_kind: string;
    };
    muthawwif?: {
      id: string;
      name: string;
      portrait_url: string;
    };
    add_ons?: Array<{
      id: string;
      name: string;
      list_amount: number;
      list_currency: string;
    }>;
    departures: Array<{
      id: string;
      departure_date: string;
      return_date: string;
      status: 'open' | 'closed';
      remaining_seats: number;
      price_per_pax?: number;
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
