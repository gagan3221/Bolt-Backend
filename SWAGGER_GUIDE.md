# Swagger Integration Guide

## Overview

Your Bolt Backend API now has complete Swagger/OpenAPI documentation integrated. This provides an interactive web interface to explore and test your APIs.

## What's Been Added

### 1. Dependencies

- `github.com/swaggo/swag` - Swagger generator
- `github.com/swaggo/fiber-swagger` - Fiber middleware for Swagger
- `github.com/swaggo/files` - Swagger static files

### 2. Documentation Annotations

API handlers now include Swagger annotations:

```go
// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with first name, last name, email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Router /api/users [post]
```

### 3. Swagger Route

Access the interactive documentation at: **http://localhost:3000/swagger/**

### 4. Auto-Generated Files

The `docs/` folder contains:

- `docs.go` - Go code for Swagger
- `swagger.json` - JSON specification
- `swagger.yaml` - YAML specification

## How to Use Swagger UI

1. **Start your server**:

   ```bash
   go run main.go
   ```

2. **Open Swagger UI**:
   Navigate to http://localhost:3000/swagger/ in your browser

3. **Try the APIs**:
   - Click on any endpoint (e.g., POST /api/users)
   - Click "Try it out"
   - Fill in the request body
   - Click "Execute"
   - View the response

## Example: Testing Create User API

1. Go to http://localhost:3000/swagger/
2. Find **POST /api/users**
3. Click "Try it out"
4. Replace the example request body:
   ```json
   {
     "first_name": "John",
     "last_name": "Doe",
     "email_id": "john@example.com",
     "password": "securePass123"
   }
   ```
5. Click "Execute"
6. See the response with the created user

## Updating Documentation

When you add new endpoints or modify existing ones:

```bash
# Method 1: Use the helper script
./generate-swagger.sh

# Method 2: Run swag init directly
~/go/bin/swag init

# Method 3: If swag is in your PATH
swag init
```

## Adding Documentation to New Endpoints

When creating new API handlers, add annotations:

```go
// HandlerName godoc
// @Summary Brief description
// @Description Detailed description
// @Tags TagName
// @Accept json
// @Produce json
// @Param paramName paramType dataType required "description"
// @Success 200 {object} responseType "Success description"
// @Failure 400 {object} errorType "Error description"
// @Router /api/path [method]
func HandlerName(c *fiber.Ctx) error {
    // handler code
}
```

## Swagger Annotation Fields

- `@Summary` - Short description (appears in list)
- `@Description` - Detailed description
- `@Tags` - Group endpoints together
- `@Accept` - Content types the API accepts
- `@Produce` - Content types the API produces
- `@Param` - Parameters (query, path, body, header)
- `@Success` - Successful response
- `@Failure` - Error responses
- `@Router` - API path and HTTP method

## Main.go Configuration

The main configuration in `main.go`:

```go
// @title Bolt Backend API
// @version 1.0
// @description REST API for Bolt Backend with MongoDB
// @host localhost:3000
// @BasePath /
// @schemes http https
```

## Benefits

✅ Interactive API documentation
✅ No need for external tools like Postman for testing
✅ Auto-generated from code comments
✅ Always in sync with your code
✅ Shareable with team members
✅ OpenAPI 2.0 compliant

## Troubleshooting

### Swagger page shows "Failed to load API definition"

- Run `~/go/bin/swag init` to regenerate docs
- Make sure docs/ folder exists
- Restart the server

### Changes not reflected in Swagger UI

- Regenerate docs: `./generate-swagger.sh`
- Hard refresh browser (Ctrl+Shift+R or Cmd+Shift+R)
- Clear browser cache

### swag command not found

- Install: `go install github.com/swaggo/swag/cmd/swag@latest`
- Or use full path: `~/go/bin/swag init`

## Resources

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Fiber Swagger](https://github.com/swaggo/fiber-swagger)
