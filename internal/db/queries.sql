-- name: GetTripByID :one
SELECT *
FROM trip.trip
WHERE id = $1;

-- name: CreateTrip :one
INSERT INTO trip.trip (name, member_id, start_date, end_date)
VALUES ($1, $2, $3, $4)
RETURNING *;