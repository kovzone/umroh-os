-- S1-E-02 / BL-CAT-001 — dev-only catalog fixtures.
--
-- Seeds a minimal but complete shape that exercises every public read
-- path: one active package (with full master refs + two upcoming
-- departures + pricing + add-ons + itinerary), one draft package
-- (must be hidden from public list and return 404 on detail), and one
-- archived package (same hiding rule). All inserts are idempotent —
-- re-running the migration on a non-empty schema is a no-op.
--
-- This migration is deliberately separated from the DDL migration
-- (000008) so a future deployment to staging / prod can run 000008
-- without the dev fixtures coming along. In MVP the distinction is
-- manual (skip 000009 in the prod migrate target); a later ADR will
-- formalise env-scoped migrations.
--
-- IDs are synthetic ULID-like strings: 26 chars of uppercase base32 +
-- the contracted prefix. Real ULIDs are generated app-side; these seed
-- values use `01JCDE...` patterns that sort reasonably under lexicographic
-- cursor pagination.

-- ---------------------------------------------------------------------
-- Master data (referenced by the active package)
-- ---------------------------------------------------------------------

INSERT INTO catalog.hotels (id, name, city, star_rating, walking_distance_m)
VALUES
    ('htl_01JCDH00000000000000000001', 'Hotel Dar Al Tawhid Intercontinental', 'mecca',  5, 50),
    ('htl_01JCDH00000000000000000002', 'Hotel Pullman Zamzam Madinah',         'medina', 5, 80)
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.airlines (id, code, name, operator_kind)
VALUES
    ('arl_01JCDI00000000000000000001', 'GA', 'Garuda Indonesia', 'airline')
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.muthawwif (id, name, portrait_url)
VALUES
    ('mtw_01JCDJ00000000000000000001', 'Ustadz Ahmad Al-Azhar', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=320&h=320&fit=crop&auto=format')
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.itinerary_templates (id, name, days, public_url)
VALUES (
    'itn_01JCDG00000000000000000001',
    'Umrah Reguler 12 Hari — Ramadan itinerary',
    '[
        {"day": 1, "title": "Keberangkatan dari Jakarta", "description": "Berkumpul di Terminal 3 CGK pukul 18.00 WIB. Manasik ringkas, pemeriksaan dokumen, dan keberangkatan ke Jeddah dengan penerbangan langsung Garuda Indonesia.", "photo_url": "https://images.unsplash.com/photo-1436491865332-7a61a109cc05?w=800&h=450&fit=crop&auto=format"},
        {"day": 2, "title": "Tiba di Madinah", "description": "Tiba di Bandara Prince Mohammad bin Abdulaziz pagi hari. Check-in hotel bintang 5 di sekitar Masjid Nabawi. Waktu istirahat dan shalat Asar berjamaah di Masjid Nabawi.", "photo_url": "https://images.unsplash.com/photo-1565552645632-d725f8bfc19a?w=800&h=450&fit=crop&auto=format"},
        {"day": 3, "title": "Ziarah Madinah — Raudhah & Makam Baqi", "description": "Shalat Subuh di Masjid Nabawi. Kunjungan ke Raudhah (pagi). Ziarah ke Makam Baqi, Masjid Quba, Masjid Qiblatayn, dan Jabal Uhud.", "photo_url": "https://images.unsplash.com/photo-1580418827493-f2b22c0a76cb?w=800&h=450&fit=crop&auto=format"},
        {"day": 4, "title": "Ibadah Madinah & Persiapan ke Makkah", "description": "Pagi hingga siang ibadah mandiri di Masjid Nabawi. Sore hari persiapan keberangkatan. Malam berangkat ke Makkah via bus dengan persinggahan di Bir Ali untuk miqat.", "photo_url": "https://images.unsplash.com/photo-1565552645632-d725f8bfc19a?w=800&h=450&fit=crop&auto=format"},
        {"day": 5, "title": "Miqat & Pelaksanaan Umrah", "description": "Berihram di Bir Ali. Tiba di Makkah dan langsung menuju Masjidil Haram untuk pelaksanaan ibadah umrah: tawaf 7 putaran, sa''i Shafa-Marwah, dan tahallul.", "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"},
        {"day": 6, "title": "Ibadah di Masjidil Haram", "description": "Hari penuh ibadah di Masjidil Haram. Tawaf sunnah, membaca Al-Qur''an di depan Ka''bah, i''tikaf, dan doa. Waktu shalat berjamaah lima waktu bersama jamaah internasional.", "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"},
        {"day": 7, "title": "Ziarah Makkah", "description": "Ziarah ke tempat-tempat bersejarah di Makkah: Jabal Tsur, Jabal Nur (Gua Hira), Masjid Ji''ranah, dan Tan''im. Sore kembali ke hotel untuk istirahat.", "photo_url": "https://images.unsplash.com/photo-1519817650390-64a93db51149?w=800&h=450&fit=crop&auto=format"},
        {"day": 8, "title": "Ibadah di Masjidil Haram & Belanja", "description": "Pagi ibadah di Masjidil Haram. Siang waktu bebas untuk berbelanja oleh-oleh di area Abraj Al-Bait dan Zamzam Tower. Malam tawaf sunnah.", "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"},
        {"day": 9, "title": "Ibadah Makkah — I''tikaf", "description": "Hari penuh untuk i''tikaf, membaca Al-Qur''an, dan shalat nafilah. Bimbingan doa khusus oleh Ustadz Ahmad Al-Azhar pada malam hari.", "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"},
        {"day": 10, "title": "Ibadah Makkah & Persiapan Pulang", "description": "Pagi ibadah di Masjidil Haram, tawaf wada'' bagi yang berencana pulang. Siang packing dan persiapan. Malam farewell dinner bersama rombongan.", "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"},
        {"day": 11, "title": "Tawaf Wada'' & Perjalanan ke Jeddah", "description": "Tawaf Wada'' (tawaf perpisahan) di Masjidil Haram. Berangkat ke Jeddah, kunjungan singkat ke Corniche dan Masjid Terapung. Check-in hotel transit Jeddah.", "photo_url": "https://images.unsplash.com/photo-1578895101408-1a36b834405b?w=800&h=450&fit=crop&auto=format"},
        {"day": 12, "title": "Kepulangan ke Jakarta", "description": "Check-out hotel pagi hari. Transfer ke Bandara King Abdulaziz Jeddah. Penerbangan pulang ke Jakarta. Tiba di CGK malam hari dengan membawa kenangan ibadah yang tak terlupakan.", "photo_url": "https://images.unsplash.com/photo-1436491865332-7a61a109cc05?w=800&h=450&fit=crop&auto=format"}
    ]'::jsonb,
    'https://umroh-os.example/itinerary/itn_01JCDG00000000000000000001'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.addons (id, name, list_amount, list_currency, settlement_currency)
VALUES
    ('addon_01JCDK00000000000000000001', 'Extra night Madinah',      2500000, 'IDR', 'IDR'),
    ('addon_01JCDK00000000000000000002', 'Upgrade room Double',      3000000, 'IDR', 'IDR')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------
-- Packages
-- ---------------------------------------------------------------------

-- Active — full tree. Surfaces in list + detail.
INSERT INTO catalog.packages (
    id, kind, name, description, highlights, cover_photo_url,
    itinerary_id, airline_id, muthawwif_id, status
) VALUES (
    'pkg_01JCDE00000000000000000001',
    'umrah_reguler',
    'Umrah Reguler 12 Hari — Ramadan 1447 H',
    'Paket umrah reguler 12 hari dengan keberangkatan di pertengahan Ramadan 1447 H. Hotel bintang 5 dekat Masjidil Haram dan Masjid Nabawi.',
    ARRAY['Direct Jakarta-Jeddah', 'Hotel 5-star Mecca 50m walking', 'Muthawwif S3 Al-Azhar'],
    'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=1280&h=720&fit=crop&auto=format',
    'itn_01JCDG00000000000000000001',
    'arl_01JCDI00000000000000000001',
    'mtw_01JCDJ00000000000000000001',
    'active'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.package_hotels (package_id, hotel_id, sort_order) VALUES
    ('pkg_01JCDE00000000000000000001', 'htl_01JCDH00000000000000000001', 0),
    ('pkg_01JCDE00000000000000000001', 'htl_01JCDH00000000000000000002', 1)
ON CONFLICT DO NOTHING;

INSERT INTO catalog.package_addons (package_id, addon_id) VALUES
    ('pkg_01JCDE00000000000000000001', 'addon_01JCDK00000000000000000001'),
    ('pkg_01JCDE00000000000000000001', 'addon_01JCDK00000000000000000002')
ON CONFLICT DO NOTHING;

-- Two upcoming departures: the earliest one drives `next_departure`.
INSERT INTO catalog.package_departures (
    id, package_id, departure_date, return_date,
    total_seats, reserved_seats, status
) VALUES
    ('dep_01JCDF00000000000000000001', 'pkg_01JCDE00000000000000000001', CURRENT_DATE + INTERVAL '30 days', CURRENT_DATE + INTERVAL '42 days', 45, 3, 'open'),
    ('dep_01JCDF00000000000000000002', 'pkg_01JCDE00000000000000000001', CURRENT_DATE + INTERVAL '75 days', CURRENT_DATE + INTERVAL '87 days', 45, 0, 'open')
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.package_pricing (
    id, package_departure_id, room_type, list_amount, list_currency, settlement_currency
) VALUES
    ('pkgpr_01JCDP00000000000000000001', 'dep_01JCDF00000000000000000001', 'quad',   38500000, 'IDR', 'IDR'),
    ('pkgpr_01JCDP00000000000000000002', 'dep_01JCDF00000000000000000001', 'triple', 41500000, 'IDR', 'IDR'),
    ('pkgpr_01JCDP00000000000000000003', 'dep_01JCDF00000000000000000001', 'double', 45500000, 'IDR', 'IDR'),
    ('pkgpr_01JCDP00000000000000000004', 'dep_01JCDF00000000000000000002', 'quad',   39500000, 'IDR', 'IDR'),
    ('pkgpr_01JCDP00000000000000000005', 'dep_01JCDF00000000000000000002', 'triple', 42500000, 'IDR', 'IDR'),
    ('pkgpr_01JCDP00000000000000000006', 'dep_01JCDF00000000000000000002', 'double', 46500000, 'IDR', 'IDR')
ON CONFLICT (id) DO NOTHING;

-- Draft — must be hidden from public list/detail.
INSERT INTO catalog.packages (id, kind, name, description, cover_photo_url, status)
VALUES (
    'pkg_01JCDE00000000000000000002',
    'hajj_khusus',
    'Hajj Khusus 2026 — Draft',
    'Paket haji khusus (belum dipublikasikan).',
    '',
    'draft'
) ON CONFLICT (id) DO NOTHING;

-- Archived — must also be hidden.
INSERT INTO catalog.packages (id, kind, name, description, cover_photo_url, status)
VALUES (
    'pkg_01JCDE00000000000000000003',
    'umrah_plus',
    'Umrah Plus Istanbul 14 Hari — arsip 2025',
    'Arsip paket 2025; tidak untuk dijual.',
    '',
    'archived'
) ON CONFLICT (id) DO NOTHING;
