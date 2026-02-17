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
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"


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
	// Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": "Failed to hash password",
    })
}

    user.Password = string(hashedPassword)

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
// Login godoc
// @Summary Login user
// @Description Login using email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Router /api/users/login [post]
func LoginUser(c *fiber.Ctx) error {
    cfg := config.LoadConfig()
    collection := database.GetCollection(cfg.DatabaseName, "users")

    var req models.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var user models.User
    err := collection.FindOne(ctx, bson.M{"email_id": req.EmailID}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    // Compare password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid password",
        })
    }

    // Generate JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID.Hex(),
        "email":   user.EmailID,
        "exp":     time.Now().Add(time.Hour * 1).Unix(),
    })

    tokenString, err := token.SignedString([]byte("your-secret-key"))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not generate token",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Login successful",
        "token":   tokenString,
    })
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Generate new token using existing token
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Router /api/users/refresh [post]
func RefreshToken(c *fiber.Ctx) error {
    tokenString := c.Get("Authorization")

    if tokenString == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Missing token",
        })
    }

    tokenString = tokenString[len("Bearer "):]

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    claims := token.Claims.(jwt.MapClaims)

    newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": claims["user_id"],
        "email":   claims["email"],
        "exp":     time.Now().Add(time.Hour * 1).Unix(),
    })

    newTokenString, _ := newToken.SignedString([]byte("your-secret-key"))

    return c.JSON(fiber.Map{
        "token": newTokenString,
    })
}
