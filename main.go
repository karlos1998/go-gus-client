package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/orkanap/regonapi"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Nie udało się wczytać .env")
	}

	app := fiber.New()

	apiKey := os.Getenv("REGON_API_KEY")
	if apiKey == "" {
		log.Fatal("Brakuje REGON_API_KEY")
	}

	app.Get("/search", func(c *fiber.Ctx) error {
		nip := c.Query("nip")
		if nip == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Brak parametru nip"})
		}

		client := regonapi.NewClient(context.Background(), apiKey)

		if err := client.Login(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd logowania"})
		}
		defer client.Logout()

		entities, err := client.SearchByNIP(nip)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd wyszukiwania"})
		}

		return c.JSON(entities)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Fatal(app.Listen(":4300"))
}
