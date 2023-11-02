package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	internalEs "qualification/internal/elastic"

	"github.com/golang/mock/gomock"
	"github.com/olivere/elastic/v7"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	postgresConn *pgxpool.Pool
	esConn       *elastic.Client
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	poolPostgres, resourcePostgres := startPostgres(ctx)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := poolPostgres.Purge(resourcePostgres); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func startPostgres(ctx context.Context) (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	if err = pool.Client.Ping(); err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=disable", hostAndPort)

	log.Printf("Connecting to database on url: %v\n", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		postgresConn, err = pgxpool.Connect(ctx, databaseUrl)
		if err != nil {
			return err
		}
		return postgresConn.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	cmd := exec.Command("flyway",
		"-user=postgres",
		"-password=postgres",
		"-locations=filesystem:../../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%v/postgres", hostAndPort),
		"migrate")
	if err = cmd.Run(); err != nil {
		log.Fatalf("There are errors in migrations: %s", err)
	}
	return pool, resource
}

func TestApp_CreateProduct(t *testing.T) {
	type fields struct {
		Router *mux.Router
		DB     *pgxpool.Pool
		ES     *internalEs.MockES
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields func(f *fields)
		args   func(a *args)
	}{
		{
			name: "test 1",
			fields: func(f *fields) {
				f.DB = postgresConn
			},
			args: func(a *args) {
				a.w = httptest.NewRecorder()
				a.r = httptest.NewRequest(http.MethodPost, "/product",
					strings.NewReader(`{"id": 1,"name": qw}`),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				ES:     internalEs.NewMockES(ctrl),
				Router: mux.NewRouter(),
			}
			if tt.fields != nil {
				tt.fields(&f)
			}

			a := args{}
			if tt.args != nil {
				tt.args(&a)
			}

			app := &App{
				Router: f.Router,
				DB:     f.DB,
				ES:     f.ES,
			}
			app.CreateProduct(a.w, a.r)
		})
	}
}
