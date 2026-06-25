CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(150) NOT NULL,
    destination VARCHAR(150),
    start_date DATE,
    end_date DATE,
    timezone VARCHAR(64) NOT NULL DEFAULT 'Asia/Jakarta',
    currency_code VARCHAR(3) NOT NULL DEFAULT 'IDR',
    status VARCHAR(30) NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS trip_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES trips(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(30) NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_trip_members_unique_active ON trip_members(trip_id, user_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS itinerary_days (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES trips(id),
    date DATE,
    day_number INT NOT NULL,
    title VARCHAR(150),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_itinerary_days_unique_active ON itinerary_days(trip_id, day_number) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS itinerary_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id UUID NOT NULL REFERENCES itinerary_days(id),
    title VARCHAR(150) NOT NULL,
    description TEXT,
    location_name VARCHAR(200),
    latitude NUMERIC(10,7),
    longitude NUMERIC(10,7),
    start_time TIME,
    end_time TIME,
    cost_estimate NUMERIC(14,2) DEFAULT 0,
    sort_order INT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_itinerary_items_unique_active ON itinerary_items(day_id, sort_order) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_trips_owner_id ON trips(owner_id);
CREATE INDEX IF NOT EXISTS idx_trips_created_by ON trips(created_by);
CREATE INDEX IF NOT EXISTS idx_trips_updated_by ON trips(updated_by);
CREATE INDEX IF NOT EXISTS idx_trips_deleted_by ON trips(deleted_by);
CREATE INDEX IF NOT EXISTS idx_trip_members_trip_id ON trip_members(trip_id);
CREATE INDEX IF NOT EXISTS idx_trip_members_user_id ON trip_members(user_id);
CREATE INDEX IF NOT EXISTS idx_trip_members_created_by ON trip_members(created_by);
CREATE INDEX IF NOT EXISTS idx_trip_members_updated_by ON trip_members(updated_by);
CREATE INDEX IF NOT EXISTS idx_trip_members_deleted_by ON trip_members(deleted_by);
CREATE INDEX IF NOT EXISTS idx_itinerary_days_trip_id ON itinerary_days(trip_id);
CREATE INDEX IF NOT EXISTS idx_itinerary_days_created_by ON itinerary_days(created_by);
CREATE INDEX IF NOT EXISTS idx_itinerary_days_updated_by ON itinerary_days(updated_by);
CREATE INDEX IF NOT EXISTS idx_itinerary_days_deleted_by ON itinerary_days(deleted_by);
CREATE INDEX IF NOT EXISTS idx_itinerary_items_day_id ON itinerary_items(day_id);
CREATE INDEX IF NOT EXISTS idx_itinerary_items_created_by ON itinerary_items(created_by);
CREATE INDEX IF NOT EXISTS idx_itinerary_items_updated_by ON itinerary_items(updated_by);
CREATE INDEX IF NOT EXISTS idx_itinerary_items_deleted_by ON itinerary_items(deleted_by);

INSERT INTO permissions (id, name, display_name, resource, action) VALUES
    (gen_random_uuid(), 'trips:list', 'List Trips', 'trips', 'list'),
    (gen_random_uuid(), 'trips:create', 'Create Trip', 'trips', 'create'),
    (gen_random_uuid(), 'trips:view', 'View Trip', 'trips', 'view'),
    (gen_random_uuid(), 'trips:update', 'Update Trip', 'trips', 'update'),
    (gen_random_uuid(), 'trips:delete', 'Delete Trip', 'trips', 'delete'),
    (gen_random_uuid(), 'trips:manage_members', 'Manage Trip Members', 'trips', 'manage_members'),
    (gen_random_uuid(), 'itineraries:create', 'Create Itinerary', 'itineraries', 'create'),
    (gen_random_uuid(), 'itineraries:update', 'Update Itinerary', 'itineraries', 'update'),
    (gen_random_uuid(), 'itineraries:delete', 'Delete Itinerary', 'itineraries', 'delete')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.resource IN ('trips', 'itineraries')
WHERE r.name IN ('admin', 'superadmin')
ON CONFLICT DO NOTHING;
