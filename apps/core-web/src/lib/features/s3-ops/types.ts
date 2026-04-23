export type FulfillmentStatus = 'queued' | 'processing' | 'shipped' | 'delivered' | 'cancelled';

export interface FulfillmentTask {
  id: string;
  booking_id: string;
  booking_code: string; // UMR-XXXXXXXX
  jamaah_name: string;
  package_name: string;
  departure_date: string;
  payment_status: 'partially_paid' | 'paid_in_full';
  fulfillment_status: FulfillmentStatus;
  tracking_number?: string;
  shipped_at?: string;
}

export interface FulfillmentTaskList {
  tasks: FulfillmentTask[];
  total: number;
}

export interface OpsFilters {
  status?: FulfillmentStatus | 'all';
  departure_id?: string;
  search?: string;
}

export interface StatusUpdate {
  status: FulfillmentStatus;
  tracking_number?: string;
}

export interface Departure {
  id: string;
  label: string;
  departure_date: string;
}
