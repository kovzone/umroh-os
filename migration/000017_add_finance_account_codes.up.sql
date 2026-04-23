-- 000017 — Finance chart of accounts table + seed data.
--
-- Adds finance.account_codes with the standard COA used by UmrohOS finance-svc:
--   1001  Kas/Bank                 (asset, debit-normal)
--   1002  Piutang Jamaah           (asset, debit-normal)
--   2001  Hutang Jamaah (DP)       (liability, credit-normal)
--   4001  Pendapatan Paket         (revenue, credit-normal)
--   5001  HPP / Biaya Paket        (expense, debit-normal)
--
-- finance.journal_lines.account_code references these codes but is NOT
-- constrained by FK so existing journal data is unaffected if this migration
-- is applied after data already exists.
--
-- Schema: finance (already created by migration 000013).

CREATE TABLE finance.account_codes (
    code           TEXT        PRIMARY KEY,
    name           TEXT        NOT NULL,
    type           TEXT        NOT NULL
        CHECK (type IN ('asset', 'liability', 'equity', 'revenue', 'expense')),
    normal_balance TEXT        NOT NULL
        CHECK (normal_balance IN ('debit', 'credit')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO finance.account_codes (code, name, type, normal_balance) VALUES
    ('1001', 'Kas/Bank',            'asset',     'debit'),
    ('1002', 'Piutang Jamaah',      'asset',     'debit'),
    ('2001', 'Hutang Jamaah (DP)',   'liability', 'credit'),
    ('4001', 'Pendapatan Paket',     'revenue',   'credit'),
    ('5001', 'HPP / Biaya Paket',    'expense',   'debit');
