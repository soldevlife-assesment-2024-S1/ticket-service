-- Insert data into tickets table
INSERT INTO tickets (id, capacity, region, event_date, created_at, updated_at, deleted_at)
VALUES
    (1, 420000, 'asean', '2022-01-01', NOW(), NULL, NULL),
    (2, 80000, 'asia', '2022-01-01', NOW(), NULL, NULL),
    (3, 70000, 'middle_east', '2022-01-01', NOW(), NULL, NULL),
    (4, 90000, 'south_america', '2022-01-01', NOW(), NULL, NULL),
    (5, 340000, 'online', '2022-01-01', NOW(), NULL, NULL);

-- Insert data into ticket_details table
INSERT INTO ticket_details (id, ticket_id, level, base_price, stock, created_at, updated_at, deleted_at)
VALUES
    (1, 1, 'wood', 50,  420000 * 0.6, NOW(), NULL, NULL),
    (2, 1, 'bronze', 75, 420000 * 0.25, NOW(), NULL, NULL),
    (3, 1, 'silver', 75, 420000 * 0.1, NOW(), NULL, NULL),
    (4, 1, 'gold', 100, 420000 * 0.05, NOW(), NULL, NULL),
    (5, 2, 'wood', 150, 80000 * 0.6, NOW(), NULL, NULL),
    (6, 2, 'bronze', 175, 80000 * 0.25, NOW(), NULL, NULL),
    (7, 2, 'silver', 175, 80000 * 0.1, NOW(), NULL, NULL),
    (8, 2, 'gold', 200, 80000 * 0.05, NOW(), NULL, NULL),
    (9, 3, 'wood', 250, 70000 * 0.6, NOW(), NULL, NULL),
    (10, 3, 'bronze', 325, 70000 * 0.25, NOW(), NULL, NULL),
    (11, 3, 'silver', 325, 70000 * 0.1, NOW(), NULL, NULL),
    (12, 3, 'gold', 325, 70000 * 0.05, NOW(), NULL, NULL),
    (13, 4, 'wood', 350, 90000 * 0.6, NOW(), NULL, NULL),
    (14, 4, 'bronze', 425, 90000 * 0.25, NOW(), NULL, NULL),
    (15, 4, 'silver', 425,  90000 * 0.1, NOW(), NULL, NULL),
    (16, 4, 'gold', 450, 90000 * 0.05, NOW(), NULL, NULL),
    (17, 5, 'online', 50, 340000,  NOW(), NULL, NULL);