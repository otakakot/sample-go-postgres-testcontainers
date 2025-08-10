package testx

import (
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

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
		"db":       db,
		"user":     user,
		"password": password,
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

	dsn := "postgres://" + user + ":" + password + "@" + host + ":" + port.Port() + "/" + db + "?sslmode=disable"

	return dsn
}
