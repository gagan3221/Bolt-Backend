# Bolt Backend

A modern REST API built with Go, Fiber, and MongoDB.

## Features

- 🚀 Fast and lightweight using Fiber framework
- 🗄️ MongoDB integration with official driver
- 🔧 Environment-based configuration
- 🛣️ RESTful API structure
- 🔄 Graceful shutdown handling
- 📝 Logger middleware
- 🌐 CORS support
- 📚 Swagger/OpenAPI documentation

## Project Structure

```
bolt-backend/
├── config/          # Configuration management
├── database/        # Database connection and utilities
├── docs/            # Swagger documentation (auto-generated)
├── handlers/        # Request handlers
├── models/          # Data models
├── routes/          # Route definitions
├── main.go          # Application entry point
├── .env.example     # Environment variables template
└── .gitignore       # Git ignore rules
```

## Prerequisites

- Go 1.25.6 or higher
- MongoDB (local or cloud instance)

## Installation

1. Clone the repository:

```bash
git clone <your-repo-url>
cd bolt-backend
```

2. Install dependencies:

```bash
go mod download
```

3. Create your environment file:

```bash
cp .env.example .env
```

4. Update the `.env` file with your MongoDB connection string:

```env
MONGO_URI=mongodb://localhost:27017
DATABASE_NAME=bolt_db
PORT=3000
```

## Running the Application

### Development Mode

```bash
go run main.go
```

### Build and Run

```bash
go build -o bolt-backend
./bolt-backend
```

The server will start on `http://localhost:3000` (or the port specified in your `.env` file).

## Swagger API Documentation

Once the server is running, you can access the interactive Swagger UI at:

**http://localhost:3000/swagger/**

The Swagger UI provides:

- Interactive API documentation
- Try-it-out functionality for all endpoints
- Request/response examples
- Schema definitions

### Regenerating Swagger Docs

After modifying API handlers or adding new endpoints, regenerate the Swagger documentation:

```bash
# Option 1: Use the helper script
./generate-swagger.sh

# Option 2: Run swag directly
~/go/bin/swag init
```

## API Endpoints

### Health Check

- `GET /` - Welcome message
- `GET /health` - Health check endpoint

### Users API

- `POST /api/users` - Create a new user
- `GET /api/users` - Get all users

## API Examples

### Create a User

```bash
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email_id": "john@example.com",
    "password": "securePassword123"
  }'
```

### Get All Users

```bash
curl http://localhost:3000/api/users
```

## MongoDB Setup

### Local MongoDB

1. Install MongoDB: https://www.mongodb.com/docs/manual/installation/
2. Start MongoDB service:

   ```bash
   # macOS
   brew services start mongodb-community

   # Linux
   sudo systemctl start mongod
   ```

### MongoDB Atlas (Cloud)

1. Create a free cluster at https://www.mongodb.com/cloud/atlas
2. Get your connection string
3. Update `MONGO_URI` in `.env` with your Atlas connection string

## Environment Variables

| Variable        | Description               | Default                     |
| --------------- | ------------------------- | --------------------------- |
| `MONGO_URI`     | MongoDB connection string | `mongodb://localhost:27017` |
| `DATABASE_NAME` | Database name             | `bolt_db`                   |
| `PORT`          | Server port               | `3000`                      |

## Development

### Adding New Models

1. Create a new file in `models/` directory
2. Define your struct with BSON tags for MongoDB

### Adding New Routes

1. Create a handler in `handlers/` directory
2. Add the route in `routes/routes.go`

## Testing the Connection

Once the server is running, visit:

- http://localhost:3000 - Should show welcome message
- http://localhost:3000/health - Should show health status with database connection

## Troubleshooting

### MongoDB Connection Issues

- Ensure MongoDB is running
- Check your connection string in `.env`
- Verify network connectivity if using MongoDB Atlas

### Port Already in Use

- Change the `PORT` in your `.env` file
- Or kill the process using the port: `lsof -ti:3000 | xargs kill`

## License

MIT

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.
