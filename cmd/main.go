// ND redis part
package main

import (
	"os"
	"time"

	"github.com/goccy/go-json"

	"github.com/aj-2000/shc-backend/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/storage/redis/v3"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	println("Starting server")
	println(os.Getenv("PORT"))

	// ND what is the meaning of starting a new service does it start a new service for each user?
	service := services.NewAppService()

	//iska mtlb h ki fiber k config me JSONEncoder aur JSONDecoder me json.Marshal aur json.Unmarshal use honge???
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// ND what is this line doing? ✅
	app.Use(logger.New())

	// i don't know what is redis
	// TODO: check if it interferes with other keys
	storage := redis.New((redis.Config{
		URL: os.Getenv("REDIS_URL"),
	}))

	// ND not understood below function ✅
	app.Use(limiter.New(limiter.Config{
		Max:        15,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c fiber.Ctx) string {
			return "rate_limit_" + c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "You have exceeded your rate limit",
			})
		},
		Storage: storage,
	}))

	// TODO: Verify if it helps
	// what is this? ✅
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	// not understood below function ✅
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost, http://localhost:3000, https://shc-app.vercel.app, https://shc.ajaysharma.dev, https://ajaysharma.dev",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	setupRoutes(app, service)

	// Prefork creates multiple independent processes, each handling a single connection, while the child model uses a single process to handle multiple connections concurrently.
	// not understood below function why we need it?
	if !fiber.IsChild() {
		runCronJobs(service)
	}

	app.Listen(":" + os.Getenv("PORT"))
}
