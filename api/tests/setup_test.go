package tests

import (
	"context"
	"finance-control-api/database"
	"finance-control-api/models"
	"finance-control-api/routes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

const (
	testDBUser = "postgres"
	testDBPass = "postgres"
	testDBName = "finance_test"
	testSchema = "../../postgres/init.sql"
)

// testUser holds the credentials of a user created at runtime by TestMain.
var (
	testUserID       int
	testUserLogin    = "testuser"
	testUserPassword = "testpassword"
)

func TestMain(m *testing.M) {
	code, err := runTests(m)
	if err != nil {
		log.Fatal("Error setting up test database: " + err.Error())
	}
	os.Exit(code)
}

func runTests(m *testing.M) (int, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	ctx := context.Background()

	container, mapping, err := startTestDatabase(ctx)
	if err != nil {
		return 1, err
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Printf("Error terminating test database container: %v", err)
		}
	}()

	if err := applySchema(ctx, mapping); err != nil {
		return 1, err
	}

	if err := database.InitPool(); err != nil {
		return 1, fmt.Errorf("initializing database pool: %w", err)
	}

	if err := database.PingDatabase(); err != nil {
		return 1, fmt.Errorf("pinging database: %w", err)
	}

	if err := createTestUser(); err != nil {
		return 1, err
	}

	return m.Run(), nil
}

// createTestUser inserts a user used as the authenticated subject of the API
// tests, since postgres/init.sql has no seed users.
func createTestUser() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(testUserPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing test user password: %w", err)
	}

	user, err := database.CreateUser("Test User", testUserLogin, string(hash))
	if err != nil {
		return fmt.Errorf("creating test user: %w", err)
	}

	testUserID = user.ID
	return nil
}

type dbMapping struct {
	host string
	port string
}

func startTestDatabase(ctx context.Context) (testcontainers.Container, dbMapping, error) {
	pg, err := postgres.Run(ctx,
		"postgres:18-trixie",
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testDBUser),
		postgres.WithPassword(testDBPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, dbMapping{}, fmt.Errorf("starting postgres container: %w", err)
	}

	host, err := pg.Host(ctx)
	if err != nil {
		return nil, dbMapping{}, fmt.Errorf("getting container host: %w", err)
	}
	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		return nil, dbMapping{}, fmt.Errorf("getting container port: %w", err)
	}

	if err := os.Setenv("DB_HOST", host); err != nil {
		return nil, dbMapping{}, fmt.Errorf("setting DB_HOST: %w", err)
	}
	if err := os.Setenv("DB_PORT", port.Port()); err != nil {
		return nil, dbMapping{}, fmt.Errorf("setting DB_PORT: %w", err)
	}
	if err := os.Setenv("DB_USER", testDBUser); err != nil {
		return nil, dbMapping{}, fmt.Errorf("setting DB_USER: %w", err)
	}
	if err := os.Setenv("DB_PASSWORD", testDBPass); err != nil {
		return nil, dbMapping{}, fmt.Errorf("setting DB_PASSWORD: %w", err)
	}
	if err := os.Setenv("DB_NAME", testDBName); err != nil {
		return nil, dbMapping{}, fmt.Errorf("setting DB_NAME: %w", err)
	}

	return pg, dbMapping{host: host, port: port.Port()}, nil
}

func applySchema(ctx context.Context, mapping dbMapping) error {
	sqlBytes, err := os.ReadFile(filepath.Clean(testSchema))
	if err != nil {
		return fmt.Errorf("reading schema file %s: %w", testSchema, err)
	}

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		testDBUser, testDBPass, mapping.host, mapping.port, testDBName)

	db, err := database.OpenExternal(uri)
	if err != nil {
		return fmt.Errorf("opening test database connection: %w", err)
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("applying schema: %w", err)
	}

	return nil
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Request().Header.Set(echo.HeaderAuthorization, "Bearer "+testToken(testUserID))
			return next(c)
		}
	})
	routes.BindRoutes(e)
	return e
}

// testToken returns a signed JWT carrying the user_id claim used by the
// JWTMiddleware. The token is only valid if a user with that ID exists, so
// tests must use testUserID (created in TestMain).
func testToken(userID int) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"sub":     testUserLogin,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func createTestArea(name string) (models.Area, error) {
	return database.CreateArea(testUserID, name)
}

func createTestBudget(year, month, areaID, amount int) (models.Budget, error) {
	props := database.BudgetProps{
		Year:   &year,
		Month:  &month,
		AreaId: &areaID,
		Amount: &amount,
	}
	return database.CreateBudget(testUserID, props)
}

func createTestCategory(name string) (models.Category, error) {
	return database.CreateCategory(testUserID, name)
}

func createTestTransaction(date string, amount, categoryID, areaID int, transType string) (models.Transaction, error) {
	props := database.TransactionProps{
		Date:       &date,
		Amount:     &amount,
		CategoryID: &categoryID,
		AreaID:     &areaID,
		Type:       &transType,
	}
	return database.CreateTransaction(testUserID, props)
}
