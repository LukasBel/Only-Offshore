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
	CreatedAt         time.Time `json:"createdAt"`
}

func (r *Repository) CreateUser(c *fiber.Ctx) error {
	spotModel := Stat{}
	err := c.BodyParser(&spotModel)
	spotModel.CreatedAt = time.Now()

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "something went wrong"})
		return err
	}

	err = r.DB.Create(spotModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to create database entry"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "spot created successfully!"})
	return nil

}

func (r *Repository) GetStats(c *fiber.Ctx) error {
	spotModels := &[]Models.Stats{}
	err := r.DB.Find(&spotModels).Error

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "failed to get surf spots"})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "spots found successfully!", "data": spotModels})
	return nil
}

func (r *Repository) GetStatsByID(c *fiber.Ctx) error {
	id := c.Params("id")
	spotModel := &Models.Stats{}

	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "empty ID"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&spotModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to retrieve spot"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "spot found successfully", "data": spotModel})
	return nil

}

func (r *Repository) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	surfModel := &Models.Stats{}

	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "empty ID"})
		return nil
	}

	err := r.DB.Where("id = ?", id).Delete(&surfModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to delete spot"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "spot deleted successfully"})
	return nil

}

func (r *Repository) UpdateStats(c *fiber.Ctx) error {
	id := c.Params("id")
	spotModel := &Models.Stats{}
	newModel := Stat{}

	err := c.BodyParser(&newModel)
	if err != nil {
		return err
	}

	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "empty ID"})
		return nil
	}
	err = r.DB.Model(spotModel).Where("id = ?", id).Updates(newModel).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to update spot"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "spot updated successfully", "data": spotModel})
	return nil
}

func (r *Repository) Progress(c *fiber.Ctx) error {
	before := &Models.Stats{}
	after := &Models.Stats{}

	err := r.DB.First(&before).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to find user"})
		return err
	}

	err = r.DB.Last(&after).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "failed to find user"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "stats fetched succesfully", "data": after})
	return nil

}

func (r *Repository) SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/stats", r.GetStats)
	api.Get("/stats/:id", r.GetStatsByID)
	api.Post("/create", r.CreateUser)
	api.Put("/update/:id", r.UpdateStats)
	api.Delete("/delete/:id", r.DeleteUser)
	api.Get("/progress", r.Progress)
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

	/*
		emails := os.Getenv("TO")
		emailAddresses := strings.Split(emails, ",")

		Handlers.SendMail(emailAddresses)
	*/

	app := fiber.New()
	r.SetUpRoutes(app)
	app.Listen(":8080")
}
