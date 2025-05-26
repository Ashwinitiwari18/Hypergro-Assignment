# Property Listing Backend

A robust backend system for managing property listings built with Go, MongoDB, and Redis.

## Features

- User authentication (register/login)
- CRUD operations for properties
- Advanced property filtering and search
- Property favorites system
- Property recommendations between users
- Redis caching for improved performance
- CSV data import functionality

## Prerequisites

- Go 1.21 or higher
- MongoDB
- Redis

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd property-listing
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```
Edit the `.env` file with your configuration.

4. Run the application:
```bash
go run main.go
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Properties
- `POST /api/properties` - Create a new property
- `GET /api/properties` - List properties with filtering
- `GET /api/properties/:id` - Get property details
- `PUT /api/properties/:id` - Update property
- `DELETE /api/properties/:id` - Delete property

### Favorites
- `POST /api/favorites/:id` - Add property to favorites
- `DELETE /api/favorites/:id` - Remove property from favorites
- `GET /api/favorites` - List user's favorite properties

### Recommendations
- `POST /api/properties/:id/recommend` - Recommend a property to another user
- `GET /api/recommendations` - List received recommendations
- `PUT /api/recommendations/:id/read` - Mark recommendation as read

## CSV Import

To import the initial property data from CSV, set the environment variable:
```bash
export IMPORT_CSV=true
```

## Deployment

The application can be deployed to any platform that supports Go applications. For example:

1. Render
2. Vercel
3. Heroku

Make sure to set up the required environment variables in your deployment platform.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 

MONGODB_URI=your-mongodb-connection-string 