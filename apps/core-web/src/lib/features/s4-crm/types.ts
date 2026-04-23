// -----------------------------------------------------------------------
// S4-CRM — Shared Types
// -----------------------------------------------------------------------

export type LeadStatus = 'new' | 'contacted' | 'qualified' | 'converted' | 'lost';

// Backend-contracted values per S4-J-01 contract (crm.leads.source CHECK constraint):
//   organic | whatsapp | instagram | facebook | referral | agent | direct
// Frontend-only display values (tiktok, landing_page, other) are used in mock data
// and UI labels; the backend coerces unknown values to 'direct' on insert.
export type LeadSource =
  | 'whatsapp'
  | 'instagram'
  | 'facebook'
  | 'tiktok'       // display-only: not in backend enum, coerced to 'direct'
  | 'organic'
  | 'agent'
  | 'referral'
  | 'direct'
  | 'landing_page' // display-only: not in backend enum, coerced to 'direct'
  | 'other';       // display-only: not in backend enum, coerced to 'direct'

export interface Lead {
  id: string;
  name: string;
  phone: string;
  email?: string;
  source: LeadSource;
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
  interest_package_id?: string;
  interest_package_name?: string;
  notes?: string;
  status: LeadStatus;
  assigned_cs_id?: string;
  assigned_cs_name?: string;
  created_at: string; // ISO 8601
  updated_at: string;
}

export interface LeadList {
  leads: Lead[];
  total: number;
}

export interface LeadFilters {
  status?: LeadStatus | 'all';
  assigned_cs_id?: string;
  search?: string;
}

export interface CSUser {
  id: string;
  name: string;
}

export interface PackageOption {
  id: string;
  name: string;
}

export interface UpdateLeadRequest {
  status?: LeadStatus; // optional — PATCH semantics: omit to keep current status
  notes?: string;
  assigned_cs_id?: string;
}

export interface CreateLeadRequest {
  name: string;
  phone: string;
  email?: string;
  source?: string;
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
  interest_package_id?: string;
  notes?: string;
}
