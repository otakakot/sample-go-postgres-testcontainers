package test_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/sample-go-postgres-testcontainers/pkg/schema"
	"github.com/otakakot/sample-go-postgres-testcontainers/test/internal/testx"
)

func TestSetupContainer(t *testing.T) {
	dsn := testx.SetupContainer(t)

	t.Logf("Postgres DSN: %s", dsn)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	t.Cleanup(pool.Close)

	sample, err := schema.New(pool).InsertSample(t.Context(), "sample")
	if err != nil {
		t.Fatalf("failed to insert sample: %v", err)
	}

	t.Logf("Inserted Sample: %+v", sample)

	samples, err := schema.New(pool).SelectSample(t.Context())
	if err != nil {
		t.Fatalf("failed to select samples: %v", err)
	}

	for _, sample := range samples {
		t.Logf("Sample: %+v", sample)
	}
}

func TestSetupPostgres(t *testing.T) {
	dsn := testx.SetupPostgres(t)

	t.Logf("Postgres DSN: %s", dsn)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	t.Cleanup(pool.Close)

	sample, err := schema.New(pool).InsertSample(t.Context(), "sample")
	if err != nil {
		t.Fatalf("failed to insert sample: %v", err)
	}

	t.Logf("Inserted Sample: %+v", sample)

	samples, err := schema.New(pool).SelectSample(t.Context())
	if err != nil {
		t.Fatalf("failed to select samples: %v", err)
	}

	for _, sample := range samples {
		t.Logf("Sample: %+v", sample)
	}
}

func TestSetupCompose(t *testing.T) {
	dsn := testx.SetupCompose(t)

	t.Logf("Postgres DSN: %s", dsn)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	t.Cleanup(pool.Close)

	sample, err := schema.New(pool).InsertSample(t.Context(), "sample")
	if err != nil {
		t.Fatalf("failed to insert sample: %v", err)
	}

	t.Logf("Inserted Sample: %+v", sample)

	samples, err := schema.New(pool).SelectSample(t.Context())
	if err != nil {
		t.Fatalf("failed to select samples: %v", err)
	}

	for _, sample := range samples {
		t.Logf("Sample: %+v", sample)
	}
}
