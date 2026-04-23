-- S4-E-02 — CRM Lead Tracker schema
-- Creates crm schema and leads table with UTM attribution fields,
-- CS assignment, and status lifecycle.

CREATE SCHEMA IF NOT EXISTS crm;

CREATE TABLE crm.leads (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source TEXT NOT NULL DEFAULT 'direct'
    CHECK (source IN ('organic','whatsapp','instagram','facebook','referral','agent','direct')),
  utm_source TEXT,
  utm_medium TEXT,
  utm_campaign TEXT,
  utm_content TEXT,
  utm_term TEXT,
  name TEXT NOT NULL,
  phone TEXT NOT NULL,
  email TEXT,
  interest_package_id UUID,
  interest_departure_id UUID,
  status TEXT NOT NULL DEFAULT 'new'
    CHECK (status IN ('new','contacted','qualified','converted','lost')),
  assigned_cs_id UUID,
  notes TEXT,
  booking_id TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON crm.leads(status);
CREATE INDEX ON crm.leads(assigned_cs_id);
CREATE INDEX ON crm.leads(created_at DESC);
