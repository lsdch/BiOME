-- name: ListAbioticParams :many
SELECT *
FROM abiotic_params
ORDER BY name ASC;

-- name: CreateAbioticParam :one
INSERT INTO abiotic_params (
        name,
        code,
        description,
        unit
    )
VALUES (
        @name,
        @code,
        @description,
        @unit
    )
RETURNING *;

-- name: CreateAbioticMeasurement :one
INSERT INTO abiotic_measurements (
        sampling_id,
        param_id,
        value
    )
VALUES (
        @sampling_id,
        @abiotic_param_id,
        @value
    )
RETURNING *;