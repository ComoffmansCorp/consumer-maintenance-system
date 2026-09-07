-- Images for the marketplace: a thumbnail per service, an avatar per master
-- profile. Both nullable -- not every service/master has one (e.g. seeded
-- before this migration existed), rendered as a placeholder client-side.
ALTER TABLE services ADD COLUMN image_url TEXT NULL;
ALTER TABLE master_profiles ADD COLUMN avatar_url TEXT NULL;
