package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"qualification/internal/usecase/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var postgresConn *pgxpool.Pool

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

func TestCreateProduct(t *testing.T) {
	type args struct {
		ctx       context.Context
		db        *pgxpool.Pool
		blProduct model.Product
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				ctx: context.Background(),
				db:  postgresConn,
				blProduct: model.Product{
					ID:   1,
					Name: "asd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateProduct(tt.args.ctx, tt.args.db, tt.args.blProduct); (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *pgxpool.Pool
		id  int
	}
	tests := []struct {
		name    string
		args    func(a *args)
		want    *model.Product
		wantErr bool
	}{
		{
			name: "test 1",
			args: func(a *args) {
				a.ctx = context.Background()
				a.db = postgresConn
				a.id = 1

				blProduct := model.Product{
					ID:   a.id,
					Name: "asd",
				}

				if err := CreateProduct(a.ctx, a.db, blProduct); err != nil {
					t.Errorf("CreateProduct() error = %v", err)
				}
			},
			want: &model.Product{ID: 1, Name: "asd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := args{}
			if tt.args != nil {
				tt.args(&a)
			}

			got, err := GetProduct(a.ctx, a.db, a.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProduct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
