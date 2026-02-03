package handlers

import (
	"context"
	"time"

	"bolt-backend/config"
	"bolt-backend/database"
	"bolt-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with first name, last name, email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing fields"
// @Failure 500 {object} map[string]interface{} "Failed to create user"
// @Router /api/users [post]
func CreateUser(c *fiber.Ctx) error {
	cfg := config.LoadConfig()
	collection := database.GetCollection(cfg.DatabaseName, "users")

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if user.FirstName == "" || user.LastName == "" || user.EmailID == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required: first_name, last_name, email_id, password",
		})
	}

	user.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users from the database
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of users with count"
// @Failure 500 {object} map[string]interface{} "Failed to fetch users"
// @Router /api/users [get]
func GetUsers(c *fiber.Ctx) error {
	cfg := config.LoadConfig()
	collection := database.GetCollection(cfg.DatabaseName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode users",
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"count": len(users),
	})
}
