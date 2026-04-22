-- Reverse of 000009_seed_catalog_dev_fixtures. Removes only the
-- explicit seed IDs so hand-inserted rows survive a rollback.

DELETE FROM catalog.package_pricing WHERE id IN (
    'pkgpr_01JCDP00000000000000000001',
    'pkgpr_01JCDP00000000000000000002',
    'pkgpr_01JCDP00000000000000000003',
    'pkgpr_01JCDP00000000000000000004',
    'pkgpr_01JCDP00000000000000000005',
    'pkgpr_01JCDP00000000000000000006'
);

DELETE FROM catalog.package_departures WHERE id IN (
    'dep_01JCDF00000000000000000001',
    'dep_01JCDF00000000000000000002'
);

DELETE FROM catalog.package_addons WHERE package_id IN (
    'pkg_01JCDE00000000000000000001'
);

DELETE FROM catalog.package_hotels WHERE package_id IN (
    'pkg_01JCDE00000000000000000001'
);

DELETE FROM catalog.packages WHERE id IN (
    'pkg_01JCDE00000000000000000001',
    'pkg_01JCDE00000000000000000002',
    'pkg_01JCDE00000000000000000003'
);

DELETE FROM catalog.addons WHERE id IN (
    'addon_01JCDK00000000000000000001',
    'addon_01JCDK00000000000000000002'
);

DELETE FROM catalog.itinerary_templates WHERE id IN (
    'itn_01JCDG00000000000000000001'
);

DELETE FROM catalog.muthawwif WHERE id IN (
    'mtw_01JCDJ00000000000000000001'
);

DELETE FROM catalog.airlines WHERE id IN (
    'arl_01JCDI00000000000000000001'
);

DELETE FROM catalog.hotels WHERE id IN (
    'htl_01JCDH00000000000000000001',
    'htl_01JCDH00000000000000000002'
);
