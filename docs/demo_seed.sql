INSERT INTO tenants (name, code, plan, active, created_at)
VALUES
    ('Demo Energy', 'demo-energy', 'BUSINESS', TRUE, CURRENT_TIMESTAMP),
    ('North Grid', 'north-grid', 'FREE', TRUE, CURRENT_TIMESTAMP);

INSERT INTO users (username, password, full_name, role, tenant_id)
VALUES
    ('superadmin', '$2a$10$lhyAbyHikJ/Xk1jy5aG2aO9..tk1BdG.kitph3egf0hPuR0gwGDfm', 'Platform Super Admin', 'SUPER_ADMIN', NULL),
    ('tenantadmin', '$2a$10$sH20ZOmmodGSqu3x25.cAuyLyHOHbJWDZo5AymzhpbOzJKLLg79zi', 'Tenant Administrator', 'TENANT_ADMIN',
        (SELECT id FROM tenants WHERE code = 'demo-energy')),
    ('dispatcher1', '$2a$10$x.NfOMCoLm5A1.vy6NXFmOejRCi55y7goSYYFMX8XKTfTmoswMKhe', 'Main Dispatcher', 'DISPATCHER',
        (SELECT id FROM tenants WHERE code = 'demo-energy')),
    ('electrician1', '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Field Electrician', 'ELECTRICIAN',
        (SELECT id FROM tenants WHERE code = 'demo-energy')),
    ('northadmin', '$2a$10$sH20ZOmmodGSqu3x25.cAuyLyHOHbJWDZo5AymzhpbOzJKLLg79zi', 'North Grid Admin', 'TENANT_ADMIN',
        (SELECT id FROM tenants WHERE code = 'north-grid'));

INSERT INTO organizations (name, type, description, tenant_id)
VALUES
    ('Romashka LLC', 'COMMERCIAL', 'Commercial consumer for demo tenant',
        (SELECT id FROM tenants WHERE code = 'demo-energy')),
    ('North School', 'GOVERNMENT', 'Municipal consumer for north tenant',
        (SELECT id FROM tenants WHERE code = 'north-grid'));

INSERT INTO addresses (street, house, building, apartment, tenant_id, consumer_id)
VALUES
    ('Lenina street', '10', '1', '15',
        (SELECT id FROM tenants WHERE code = 'demo-energy'),
        (SELECT id FROM organizations WHERE name = 'Romashka LLC' AND tenant_id = (SELECT id FROM tenants WHERE code = 'demo-energy'))),
    ('Mira avenue', '25', NULL, NULL,
        (SELECT id FROM tenants WHERE code = 'demo-energy'),
        (SELECT id FROM organizations WHERE name = 'Romashka LLC' AND tenant_id = (SELECT id FROM tenants WHERE code = 'demo-energy'))),
    ('School lane', '7', NULL, NULL,
        (SELECT id FROM tenants WHERE code = 'north-grid'),
        (SELECT id FROM organizations WHERE name = 'North School' AND tenant_id = (SELECT id FROM tenants WHERE code = 'north-grid')));

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
VALUES
    ('INSPECTION',
        (SELECT id FROM tenants WHERE code = 'demo-energy'),
        (SELECT id FROM addresses WHERE street = 'Lenina street' AND tenant_id = (SELECT id FROM tenants WHERE code = 'demo-energy')),
        'PENDING',
        CURRENT_DATE + INTERVAL '2 day',
        (SELECT id FROM users WHERE username = 'electrician1' AND tenant_id = (SELECT id FROM tenants WHERE code = 'demo-energy'))),
    ('REPLACEMENT',
        (SELECT id FROM tenants WHERE code = 'demo-energy'),
        (SELECT id FROM addresses WHERE street = 'Mira avenue' AND tenant_id = (SELECT id FROM tenants WHERE code = 'demo-energy')),
        'IN_PROGRESS',
        CURRENT_DATE + INTERVAL '5 day',
        (SELECT id FROM users WHERE username = 'electrician1' AND tenant_id = (SELECT id FROM tenants WHERE code = 'demo-energy'))),
    ('INSPECTION',
        (SELECT id FROM tenants WHERE code = 'north-grid'),
        (SELECT id FROM addresses WHERE street = 'School lane' AND tenant_id = (SELECT id FROM tenants WHERE code = 'north-grid')),
        'PENDING',
        CURRENT_DATE + INTERVAL '3 day',
        NULL);
