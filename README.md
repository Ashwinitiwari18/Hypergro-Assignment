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
