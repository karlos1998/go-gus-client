package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken != "" {
		app.Use(func(c *fiber.Ctx) error {
			auth := c.Query("auth")

			if auth == "" {
				authHeader := c.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					auth = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if auth != accessToken {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Nieprawidłowy klucz autoryzacyjny. Upewnij się czy w parametrzy query podałeś auth taki jak ustawiłeś w .env jako ACCESS_TOKEN. zapytanie powinno wyglądać tak: http://localhost:4300/search/?nip=1234567890&auth=sekretnytoken lub auth możesz podać jako header Authorization: Bearer sekretnytoken Jeśli kontener jest używany tylko lokalnie możesz też pozbyć się tego tokenu.",
				})
			}

			return c.Next()
		})
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

	app.Get("/legal-pkd-list", func(c *fiber.Ctx) error {
		regon := c.Query("regon")
		if regon == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Brak parametru regon"})
		}

		client := regonapi.NewClient(context.Background(), apiKey)

		if err := client.Login(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd logowania"})
		}
		defer client.Logout()

		entities, err := client.LegalPersonPKDList(regon)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd wyszukiwania"})
		}

		return c.JSON(entities)
	})

	app.Get("/natural-pkd-list", func(c *fiber.Ctx) error {
		regon := c.Query("regon")
		if regon == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Brak parametru regon"})
		}

		client := regonapi.NewClient(context.Background(), apiKey)

		if err := client.Login(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd logowania"})
		}
		defer client.Logout()

		entities, err := client.NaturalPersonPKDList(regon)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd wyszukiwania"})
		}

		return c.JSON(entities)
	})

	app.Get("/details", func(c *fiber.Ctx) error {
		regon := c.Query("regon")
		if regon == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Brak parametru regon"})
		}

		client := regonapi.NewClient(context.Background(), apiKey)

		if err := client.Login(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd logowania"})
		}
		defer client.Logout()

		entities, err := client.NaturalPersonDetails(regon)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Błąd wyszukiwania"})
		}

		return c.JSON(entities)
	})

	log.Fatal(app.Listen(":4300"))
}
