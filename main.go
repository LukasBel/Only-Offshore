package main

import (
	"github.com/LukasBel/Only-Offshore.git/Models"
	"github.com/LukasBel/Only-Offshore.git/Storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type Repository struct {
	DB *gorm.DB
}

type Stat struct {
	User              string    `json:"user"`
	FlatBench         int       `json:"flatBench"`
	InclineBench      int       `json:"inclineBench"`
	Squat             int       `json:"squat"`
	PullUps           int       `json:"pullUps"`
	WeightedPullUpMax int       `json:"weightedPullUpMax"`
	BodyWeight        int       `json:"bodyWeight"`
	CreationDate      time.Time `json:"creationDate"`
}

func (r *Repository) GetStats(c *fiber.Ctx) error {
	statsModel := &[]Models.Stats{}

	err := r.DB.Find(&statsModel).Error
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "failed to fetch stats"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "stats fetched succesfully", "data": statsModel})
	return nil
}

func (r *Repository) GetStatsByID(c *fiber.Ctx) error {
	statsModel := &Models.Stats{}
	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "empty ID"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&statsModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to fetch stats"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "stats fetched succesfully", "data": statsModel})
	return nil
}

func (r *Repository) CreateUser(c *fiber.Ctx) error {
	statsModel := Stat{}
	err := c.BodyParser(&statsModel)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "failed to parse body"})
		return err
	}

	statsModel.CreationDate = time.Now()

	err = r.DB.Create(statsModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to create entry"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "user created succesfully"})
	return nil
}

func (r *Repository) UpdateStats(c *fiber.Ctx) error {
	id := c.Params("id")
	statsModel := &Models.Stats{}
	newModel := Stat{}

	err := c.BodyParser(&newModel)
	if err != nil {
		return err
	}

	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "empty ID"})
		return nil
	}

	err = r.DB.Model(statsModel).Where("id = ?").Updates(newModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to update stats"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "stats updated succesfully"})
	return nil
}

func (r *Repository) DeleteUser(c *fiber.Ctx) error {
	statsModel := &Models.Stats{}
	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "empty ID"})
		return nil
	}

	err := r.DB.Where("id = ?").Delete(&statsModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to delete user"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "user deleted succesfully"})
	return nil
}

func (r *Repository) SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/stats", r.GetStats)
	api.Get("/stats/:id", r.GetStatsByID)
	api.Post("/create", r.CreateUser)
	api.Put("/update/:id", r.UpdateStats)
	api.Delete("/delete/:id", r.DeleteUser)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load environment variables")
	}

	config := &Storage.Config{
		Host:    os.Getenv("DB_HOST"),
		Port:    os.Getenv("DB_PORT"),
		User:    os.Getenv("DB_USER"),
		Pass:    os.Getenv("DB_PASS"),
		DBName:  os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	db, err := Storage.NewConnection(config)
	if err != nil {
		log.Fatal("Configuration File Issue")
	}

	err = Models.MigrateStats(db)
	if err != nil {
		log.Fatal("Failed to migrate stats")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetUpRoutes(app)
	app.Listen(":8080")

}
