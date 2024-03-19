-- Insert data into tickets table
INSERT INTO tickets (id, capacity, region, level, event_date, created_at, updated_at, deleted_at)
VALUES
    (1, 420000 * 0.6, 'asean', 'wood', '2022-01-01', NOW(), NULL, NULL),
    (2, 420000 * 0.25, 'asean', 'bronze', '2022-01-01', NOW(), NULL, NULL),
    (3, 420000 * 0.1, 'asean', 'silver', '2022-01-01', NOW(), NULL, NULL),
    (4, 420000 * 0.05, 'asean', 'gold', '2022-01-01', NOW(), NULL, NULL),
    (5, 80000 * 0.6, 'asia', 'wood', '2022-01-01', NOW(), NULL, NULL),
    (6, 80000 * 0.25, 'asia', 'bronze', '2022-01-01', NOW(), NULL, NULL),
    (7, 80000 * 0.1, 'asia', 'silver', '2022-01-01', NOW(), NULL, NULL),
    (8, 80000 * 0.05, 'asia', 'gold', '2022-01-01', NOW(), NULL, NULL),
    (9, 70000 * 0.6, 'middle_east', 'wood', '2022-01-01', NOW(), NULL, NULL),
    (10, 70000 * 0.25, 'middle_east', 'bronze', '2022-01-01', NOW(), NULL, NULL),
    (11, 70000 * 0.1, 'middle_east', 'silver', '2022-01-01', NOW(), NULL, NULL),
    (12, 70000 * 0.05, 'middle_east', 'gold', '2022-01-01', NOW(), NULL, NULL),
    (13, 90000 * 0.6, 'south_america', 'wood', '2022-01-01', NOW(), NULL, NULL),
    (14, 90000 * 0.25, 'south_america', 'bronze', '2022-01-01', NOW(), NULL, NULL),
    (15, 90000 * 0.1, 'south_america', 'silver', '2022-01-01', NOW(), NULL, NULL),
    (16, 90000 * 0.05, 'south_america', 'gold', '2022-01-01', NOW(), NULL, NULL),
    (17, 340000, 'online', 'online', '2022-01-01', NOW(), NULL, NULL);

-- Insert data into ticket_details table
INSERT INTO ticket_details (id, ticket_id, base_price, created_at, updated_at, deleted_at)
VALUES
    (1, 1, 50, NOW(), NULL, NULL),
    (2, 2, 75, NOW(), NULL, NULL),
    (3, 3, 75, NOW(), NULL, NULL),
    (4, 4, 100, NOW(), NULL, NULL),
    (5, 5, 150, NOW(), NULL, NULL),
    (6, 6, 175, NOW(), NULL, NULL),
    (7, 7, 175, NOW(), NULL, NULL),
    (8, 8, 200, NOW(), NULL, NULL),
    (9, 9, 250, NOW(), NULL, NULL),
    (10, 10, 325, NOW(), NULL, NULL),
    (11, 11, 325, NOW(), NULL, NULL),
    (12, 12, 325, NOW(), NULL, NULL),
    (13, 13, 350, NOW(), NULL, NULL),
    (14, 14, 425, NOW(), NULL, NULL),
    (15, 15, 425, NOW(), NULL, NULL),
    (16, 16, 450, NOW(), NULL, NULL),
    (17, 17, 50, NOW(), NULL, NULL);