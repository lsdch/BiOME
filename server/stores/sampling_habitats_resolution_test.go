package stores

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/stretchr/testify/require"
)

// Only the related production schemas are installed in a disposable schema.
// BIOME_TEST_DATABASE_URL must reference PostgreSQL with citext installed.
func habitatResolutionTestDB(t *testing.T) (*db.DB, *pgxpool.Pool) {
	t.Helper()
	dsn := os.Getenv("BIOME_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("BIOME_TEST_DATABASE_URL is required")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(admin.Close)
	var citext bool
	require.NoError(t, admin.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'citext')").Scan(&citext))
	require.True(t, citext, "the test database requires citext")
	schema := "habitat_resolution_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	_, err = admin.Exec(ctx, "CREATE SCHEMA "+schema)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
		require.NoError(t, err)
	})
	cfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	cfg.ConnConfig.RuntimeParams["search_path"] = schema + ",public"
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	database := db.NewDB(pool)
	t.Cleanup(database.Close)
	_, err = pool.Exec(ctx, `
        CREATE TYPE import_batch_status AS ENUM ('created', 'staged', 'completed', 'failed', 'canceled');
        CREATE TABLE import_batches (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), status import_batch_status NOT NULL DEFAULT 'staged');
        CREATE TABLE samplings (id UUID PRIMARY KEY DEFAULT gen_random_uuid());
        CREATE TABLE sampling_methods (id UUID PRIMARY KEY);
        CREATE TABLE fixatives (id UUID PRIMARY KEY);
        CREATE TABLE taxon_resolution (id UUID PRIMARY KEY, import_id UUID, resolved_candidate_id UUID);
        CREATE TABLE publication_resolution (import_id UUID, status TEXT);
        CREATE TABLE import_samplings_occurrences (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            import_id UUID NOT NULL REFERENCES import_batches(id),
            sampling_hash TEXT NOT NULL,
            habitats CITEXT[],
            materialized_sampling_id UUID REFERENCES samplings(id)
        );
    `)
	require.NoError(t, err)
	for _, path := range []string{"../db/migrations/0022_habitats.sql", "../db/migrations/0032_sampling_metadata_resolution.sql"} {
		sql, err := os.ReadFile(path)
		require.NoError(t, err)
		_, err = pool.Exec(ctx, string(sql))
		require.NoError(t, err)
	}
	return database, pool
}

func habitatBatch(t *testing.T, database *db.DB, values ...[]string) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	var batchID, samplingID uuid.UUID
	require.NoError(t, database.QueryRow(ctx, "INSERT INTO import_batches DEFAULT VALUES RETURNING id").Scan(&batchID))
	require.NoError(t, database.QueryRow(ctx, "INSERT INTO samplings DEFAULT VALUES RETURNING id").Scan(&samplingID))
	for _, habitats := range values {
		require.NoError(t, database.QueryRow(ctx, `
            INSERT INTO import_samplings_occurrences (import_id, sampling_hash, habitats, materialized_sampling_id)
            VALUES ($1, 'shared_sampling', $2::text[]::citext[], $3) RETURNING id`, batchID, habitats, samplingID).Scan(new(uuid.UUID)))
	}
	store := NewSamplingStore()
	require.NoError(t, database.WithTx(ctx, func(tx *db.Tx) error {
		return store.InitHabitatResolution(ctx, tx, batchID)
	}))
	return batchID, samplingID
}

func TestSamplingHabitatResolutionIntegration(t *testing.T) {
	database, pool := habitatResolutionTestDB(t)
	ctx := context.Background()
	store := NewSamplingStore()
	lotic, err := database.Queries().InsertHabitat(ctx, biomedb.InsertHabitatParams{Label: "Lotic"})
	require.NoError(t, err)
	fresh, err := database.Queries().InsertHabitat(ctx, biomedb.InsertHabitatParams{Label: "Freshwater"})
	require.NoError(t, err)
	batch, sampling := habitatBatch(t, database, []string{"LOTIC", "river"}, []string{"Freshwater", "river", "ignore", "", " "})
	other, _ := habitatBatch(t, database, []string{"river"})
	ready, err := database.Queries().CheckSamplingHabitatsReady(ctx, batch)
	require.NoError(t, err)
	require.False(t, ready)
	check, err := database.Queries().CheckReadyToMaterialize(ctx, batch)
	require.NoError(t, err)
	require.False(t, models.MaterializationReadyCheckFromDB(check).IsReady())

	decision := models.SamplingHabitatResolutionInput{
		InputText: "river", Status: biomedb.VocabResolutionStatusRequestCreation,
		HabitatIDs:     []uuid.UUID{lotic.ID},
		CreateHabitats: []models.HabitatInput{{Name: "Cave", Facet: models.NewOptional("New facet"), Description: models.NewOptional("Requested description")}},
	}
	_, err = store.ResolveHabitats(ctx, database, batch, decision)
	require.NoError(t, err)
	var count int
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM habitats WHERE label = 'Cave'").Scan(&count))
	require.Zero(t, count, "creation must remain deferred")
	_, err = store.ResolveHabitats(ctx, database, batch, models.SamplingHabitatResolutionInput{
		InputText: "ignore", Status: biomedb.VocabResolutionStatusDiscard,
	})
	require.NoError(t, err)
	require.NoError(t, database.WithTx(ctx, func(tx *db.Tx) error { return store.InitHabitatResolution(ctx, tx, batch) }))
	resolutions, err := store.GetHabitatsResolution(ctx, database, batch)
	require.NoError(t, err)
	require.Len(t, resolutions, 4, "blank tokens must not need resolution")
	for _, resolution := range resolutions {
		if resolution.InputText == "river" {
			require.Equal(t, decision.Status, resolution.Status)
			require.Equal(t, decision.HabitatIDs, resolution.HabitatIDs)
			require.Equal(t, decision.CreateHabitats, resolution.CreateHabitats)
		}
	}
	ready, err = database.Queries().CheckSamplingHabitatsReady(ctx, batch)
	require.NoError(t, err)
	require.True(t, ready)
	ready, err = database.Queries().CheckSamplingHabitatsReady(ctx, other)
	require.NoError(t, err)
	require.False(t, ready, "decisions must be scoped to their batch")

	conflicting := models.SamplingHabitatResolutionInput{InputText: "Freshwater", Status: biomedb.VocabResolutionStatusRequestCreation,
		CreateHabitats: []models.HabitatInput{{Name: "cAVE", Facet: models.NewOptional("Conflicting facet")}}}
	_, err = store.ResolveHabitats(ctx, database, batch, conflicting)
	require.ErrorIs(t, err, ErrConflictingHabitatCreations)
	// Both a bad target and an unknown source roll back the previous decision.
	_, err = store.ResolveHabitats(ctx, database, batch, models.SamplingHabitatResolutionInput{
		InputText: "river", Status: biomedb.VocabResolutionStatusSelected, HabitatIDs: []uuid.UUID{uuid.New()},
	})
	require.Error(t, err)
	_, err = store.ResolveHabitats(ctx, database, other, models.SamplingHabitatResolutionInput{
		InputText: "LOTIC", Status: biomedb.VocabResolutionStatusDiscard,
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)

	rollback := errors.New("abort after linking")
	err = database.WithTx(ctx, func(tx *db.Tx) error {
		if err := store.MaterializeSamplingHabitats(ctx, tx, batch); err != nil {
			return err
		}
		return rollback
	})
	require.ErrorIs(t, err, rollback)
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM habitats WHERE label = 'Cave'").Scan(&count))
	require.Zero(t, count)
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM samplings_habitats WHERE sampling_id = $1", sampling).Scan(&count))
	require.Zero(t, count)

	require.NoError(t, database.WithTx(ctx, func(tx *db.Tx) error { return store.MaterializeSamplingHabitats(ctx, tx, batch) }))
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM samplings_habitats WHERE sampling_id = $1", sampling).Scan(&count))
	require.Equal(t, 3, count, "merge all source lines and deduplicate repeated targets")
	var caveID uuid.UUID
	var facet, description string
	require.NoError(t, pool.QueryRow(ctx, "SELECT id, facet, description FROM habitats WHERE label = 'Cave'").Scan(&caveID, &facet, &description))
	require.Equal(t, "New facet", facet)
	require.Equal(t, "Requested description", description)
	require.NoError(t, database.WithTx(ctx, func(tx *db.Tx) error { return store.MaterializeSamplingHabitats(ctx, tx, batch) }))
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM samplings_habitats WHERE sampling_id = $1", sampling).Scan(&count))
	require.Equal(t, 3, count, "linking is idempotent")

	// Another batch can request the same name; canonical metadata wins.
	decision.HabitatIDs = []uuid.UUID{caveID}
	decision.CreateHabitats = []models.HabitatInput{{Name: "cave", Facet: models.NewOptional("Other facet")}}
	_, err = store.ResolveHabitats(ctx, database, other, decision)
	require.NoError(t, err)
	require.NoError(t, database.WithTx(ctx, func(tx *db.Tx) error { return store.MaterializeSamplingHabitats(ctx, tx, other) }))
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM habitats WHERE label = 'Cave'").Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, pool.QueryRow(ctx, "SELECT facet FROM habitats WHERE id = $1", caveID).Scan(&facet))
	require.Equal(t, "New facet", facet)
	resolutions, err = store.GetHabitatsResolution(ctx, database, other)
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{caveID}, resolutions[0].HabitatIDs)
	require.Empty(t, resolutions[0].CreateHabitats)

	// Deleting a selected habitat must invalidate readiness instead of dropping it silently.
	require.NoError(t, database.Queries().DeleteHabitatByID(ctx, fresh.ID))
	ready, err = database.Queries().CheckSamplingHabitatsReady(ctx, batch)
	require.NoError(t, err)
	require.False(t, ready)
	require.NoError(t, pool.QueryRow(ctx, "UPDATE import_batches SET status = 'completed' WHERE id = $1 RETURNING id", batch).Scan(new(uuid.UUID)))
	_, err = store.ResolveHabitats(ctx, database, batch, decision)
	require.ErrorIs(t, err, ErrHabitatBatchNotEditable)
}

func TestSamplingHabitatConcurrentCreation(t *testing.T) {
	database, pool := habitatResolutionTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	store := NewSamplingStore()
	var batches, samplings []uuid.UUID
	for range 2 {
		batch, sampling := habitatBatch(t, database, []string{"source"})
		batches, samplings = append(batches, batch), append(samplings, sampling)
		_, err := store.ResolveHabitats(ctx, database, batch, models.SamplingHabitatResolutionInput{
			InputText: "source", Status: biomedb.VocabResolutionStatusRequestCreation,
			CreateHabitats: []models.HabitatInput{{Name: "Concurrent cave", Facet: models.NewOptional("New facet")}},
		})
		require.NoError(t, err)
	}
	start := make(chan struct{})
	results := make(chan error, len(batches))
	for _, batch := range batches {
		go func() {
			<-start
			results <- database.WithTx(ctx, func(tx *db.Tx) error { return store.MaterializeSamplingHabitats(ctx, tx, batch) })
		}()
	}
	close(start)
	for range batches {
		require.NoError(t, <-results)
	}
	var count int
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM habitats WHERE label = 'Concurrent cave'").Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM samplings_habitats WHERE sampling_id = ANY($1::uuid[])", samplings).Scan(&count))
	require.Equal(t, 2, count)
}

func TestSamplingHabitatReadinessRejectsIncompleteTargets(t *testing.T) {
	database, pool := habitatResolutionTestDB(t)
	ctx := context.Background()
	batch, _ := habitatBatch(t, database, []string{"source"})
	assertReady := func(expected bool) {
		t.Helper()
		ready, err := database.Queries().CheckSamplingHabitatsReady(ctx, batch)
		require.NoError(t, err)
		require.Equal(t, expected, ready)
	}
	// A LEFT JOIN with no targets produces NULL aggregates, not zeroes.
	// Neither selected nor request_creation may become ready in that situation.
	for _, status := range []biomedb.VocabResolutionStatus{biomedb.VocabResolutionStatusSelected, biomedb.VocabResolutionStatusRequestCreation} {
		_, err := database.Queries().UpdateSamplingHabitatsResolution(ctx, biomedb.UpdateSamplingHabitatsResolutionParams{
			ImportID: batch, InputText: "source", Status: status,
		})
		require.NoError(t, err)
		assertReady(false)
	}
	_, err := pool.Exec(ctx, "UPDATE sampling_habitats_resolution SET status = 'discard' WHERE import_id = $1", batch)
	require.NoError(t, err)
	assertReady(true)
	_, err = pool.Exec(ctx, "INSERT INTO sampling_habitat_targets (import_id, input_text) VALUES ($1, 'source')", batch)
	require.NoError(t, err)
	assertReady(false)
	_, err = pool.Exec(ctx, "DELETE FROM sampling_habitats_resolution WHERE import_id = $1", batch)
	require.NoError(t, err)
	assertReady(false)
	_, err = pool.Exec(ctx, "UPDATE import_samplings_occurrences SET habitats = NULL WHERE import_id = $1", batch)
	require.NoError(t, err)
	assertReady(true)
}

func TestSamplingHabitatDecisionCannotFollowMaterialization(t *testing.T) {
	database, _ := habitatResolutionTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	store := NewSamplingStore()
	batch, _ := habitatBatch(t, database, []string{"source"})
	input := models.SamplingHabitatResolutionInput{InputText: "source", Status: biomedb.VocabResolutionStatusDiscard}
	_, err := store.ResolveHabitats(ctx, database, batch, input)
	require.NoError(t, err)
	locked := make(chan struct{})
	finish := make(chan struct{})
	materialized := make(chan error, 1)
	go func() {
		materialized <- database.WithTx(ctx, func(tx *db.Tx) error {
			if _, err := tx.Queries().LockImportBatchForHabitats(ctx, batch); err != nil {
				return err
			}
			close(locked)
			select {
			case <-finish:
			case <-ctx.Done():
				return ctx.Err()
			}
			if err := store.MaterializeSamplingHabitats(ctx, tx, batch); err != nil {
				return err
			}
			rows, err := tx.Query(ctx, "UPDATE import_batches SET status = 'completed' WHERE id = $1", batch)
			if err != nil {
				return err
			}
			rows.Close()
			return rows.Err()
		})
	}()
	select {
	case <-locked:
	case err := <-materialized:
		require.NoError(t, err)
		t.Fatal("batch lock failed")
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	edited := make(chan error, 1)
	go func() { _, err := store.ResolveHabitats(ctx, database, batch, input); edited <- err }()
	close(finish)
	require.NoError(t, <-materialized)
	require.ErrorIs(t, <-edited, ErrHabitatBatchNotEditable)
}
