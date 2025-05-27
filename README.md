# Property Listing Backend

A robust backend system for managing property listings, built with **Go**, **MongoDB**, and **Redis**.

---

## Features

- **User Authentication:** Register and login with email/password.
- **Property CRUD:** Create, read, update, and delete property listings.
- **Advanced Filtering:** Search properties by 10+ attributes (state, city, area, amenities, etc.).
- **Favorites:** Users can favorite/unfavorite properties.
- **Recommendations:** Users can recommend properties to other registered users by email.
- **Caching:** Redis caching for property lists and details.
- **CSV Import:** Bulk import properties from a CSV file.
- **Secure:** All secrets managed via environment variables.

---

## Prerequisites

- Go 1.21 or higher
- MongoDB (Atlas or local)
- Redis (optional, for caching)

---

## Getting Started

### 1. **Clone the Repository**

```bash
git clone https://github.com/Ashwinitiwari18/Hypergro-Assignment.git
cd Hypergro-assignment
```

### 2. **Install Dependencies**

```bash
go mod download
```

### 3. **Set Up Environment Variables**

Create a `.env` file in the root directory:

```
MONGODB_URI=your-mongodb-connection-string
GIN_MODE=debug
# PORT=8080
# REDIS_URL=redis://localhost:6379
```

> **Note:** Never commit your `.env` file. Only share `.env.example`!

### 4. **Run the Application**

```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

---

## Live Demo

The application is currently deployed and accessible at: [https://hypergro-assignment-production.up.railway.app](https://hypergro-assignment-production.up.railway.app)

---


## API Endpoints

### **Authentication**
- `POST /api/auth/register` — Register a new user
- `POST /api/auth/login` — Login and get a JWT token

### **Properties**
- `POST /api/properties` — Create a property (auth required)
- `GET /api/properties` — List/search properties (supports advanced filtering)
- `GET /api/properties/:id` — Get property details
- `PUT /api/properties/:id` — Update property (only creator)
- `DELETE /api/properties/:id` — Delete property (only creator)

### **Favorites**
- `POST /api/favorites/:id` — Add property to favorites
- `DELETE /api/favorites/:id` — Remove property from favorites
- `GET /api/favorites` — List user's favorite properties

### **Recommendations**
- `POST /api/properties/:id/recommend` — Recommend a property to another user by email
- `GET /api/recommendations` — List received recommendations
- `PUT /api/recommendations/:id/read` — Mark recommendation as read

---

## Detailed API Documentation

### Base URL
```
https://hypergro-assignment.railway.app
```

### Authentication Endpoints

#### Register User
```
POST /api/auth/register
```
Request Body:
```json
{
    "email": "test@example.com",
    "password": "Test@123",
    "name": "Test User"
}
```

#### Login
```
POST /api/auth/login
```
Request Body:
```json
{
    "email": "test@example.com",
    "password": "Test@123"
}
```

### Property Endpoints

#### Create Property
```
POST /api/properties
```
Headers:
```
Authorization: Bearer <your_jwt_token>
```
Request Body:
```json
{
    "title": "Luxury Apartment",
    "description": "Beautiful 3BHK apartment in prime location",
    "state": "California",
    "city": "Los Angeles",
    "areaSqFt": 2000,
    "bedrooms": 3,
    "bathrooms": 2,
    "furnished": true,
    "availableFrom": "2024-03-01",
    "price": 500000,
    "listingType": "sale",
    "amenities": ["pool", "gym", "parking"],
    "tags": ["luxury", "modern"],
    "colorTheme": "light",
    "isVerified": true
}
```

#### List/Search Properties
```
GET /api/properties
```
Query Parameters (all optional):
```
?state=California
&city=Los Angeles
&areaSqFtMin=1000
&areaSqFtMax=5000
&bedrooms=3
&bathrooms=2
&furnished=true
&availableFrom=2024-03-01
&tags=luxury
&colorTheme=light
&ratingMin=4
&ratingMax=5
&isVerified=true
&listingType=sale
&amenities=pool,gym
```

### Favorites Endpoints

#### Add to Favorites
```
POST /api/favorites/:id
```
Headers:
```
Authorization: Bearer <your_jwt_token>
```

#### List Favorites
```
GET /api/favorites
```
Headers:
```
Authorization: Bearer <your_jwt_token>
```

### Recommendations Endpoints

#### Recommend Property
```
POST /api/properties/:id/recommend
```
Headers:
```
Authorization: Bearer <your_jwt_token>
```
Request Body:
```json
{
    "email": "recipient@example.com",
    "message": "Check out this amazing property!"
}
```

### Response Formats

#### Success Response
```json
{
    "data": {
        // Response data
    },
    "message": "Success message"
}
```

#### Error Response
```json
{
    "error": {
        "code": "ERROR_CODE",
        "message": "Human readable error message"
    }
}
```

### Testing Tips

1. **Authentication**
   - Always include JWT token in Authorization header for protected routes
   - Test with both valid and invalid tokens
   - Test token expiration scenarios

2. **Data Validation**
   - Test with valid and invalid data formats
   - Test required fields
   - Test boundary values
   - Test special characters

3. **Error Handling**
   - Verify error message format
   - Check error status codes
   - Validate error details

4. **Performance**
   - Test pagination with different page sizes
   - Test response times
   - Test with large datasets

5. **Security**
   - Test unauthorized access
   - Test input validation

---

## Advanced Property Filtering

You can filter properties by any combination of:
- `state`, `city`, `areaSqFtMin`, `areaSqFtMax`, `bedrooms`, `bathrooms`, `furnished`, `availableFrom`, `listedBy`, `tags`, `colorTheme`, `ratingMin`, `ratingMax`, `isVerified`, `listingType`, `amenities`, etc.

**Example:**
```
GET /api/properties?state=California&city=Los Angeles&areaSqFtMin=1000&areaSqFtMax=5000&furnished=yes&tags=luxury&isVerified=true
```

---

## CSV Import

To bulk import properties from a CSV file, temporarily uncomment the import block in `main.go` and set the path to your CSV.  
Set the environment variable:

```bash
export IMPORT_CSV=true
```

Run the app once, then comment the import block again to avoid duplicate imports.

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Create a new Pull Request

---

## Author

- [Ashwinitiwari18](https://github.com/Ashwinitiwari18)

---

**For any questions or issues, please open an issue on GitHub.** 
