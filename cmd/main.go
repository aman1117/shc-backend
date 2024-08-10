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
	service := services.NewAppService()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(logger.New())

	// for rate limiting
	// TODO: check if it interferes with other keys
	storage := redis.New((redis.Config{
		URL: os.Getenv("REDIS_URL"),
	}))

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
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost, http://localhost:3000, https://shc-app.vercel.app, https://shc.ajaysharma.dev, https://ajaysharma.dev",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	setupRoutes(app, service)

	if !fiber.IsChild() {
		runCronJobs(service)
	}

	app.Listen(":" + os.Getenv("PORT"))
}
