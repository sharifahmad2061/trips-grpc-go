CREATE SCHEMA IF NOT EXISTS trip;

CREATE TABLE trip.trip (
   id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
   name VARCHAR(100) NOT NULL,
   member_id BIGINT NOT NULL,
   start_date DATE NOT NULL,
   end_date DATE NOT NULL
);

-- Create index on member_id and end_date for efficient queries
CREATE INDEX idx_trip_member_id_end_date ON trip.trip(member_id, end_date);

-- Add this for efficient deletion queries
CREATE INDEX idx_trip_end_date ON trip.trip(end_date);
