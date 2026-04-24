-- 000027_fix_dev_seed_content.up.sql
--
-- Fixes dev-fixture content issues found during UAT:
--   ISSUE-020: cover_photo_url used a fake cdn.umroh-os.example domain → broken image
--   ISSUE-023: itinerary had only 3 days for a 12-day trip + fake photo URLs
--              muthawwif portrait URL also pointed to fake CDN
--
-- Uses publicly accessible Unsplash images (no auth required).
-- All UPDATE statements are idempotent — re-running is safe.

-- Fix muthawwif portrait (ISSUE-023 / fake CDN)
UPDATE catalog.muthawwif
SET portrait_url = 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=320&h=320&fit=crop&auto=format'
WHERE id = 'mtw_01JCDJ00000000000000000001'
  AND portrait_url LIKE '%umroh-os.example%';

-- Fix package cover photo (ISSUE-020)
UPDATE catalog.packages
SET cover_photo_url = 'https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=1280&h=720&fit=crop&auto=format'
WHERE id = 'pkg_01JCDE00000000000000000001'
  AND cover_photo_url LIKE '%umroh-os.example%';

-- Fix itinerary: expand from 3 placeholder days to full 12-day program (ISSUE-023)
UPDATE catalog.itinerary_templates
SET days = '[
  {
    "day": 1,
    "title": "Keberangkatan dari Jakarta",
    "description": "Berkumpul di Terminal 3 CGK pukul 18.00 WIB. Manasik ringkas, pemeriksaan dokumen, dan keberangkatan ke Jeddah dengan penerbangan langsung Garuda Indonesia.",
    "photo_url": "https://images.unsplash.com/photo-1436491865332-7a61a109cc05?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 2,
    "title": "Tiba di Madinah",
    "description": "Tiba di Bandara Prince Mohammad bin Abdulaziz pagi hari. Check-in hotel bintang 5 di sekitar Masjid Nabawi. Waktu istirahat dan shalat Asar berjamaah di Masjid Nabawi.",
    "photo_url": "https://images.unsplash.com/photo-1565552645632-d725f8bfc19a?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 3,
    "title": "Ziarah Madinah — Raudhah & Makam Baqi",
    "description": "Shalat Subuh di Masjid Nabawi. Kunjungan ke Raudhah (pagi). Ziarah ke Makam Baqi, Masjid Quba, Masjid Qiblatayn, dan Jabal Uhud.",
    "photo_url": "https://images.unsplash.com/photo-1580418827493-f2b22c0a76cb?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 4,
    "title": "Ibadah Madinah & Persiapan ke Makkah",
    "description": "Pagi hingga siang ibadah mandiri di Masjid Nabawi. Sore hari persiapan keberangkatan. Malam berangkat ke Makkah via bus dengan persinggahan di Bir Ali untuk miqat.",
    "photo_url": "https://images.unsplash.com/photo-1565552645632-d725f8bfc19a?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 5,
    "title": "Miqat & Pelaksanaan Umrah",
    "description": "Berihram di Bir Ali. Tiba di Makkah dan langsung menuju Masjidil Haram untuk pelaksanaan ibadah umrah: tawaf 7 putaran, sa''i Shafa–Marwah, dan tahallul.",
    "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 6,
    "title": "Ibadah di Masjidil Haram",
    "description": "Hari penuh ibadah di Masjidil Haram. Tawaf sunnah, membaca Al-Qur''an di depan Ka''bah, i''tikaf, dan doa. Waktu shalat berjamaah lima waktu bersama jamaah internasional.",
    "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 7,
    "title": "Ziarah Makkah",
    "description": "Ziarah ke tempat-tempat bersejarah di Makkah: Jabal Tsur, Jabal Nur (Gua Hira), Masjid Ji''ranah, dan Tan''im. Sore kembali ke hotel untuk istirahat.",
    "photo_url": "https://images.unsplash.com/photo-1519817650390-64a93db51149?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 8,
    "title": "Ibadah di Masjidil Haram & Belanja",
    "description": "Pagi ibadah di Masjidil Haram. Siang waktu bebas untuk berbelanja oleh-oleh di area Abraj Al-Bait dan Zamzam Tower. Malam tawaf sunnah.",
    "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 9,
    "title": "Ibadah Makkah — I''tikaf",
    "description": "Hari penuh untuk i''tikaf, membaca Al-Qur''an, dan shalat nafilah. Bimbingan doa khusus oleh Ustadz Ahmad Al-Azhar pada malam hari.",
    "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 10,
    "title": "Ibadah Makkah & Persiapan Pulang",
    "description": "Pagi ibadah di Masjidil Haram, tawaf wada'' bagi yang berencana pulang. Siang packing dan persiapan. Malam farewell dinner bersama rombongan.",
    "photo_url": "https://images.unsplash.com/photo-1591604129939-f1efa4d9f7fa?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 11,
    "title": "Tawaf Wada'' & Perjalanan ke Jeddah",
    "description": "Tawaf Wada'' (tawaf perpisahan) di Masjidil Haram. Berangkat ke Jeddah, kunjungan singkat ke Corniche dan Masjid Terapung. Check-in hotel transit Jeddah.",
    "photo_url": "https://images.unsplash.com/photo-1578895101408-1a36b834405b?w=800&h=450&fit=crop&auto=format"
  },
  {
    "day": 12,
    "title": "Kepulangan ke Jakarta",
    "description": "Check-out hotel pagi hari. Transfer ke Bandara King Abdulaziz Jeddah. Penerbangan pulang ke Jakarta. Tiba di CGK malam hari dengan membawa kenangan ibadah yang tak terlupakan.",
    "photo_url": "https://images.unsplash.com/photo-1436491865332-7a61a109cc05?w=800&h=450&fit=crop&auto=format"
  }
]'::jsonb
WHERE id = 'itn_01JCDG00000000000000000001';
