package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"livescore/internal/api"
	"livescore/internal/db"
	"livescore/internal/realtime"
	"livescore/internal/scheduler"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load .env file if present (for local dev)
	_ = godotenv.Load()

	redisURL := os.Getenv("REDIS_URL")
	pgURL := os.Getenv("DATABASE_URL")

	if err := realtime.InitRedis(redisURL); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Retry logic for database connection
	var dbErr error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		dbErr = db.InitDB(pgURL)
		if dbErr == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, dbErr)
		time.Sleep(3 * time.Second)
	}

	if dbErr != nil {
		log.Fatalf("Failed to initialize database after %d attempts: %v", maxRetries, dbErr)
	}

	defer db.Pool.Close()

	app := fiber.New()

	// Add Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	// Add CORS middleware to allow admin dashboard requests
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Allow all origins in development
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders:    "Content-Range, X-Total-Count",
		AllowCredentials: true,
	}))

	api.RegisterRoutes(app)

	// Initialize scheduler
	scheduler.InitScheduler()

	// Get scheduler configuration from environment
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:3000"
	}

	cronSchedule := os.Getenv("SPORTMONKS_CRON_SCHEDULE")
	if cronSchedule == "" {
		// Default: Run every 3 months on the 1st day at midnight (0 0 0 1 */3 *)
		cronSchedule = "0 0 0 1 */3 *"
	}

	runImmediately := os.Getenv("RUN_SPORTMONKS_IMMEDIATELY") != "false"

	// Schedule the job
	err := scheduler.ScheduleFetchSportmonksJob(cronSchedule, apiBaseURL)
	if err != nil {
		log.Printf("Failed to schedule Sportmonks fetch job: %v", err)
	} else {
		log.Printf("Scheduled Sportmonks fetch job with cron schedule: %s", cronSchedule)

		// Also run the job immediately on startup if configured to do so
		if runImmediately {
			go func() {
				// Wait a short time for the server to start
				time.Sleep(5 * time.Second)
				if err := scheduler.RunFetchSportmonksNow(apiBaseURL); err != nil {
					log.Printf("Failed to run immediate Sportmonks fetch job: %v", err)
				}
			}()
		}
	}

	fmt.Println("Server running on :3000")
	log.Fatal(app.Listen(":3000"))
}
