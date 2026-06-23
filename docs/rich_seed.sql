-- =============================================================================
-- RICH DEMO SEED  —  Consumer Maintenance System
-- =============================================================================
-- Запускать ПОСЛЕ старта приложения (Hibernate создаст таблицы).
-- Пароль для всех пользователей ниже: тот же, что у electrician1 из demo_seed.sql
--   Чтобы создать superadmin с известным паролем, вызови:
--   POST /api/auth/bootstrap-super-admin
--   {"username":"superadmin","password":"super123"}
-- =============================================================================

-- ── Tenants ──────────────────────────────────────────────────────────────────

INSERT INTO tenants (name, code, plan, active, created_at) VALUES
    ('ЭнергоСервис Урал', 'esc-ural',     'ENTERPRISE', TRUE, NOW() - INTERVAL '90 day'),
    ('ГазМонтаж Плюс',    'gazmontazh',   'BUSINESS',   TRUE, NOW() - INTERVAL '60 day'),
    ('СтройКонтроль',     'stroykontrol', 'BUSINESS',   TRUE, NOW() - INTERVAL '45 day'),
    ('МегаВатт Сервис',   'megawatt',     'ENTERPRISE', TRUE, NOW() - INTERVAL '120 day'),
    ('ЖКХ Партнёр',       'zhkh-partner', 'FREE',       TRUE, NOW() - INTERVAL '30 day'),
    ('North Grid',        'north-grid',   'FREE',       TRUE, NOW() - INTERVAL '15 day');

-- ── Users ─────────────────────────────────────────────────────────────────────
-- Используем тот же bcrypt хэш, что и у electrician1 из demo_seed.sql

INSERT INTO users (username, password, full_name, role, tenant_id) VALUES
-- esc-ural (ENTERPRISE)
('ivanov_p',   '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Иванов Павел Алексеевич',     'TENANT_ADMIN', (SELECT id FROM tenants WHERE code='esc-ural')),
('smirnova_o', '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Смирнова Оксана Дмитриевна',  'DISPATCHER',   (SELECT id FROM tenants WHERE code='esc-ural')),
('kozlov_d',   '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Козлов Дмитрий Игоревич',     'ELECTRICIAN',  (SELECT id FROM tenants WHERE code='esc-ural')),
('fedorova_a', '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Фёдорова Анна Сергеевна',     'ELECTRICIAN',  (SELECT id FROM tenants WHERE code='esc-ural')),
('volkov_r',   '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Волков Роман Николаевич',     'ELECTRICIAN',  (SELECT id FROM tenants WHERE code='esc-ural')),
-- gazmontazh (BUSINESS)
('petrov_a',   '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Петров Алексей Витальевич',   'TENANT_ADMIN', (SELECT id FROM tenants WHERE code='gazmontazh')),
('novikov_s',  '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Новиков Сергей Павлович',     'ELECTRICIAN',  (SELECT id FROM tenants WHERE code='gazmontazh')),
('larina_m',   '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Ларина Марина Юрьевна',       'DISPATCHER',   (SELECT id FROM tenants WHERE code='gazmontazh')),
-- stroykontrol (BUSINESS)
('morozova_e', '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Морозова Елена Владимировна', 'TENANT_ADMIN', (SELECT id FROM tenants WHERE code='stroykontrol')),
('titov_k',    '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Титов Константин Андреевич',  'ELECTRICIAN',  (SELECT id FROM tenants WHERE code='stroykontrol')),
-- megawatt (ENTERPRISE)
('belov_v',    '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Белов Виктор Михайлович',     'TENANT_ADMIN', (SELECT id FROM tenants WHERE code='megawatt')),
('gorbunov_n', '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Горбунов Никита Олегович',    'ELECTRICIAN',  (SELECT id FROM tenants WHERE code='megawatt')),
('savina_t',   '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Савина Татьяна Ивановна',     'DISPATCHER',   (SELECT id FROM tenants WHERE code='megawatt')),
('ershov_p',   '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Ершов Павел Геннадьевич',     'ELECTRICIAN',  (SELECT id FROM tenants WHERE code='megawatt')),
-- zhkh-partner (FREE)
('zhkh_admin', '$2a$10$xSwqxGHdbH6L7tCahm91SuEt0eAv0JaJBJ5hw1PM7CcgdqpqYylF.', 'Голубев Иван Петрович',       'TENANT_ADMIN', (SELECT id FROM tenants WHERE code='zhkh-partner')),
-- north-grid (FREE)
('northadmin', '$2a$10$sH20ZOmmodGSqu3x25.cAuyLyHOHbJWDZo5AymzhpbOzJKLLg79zi',  'North Grid Admin',            'TENANT_ADMIN', (SELECT id FROM tenants WHERE code='north-grid'));

-- ── Organizations ─────────────────────────────────────────────────────────────

INSERT INTO organizations (name, type, description, tenant_id) VALUES
    ('ОАО ЭнергоПоставка',  'COMMERCIAL',   'Промышленный потребитель', (SELECT id FROM tenants WHERE code='esc-ural')),
    ('Школа №14',           'GOVERNMENT',   'Муниципальная школа',      (SELECT id FROM tenants WHERE code='esc-ural')),
    ('ТЦ Мегаполис',        'COMMERCIAL',   'Торговый центр',           (SELECT id FROM tenants WHERE code='esc-ural')),
    ('Больница №2',         'GOVERNMENT',   'Городская больница',       (SELECT id FROM tenants WHERE code='esc-ural')),
    ('ГазТрейд ООО',        'COMMERCIAL',   'Газовый поставщик',        (SELECT id FROM tenants WHERE code='gazmontazh')),
    ('МКД Ленина 15',       'RESIDENTIAL',  'Многоквартирный дом',      (SELECT id FROM tenants WHERE code='gazmontazh')),
    ('ЖК Панорама',         'RESIDENTIAL',  'Жилой комплекс',           (SELECT id FROM tenants WHERE code='stroykontrol')),
    ('North School',        'GOVERNMENT',   'Municipal consumer',       (SELECT id FROM tenants WHERE code='north-grid')),
    ('Завод Прогресс',      'COMMERCIAL',   'Завод приборостроения',    (SELECT id FROM tenants WHERE code='megawatt')),
    ('МегаЭнерго',          'COMMERCIAL',   'Крупный промышленник',     (SELECT id FROM tenants WHERE code='megawatt'));

-- ── Addresses ─────────────────────────────────────────────────────────────────

INSERT INTO addresses (street, house, building, apartment, tenant_id, consumer_id) VALUES
-- esc-ural
('ул. Ленина',       '12',  NULL, '34', (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='ОАО ЭнергоПоставка' LIMIT 1)),
('ул. Ленина',       '12',  NULL, '7',  (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='ОАО ЭнергоПоставка' LIMIT 1)),
('пр. Мира',         '5',   '1',  '7',  (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='ОАО ЭнергоПоставка' LIMIT 1)),
('ул. Гагарина',     '3',   NULL, '15', (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='Школа №14' LIMIT 1)),
('ул. Гагарина',     '44',  NULL, '24', (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='ТЦ Мегаполис' LIMIT 1)),
('ул. Советская',    '88',  NULL, '2',  (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='Больница №2' LIMIT 1)),
('ул. Пушкина',      '21',  NULL, NULL, (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='ТЦ Мегаполис' LIMIT 1)),
('пр. Победы',       '44',  NULL, '9',  (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='Школа №14' LIMIT 1)),
('ул. Кирова',       '17',  '2',  '33', (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='ОАО ЭнергоПоставка' LIMIT 1)),
('ул. Строителей',   '9',   NULL, '11', (SELECT id FROM tenants WHERE code='esc-ural'), (SELECT id FROM organizations WHERE name='Больница №2' LIMIT 1)),
-- gazmontazh
('пр. Октября',      '33',  NULL, '5',  (SELECT id FROM tenants WHERE code='gazmontazh'), (SELECT id FROM organizations WHERE name='ГазТрейд ООО' LIMIT 1)),
('ул. Молодёжная',   '8',   '1',  NULL, (SELECT id FROM tenants WHERE code='gazmontazh'), (SELECT id FROM organizations WHERE name='МКД Ленина 15' LIMIT 1)),
-- stroykontrol
('ул. Новостроек',   '1',   'А',  '42', (SELECT id FROM tenants WHERE code='stroykontrol'), (SELECT id FROM organizations WHERE name='ЖК Панорама' LIMIT 1)),
-- north-grid
('School lane',      '7',   NULL, NULL, (SELECT id FROM tenants WHERE code='north-grid'),    (SELECT id FROM organizations WHERE name='North School' LIMIT 1)),
-- megawatt
('пр. Автозаводцев', '100', NULL, NULL, (SELECT id FROM tenants WHERE code='megawatt'), (SELECT id FROM organizations WHERE name='Завод Прогресс' LIMIT 1)),
('ул. Промышленная', '55',  '3',  NULL, (SELECT id FROM tenants WHERE code='megawatt'), (SELECT id FROM organizations WHERE name='МегаЭнерго' LIMIT 1));

-- ── Tasks ─────────────────────────────────────────────────────────────────────

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'COMPLETED', CURRENT_DATE - 15, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Ленина' AND a.apartment='34' AND a.tenant_id=t.id AND u.username='kozlov_d';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'COMPLETED', CURRENT_DATE - 10, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='пр. Мира' AND a.tenant_id=t.id AND u.username='fedorova_a';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'COMPLETED', CURRENT_DATE - 8, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Гагарина' AND a.apartment='15' AND a.tenant_id=t.id AND u.username='kozlov_d';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'COMPLETED', CURRENT_DATE - 5, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Ленина' AND a.apartment='7' AND a.tenant_id=t.id AND u.username='volkov_r';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'COMPLETED', CURRENT_DATE - 3, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Кирова' AND a.tenant_id=t.id AND u.username='fedorova_a';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'IN_PROGRESS', CURRENT_DATE + 2, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Гагарина' AND a.apartment='24' AND a.tenant_id=t.id AND u.username='kozlov_d';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'IN_PROGRESS', CURRENT_DATE + 3, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Советская' AND a.tenant_id=t.id AND u.username='fedorova_a';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'IN_PROGRESS', CURRENT_DATE + 4, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='пр. Победы' AND a.tenant_id=t.id AND u.username='volkov_r';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'PENDING', CURRENT_DATE + 5, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Пушкина' AND a.tenant_id=t.id AND u.username='kozlov_d';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'PENDING', CURRENT_DATE + 7, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Строителей' AND a.tenant_id=t.id AND u.username='fedorova_a';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'PENDING', CURRENT_DATE + 9, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='ул. Ленина' AND a.apartment='34' AND a.tenant_id=t.id AND u.username='volkov_r';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'CANCELED', CURRENT_DATE - 1, u.id
FROM tenants t, addresses a, users u
WHERE t.code='esc-ural' AND a.street='пр. Победы' AND a.tenant_id=t.id AND u.username='volkov_r';

-- gazmontazh
INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'COMPLETED', CURRENT_DATE - 7, u.id
FROM tenants t, addresses a, users u
WHERE t.code='gazmontazh' AND a.street='пр. Октября' AND a.tenant_id=t.id AND u.username='novikov_s';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'IN_PROGRESS', CURRENT_DATE + 3, u.id
FROM tenants t, addresses a, users u
WHERE t.code='gazmontazh' AND a.street='ул. Молодёжная' AND a.tenant_id=t.id AND u.username='novikov_s';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'PENDING', CURRENT_DATE + 6, u.id
FROM tenants t, addresses a, users u
WHERE t.code='gazmontazh' AND a.street='пр. Октября' AND a.tenant_id=t.id AND u.username='novikov_s';

-- stroykontrol
INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'PENDING', CURRENT_DATE + 4, u.id
FROM tenants t, addresses a, users u
WHERE t.code='stroykontrol' AND a.street='ул. Новостроек' AND a.tenant_id=t.id AND u.username='titov_k';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'COMPLETED', CURRENT_DATE - 2, u.id
FROM tenants t, addresses a, users u
WHERE t.code='stroykontrol' AND a.street='ул. Новостроек' AND a.tenant_id=t.id AND u.username='titov_k';

-- megawatt
INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'COMPLETED', CURRENT_DATE - 20, u.id
FROM tenants t, addresses a, users u
WHERE t.code='megawatt' AND a.street='пр. Автозаводцев' AND a.tenant_id=t.id AND u.username='gorbunov_n';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'IN_PROGRESS', CURRENT_DATE + 1, u.id
FROM tenants t, addresses a, users u
WHERE t.code='megawatt' AND a.street='ул. Промышленная' AND a.tenant_id=t.id AND u.username='ershov_p';

INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'REPLACEMENT', t.id, a.id, 'PENDING', CURRENT_DATE + 8, u.id
FROM tenants t, addresses a, users u
WHERE t.code='megawatt' AND a.street='пр. Автозаводцев' AND a.tenant_id=t.id AND u.username='gorbunov_n';

-- north-grid
INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
SELECT 'INSPECTION', t.id, a.id, 'PENDING', CURRENT_DATE + 3, NULL
FROM tenants t, addresses a
WHERE t.code='north-grid' AND a.street='School lane' AND a.tenant_id=t.id;

-- ── Inspection Acts ───────────────────────────────────────────────────────────

INSERT INTO inspection_acts (task_id, tenant_id, address_id, inspection_date, consumer_id, inspection_type, notes)
SELECT t.id, t.tenant_id, t.address_id, t.due_date,
       a.consumer_id, 'SCHEDULED', 'Плановый осмотр выполнен. Нарушений не выявлено.'
FROM tasks t
JOIN addresses a ON a.id = t.address_id
WHERE t.status = 'COMPLETED' AND t.type = 'INSPECTION';

-- ── Meters for inspection acts ────────────────────────────────────────────────

INSERT INTO meters (type, serial_number, manufacture_year, verification_date, seal_state, inspection_act_id)
SELECT 'SINGLE_PHASE',
       CONCAT('CE102-', LPAD(ia.id::TEXT, 8, '0')),
       2019,
       ia.inspection_date + INTERVAL '2 year',
       'Не сорвана',
       ia.id
FROM inspection_acts ia;

INSERT INTO meters (type, serial_number, manufacture_year, verification_date, seal_state, inspection_act_id)
SELECT 'THREE_PHASE_DIRECT',
       CONCAT('M234-', LPAD(ia.id::TEXT, 8, '0')),
       2020,
       ia.inspection_date + INTERVAL '1 year',
       'Сорвана',
       ia.id
FROM inspection_acts ia
LIMIT 2;

-- ── Replacement Acts ──────────────────────────────────────────────────────────

INSERT INTO replacement_acts (
    task_id, tenant_id, address_id, account_number, installation_date,
    old_brand, old_serial_number, old_readings,
    new_brand, new_serial_number, new_readings
)
SELECT
    t.id,
    t.tenant_id,
    t.address_id,
    CONCAT('ЛС-', LPAD(t.id::TEXT, 6, '0')),
    t.due_date,
    'Энергомера CE102',  CONCAT('OLD-', LPAD(t.id::TEXT, 8, '0')), 4512.75 + t.id,
    'Меркурий 234',      CONCAT('NEW-', LPAD(t.id::TEXT, 8, '0')), 0.00
FROM tasks t
WHERE t.status = 'COMPLETED' AND t.type = 'REPLACEMENT';
