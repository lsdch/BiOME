-- name: DetectExistingSamplings :many 
SELECT latitude,
    longitude,
    event_date,
    event_date_precision,
    coordinates_precision,
    ST_Distance(
        coordinates::geography,
        ST_SetSRID(
            ST_MakePoint(@longitude::float8, @latitude::float8),
            4326
        )::geography
    )::integer AS distance_meters
FROM samplings
WHERE ST_DWithin(
        coordinates::geography,
        ST_SetSRID(
            ST_MakePoint(@longitude::float8, @latitude::float8),
            4326
        )::geography,
        @radius_meters::integer
    )
    AND (
        event_date_precision IS NOT NULL
        AND event_date IS NOT NULL
        AND event_date BETWEEN (
            @event_date::date - INTERVAL @date_interval_days::integer || 'day'
        )
        AND (
            @event_date::date + INTERVAL @date_interval_days::integer || 'day'
        )
    );