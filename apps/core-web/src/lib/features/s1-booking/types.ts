export type BookingChannel = 'b2c_self';

export type DraftLead = {
  fullName: string;
  email: string;
  whatsapp: string;
  domicile: string;
};

export type DraftBookingPayload = {
  channel: BookingChannel;
  packageId: string;
  departureId: string;
  roomType: 'quad' | 'triple' | 'double';
  lead: DraftLead;
  jamaahCount: number;
};

export type DraftBookingResult = {
  bookingId: string;
  status: 'draft';
  createdAt: string;
  replayed: boolean;
};

export type DraftBookingError = {
  code: string;
  message: string;
};
