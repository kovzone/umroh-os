-- 000027_fix_dev_seed_content.down.sql
-- Reverts cover photo, muthawwif portrait, and itinerary back to the
-- original (broken) placeholder values from migration 000009.
-- Used only in development when rolling back this specific migration.

UPDATE catalog.muthawwif
SET portrait_url = 'https://cdn.umroh-os.example/mtw/ahmad.jpg'
WHERE id = 'mtw_01JCDJ00000000000000000001';

UPDATE catalog.packages
SET cover_photo_url = 'https://cdn.umroh-os.example/pkg/01JCDE.../cover.jpg'
WHERE id = 'pkg_01JCDE00000000000000000001';

UPDATE catalog.itinerary_templates
SET days = '[
    {"day": 1, "title": "Keberangkatan dari Jakarta", "description": "Berkumpul di CGK, manasik ringkas, keberangkatan ke Jeddah.", "photo_url": "https://cdn.umroh-os.example/itn/day1.jpg"},
    {"day": 2, "title": "Ziarah Madinah", "description": "Masjid Nabawi, Raudhah, kota Madinah.", "photo_url": "https://cdn.umroh-os.example/itn/day2.jpg"},
    {"day": 3, "title": "Ibadah Umrah", "description": "Miqat Bir Ali, tawaf, sa''i, tahallul."}
]'::jsonb
WHERE id = 'itn_01JCDG00000000000000000001';
