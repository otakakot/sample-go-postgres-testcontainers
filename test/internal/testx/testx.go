package testx

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupContainer(
	t *testing.T,
) string {
	t.Helper()

	user := "test"
	password := "test"
	db := "test"

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	schemaPath := filepath.Join(filepath.Dir(pwd), "schema")

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Fatalf("Schema directory does not exist: %s", schemaPath)
	}

	sqlFiles, err := filepath.Glob(filepath.Join(schemaPath, "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	if len(sqlFiles) == 0 {
		t.Fatalf("No SQL files found in schema directory: %s", schemaPath)
	}

	containerFiles := make([]testcontainers.ContainerFile, 0, len(sqlFiles))

	for _, sqlFile := range sqlFiles {
		fileName := filepath.Base(sqlFile)
		containerFiles = append(containerFiles, testcontainers.ContainerFile{
			HostFilePath:      sqlFile,
			ContainerFilePath: "/docker-entrypoint-initdb.d/" + fileName,
			FileMode:          0644,
		})
	}

	container, err := testcontainers.GenericContainer(
		t.Context(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "postgres:17-alpine",
				ExposedPorts: []string{"5432/tcp"},
				Env: map[string]string{
					"TZ":                        "UTC",
					"LANG":                      "ja_JP.UTF-8",
					"POSTGRES_DB":               db,
					"POSTGRES_USER":             user,
					"POSTGRES_PASSWORD":         password,
					"POSTGRES_INITDB_ARGS":      "--encoding=UTF-8",
					"POSTGRES_HOST_AUTH_METHOD": "trust",
				},
				Cmd: []string{"postgres", "-c", "log_statement=all"},
				WaitingFor: wait.ForAll(
					wait.ForListeningPort("5432/tcp"),
					wait.ForExec([]string{"pg_isready", "-U", user, "-d", db}).
						WithPollInterval(1*time.Second).
						WithExitCodeMatcher(func(exitCode int) bool {
							return exitCode == 0
						}).
						WithStartupTimeout(30*time.Second),
				),
				Files: containerFiles,
			},
			Started: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	testcontainers.CleanupContainer(t, container)

	host, err := container.Host(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(t.Context(), "5432/tcp")
	if err != nil {
		t.Fatal(err)
	}

	return "postgres://" + user + ":" + password + "@" + host + ":" + port.Port() + "/" + db + "?sslmode=disable"
}

func SetupPostgres(
	t *testing.T,
) string {
	t.Helper()

	user := "test"
	password := "test"
	db := "test"

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	schemaPath := filepath.Join(filepath.Dir(pwd), "schema")

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Fatalf("Schema directory does not exist: %s", schemaPath)
	}

	sqlFiles, err := filepath.Glob(filepath.Join(schemaPath, "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	if len(sqlFiles) == 0 {
		t.Fatalf("No SQL files found in schema directory: %s", schemaPath)
	}

	container, err := postgres.Run(
		t.Context(),
		"postgres:17-alpine",
		postgres.WithDatabase(db),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.WithInitScripts(sqlFiles...),
		testcontainers.WithEnv(map[string]string{
			"TZ":                        "UTC",
			"LANG":                      "ja_JP.UTF-8",
			"POSTGRES_INITDB_ARGS":      "--encoding=UTF-8",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForExec([]string{"pg_isready", "-U", user, "-d", db}).
					WithPollInterval(1*time.Second).
					WithExitCodeMatcher(func(exitCode int) bool {
						return exitCode == 0
					}).
					WithStartupTimeout(30*time.Second),
			),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return dsn
}

func SetupCompose(
	t *testing.T,
) string {
	user := "test"
	password := "test"
	db := "test"

	pwd, _ := os.Getwd()

	stack, err := compose.NewDockerComposeWith(
		compose.WithStackFiles(pwd + "/compose.yaml"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := stack.WithEnv(map[string]string{
		"DB":       db,
		"USER":     user,
		"PASSWORD": password,
	}).Up(
		t.Context(),
		compose.RunServices("postgres"),
		compose.Wait(true),
	); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = stack.Down(t.Context())
	})

	container, err := stack.ServiceContainer(t.Context(), "postgres")
	if err != nil {
		t.Fatal(err)
	}

	host, err := container.Host(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(t.Context(), "5432/tcp")
	if err != nil {
		t.Fatal(err)
	}

	return "postgres://" + user + ":" + password + "@" + host + ":" + port.Port() + "/" + db + "?sslmode=disable"
}
