DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE resource IN ('trips', 'itineraries')
);

DELETE FROM permissions WHERE resource IN ('trips', 'itineraries');

DROP TABLE IF EXISTS itinerary_items;
DROP TABLE IF EXISTS itinerary_days;
DROP TABLE IF EXISTS trip_members;
DROP TABLE IF EXISTS trips;
