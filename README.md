# Golang API Documentation

## Project Overview

This project is a RESTful API built with Go (Golang) that provides user authentication and watchlist management functionality. It uses MongoDB as the database and JWT for authentication.

## Folder Structure

The project follows a standard Go project structure:

```
golang-api-tutorial/
├── config/             # Database configuration
│   └── db.go
├── controller/         # HTTP request handlers
│   ├── user.go
│   └── watchlist.go
├── middleware/         # HTTP middleware
│   └── authMiddleware.go
├── model/              # Data models
│   ├── user.go
│   └── watchlist.go
├── routes/             # API route definitions
│   ├── index.go
│   ├── user.go
│   └── watchlist.go
├── services/           # Business logic services
│   └── jwt.go
├── .gitignore
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── main.go             # Application entry point
```

## Go Concepts and Packages Used

### Core Go Concepts

1. **Structs and Methods**: Used to define data models and their behavior.
2. **Interfaces**: Used for abstraction and polymorphism.
3. **Goroutines**: Lightweight threads for concurrent execution.
4. **Channels**: Used for communication between goroutines.
5. **Context**: Used for carrying deadlines, cancellation signals, and request-scoped values.
6. **Error Handling**: Go's approach to error handling with explicit error returns.
7. **Defer**: Used to ensure resources are properly released.

### External Packages

1. **github.com/gorilla/mux**: A powerful HTTP router and URL matcher for building Go web servers.
2. **go.mongodb.org/mongo-driver**: Official MongoDB driver for Go.
3. **github.com/golang-jwt/jwt/v5**: Implementation of JSON Web Tokens (JWT).
4. **github.com/joho/godotenv**: Loads environment variables from a .env file.
5. **golang.org/x/crypto/bcrypt**: Used for password hashing.

## API Endpoints

### Authentication Endpoints

#### 1. User Signup
- **URL**: `/signup`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **Response**:
  ```json
  {
    "message": "user created successfully",
    "user": {
      "id": "60d21b4667d0d8992e610c85",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
  ```
- **Status Codes**:
  - 200: Success
  - 400: Invalid JSON
  - 409: User already exists
  - 500: Server error

#### 2. User Login
- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- **Response**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "User authentication successful"
  }
  ```
- **Status Codes**:
  - 200: Success
  - 400: Invalid JSON
  - 404: User not found
  - 403: Incorrect password
  - 502: Token creation error

### Watchlist Endpoints (Protected - Require Authentication)

#### 1. Add Movie to Watchlist
- **URL**: `/watchlist`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "movie_name": "The Matrix",
    "watched": false
  }
  ```
- **Response**:
  ```json
  {
    "message": "Watchlist added successfully",
    "watchlist": {
      "id": "60d21b4667d0d8992e610c85",
      "movieName": "The Matrix",
      "watched": false
    }
  }
  ```
- **Status Codes**:
  - 200: Success
  - 400: Invalid JSON or user ID
  - 409: Movie already in watchlist
  - 500: Server error

#### 2. Get All Watchlist Items
- **URL**: `/watchlist`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  ```json
  {
    "message": "Watchlists fetched successfully",
    "watchlists": [
      {
        "id": "60d21b4667d0d8992e610c85",
        "user_id": "60d21b4667d0d8992e610c85",
        "movie_name": "The Matrix",
        "watched": false
      },
      {
        "id": "60d21b4667d0d8992e610c86",
        "user_id": "60d21b4667d0d8992e610c85",
        "movie_name": "Inception",
        "watched": true
      }
    ]
  }
  ```
- **Status Codes**:
  - 200: Success
  - 400: Invalid user ID
  - 500: Server error

#### 3. Update Watchlist Item
- **URL**: `/watchlist/{id}`
- **Method**: `PATCH`
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "watched": true
  }
  ```
- **Response**:
  ```json
  {
    "message": "Watchlist updated successfully"
  }
  ```
- **Status Codes**:
  - 200: Success
  - 400: Invalid JSON
  - 404: Watchlist not found
  - 500: Server error

#### 4. Delete Watchlist Item
- **URL**: `/watchlist/{id}`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  ```json
  {
    "message": "Watchlist deleted successfully"
  }
  ```
- **Status Codes**:
  - 200: Success
  - 404: Watchlist not found
  - 500: Server error

## Important Code Explanations

### 1. Application Entry Point (main.go)

```go
func main() {
    fmt.Println("Welcome to mongodb api")
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    config.ConnectDB()
    router := mux.NewRouter()
    routes.RegisterAllRoutes(router);
    log.Fatal(http.ListenAndServe(":4000", router))
}
```

This is the entry point of the application. It:
1. Loads environment variables from a .env file
2. Connects to MongoDB
3. Sets up a router using gorilla/mux
4. Registers all routes
5. Starts the HTTP server on port 4000

### 2. Route Registration (routes/index.go)

```go
func RegisterAllRoutes(r *mux.Router) {
    RegisterUserRoutes(r)
    authenticationRouter := r.PathPrefix("/").Subrouter()
    authenticationRouter.Use(middleware.AuthMiddleware)
    WatchlistRoutes(authenticationRouter)
    // Register more route groups as needed
}
```

This function registers all routes in the application. It:
1. Registers user routes (signup and login)
2. Creates a subrouter with authentication middleware for protected routes
3. Registers watchlist routes with the authenticated subrouter

### 3. Authentication Middleware (middleware/authMiddleware.go)
---

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
            http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        userId, err := services.VerifyToken(tokenString)
        if err != nil {
            http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), UserIDKey, userId)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

### 🔹 Function Signature
```go
func AuthMiddleware(next http.Handler) http.Handler
```
This defines a middleware function for HTTP requests. It accepts the `next` handler in the middleware chain and returns a new `http.Handler` that wraps the provided handler with authentication logic.

---

### 🔹 Extracting the Authorization Header
```go
authHeader := r.Header.Get("Authorization")
```
The middleware retrieves the `Authorization` header from the incoming HTTP request. This is where the Bearer token is expected to be passed by the client.

---

### 🔹 Validating Header Presence and Format
```go
if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
    http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
    return
}
```
Checks two things:
- Whether the `Authorization` header is missing (empty string).
- Whether it starts with the expected `"Bearer"` prefix.

If either check fails, it sends an `HTTP 401 Unauthorized` response and halts further processing.

---

### 🔹 Extracting the Token from the Header
```go
tokenString := strings.TrimPrefix(authHeader, "Bearer ")
```
Removes the `"Bearer "` prefix and extracts the raw token string from the header. This token will be passed to the verification service.

---

### 🔹 Verifying the Token
```go
userId, err := services.VerifyToken(tokenString)
if err != nil {
    http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
    return
}
```
The extracted token is passed to the `services.VerifyToken` function, which validates it and returns the `userId` if successful. If the token is invalid or expired, a `401 Unauthorized` is returned along with the error message.

---

### 🔹 Adding User ID to the Request Context
```go
ctx := context.WithValue(r.Context(), UserIDKey, userId)
```
Once the token is verified, the user's ID is stored in the request context. This allows any downstream handlers to access the user ID via `r.Context()`.

---

### 🔹 Proceeding to the Next Handler
```go
next.ServeHTTP(w, r.WithContext(ctx))
```
The middleware then calls the next HTTP handler in the chain, passing along the updated request with the user ID context attached.

---

### 4. JWT Token Creation and Verification (`services/jwt.go`)

```go
func CreateToken(userId string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
    jwt.MapClaims{
        "userId": userId,
        "exp": time.Now().Add(time.Hour*24).Unix(),
    })
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return secretKey, nil
    })
    if err != nil || !token.Valid {
        return "", errors.New("invalid token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", errors.New("could not parse claims")
    }
    userId, ok := claims["userId"].(string)
    if !ok {
        return "", errors.New("user_id not found in token")
    }
    return userId, nil
}
```

---

#### `CreateToken`

1. Creates a new JWT token using the HMAC SHA-256 signing method (`jwt.SigningMethodHS256`).
2. Sets the claims inside the token:
   - `"userId"`: the ID of the user.
   - `"exp"`: the expiration time, which is set to 24 hours from the time of creation.
3. Signs the token using the `secretKey`.
4. Returns the signed token string, or an error if signing fails.

#### `VerifyToken`

1. Parses the JWT token string using `jwt.Parse`.
2. Checks that the token's signing method is HMAC (`jwt.SigningMethodHMAC`).
3. If parsing fails or the token is invalid (e.g., expired or tampered), it returns an error.
4. Extracts the claims from the token using `jwt.MapClaims`.
5. Retrieves the `"userId"` from the claims. If not present or invalid, returns an error.
6. If successful, returns the extracted `userId`.

---

### 5. Database Connection (config/db.go)

```go
func ConnectDB() {
    mongoURI := os.Getenv("MONGO_URI")
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("error in connecting to mongodb", err)
    }
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("MongoDB ping failed:", err)
    }
    DB = client
    log.Println("MongoDb connected successfully")
}

func GetCollection(collectionName string) *mongo.Collection {
    return DB.Database("mongodbapi").Collection(collectionName)
}
```

These functions:
1. ConnectDB: Connects to MongoDB using the URI from environment variables
2. GetCollection: Returns a collection from the "mongodbapi" database

### 6. User Signup (controller/user.go)

```go
func SignupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }
    var user model.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid json", http.StatusBadRequest)
    }
    
    userCollection := config.GetCollection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var existingUser model.User
    
    if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser); err == nil {
        http.Error(w, "User already exists", http.StatusConflict)
        return
    } else if err != mongo.ErrNoDocuments {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error in hashing password", http.StatusInternalServerError)
    }
    user.ID = primitive.NewObjectID()
    user.Password = string(hashedPassword)
    _, err = userCollection.InsertOne(ctx, user)
    if err != nil {
        fmt.Println("Insert error:", err)
        http.Error(w, "Error saving user", http.StatusInternalServerError)
        return
    }

    response := map[string]any{
        "message": "user created successfully",
        "user": map[string]any{
            "id":    user.ID.Hex(),
            "name":  user.Name,
            "email": user.Email,
        },
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

This function:
1. Validates that the request is a POST request
2. Parses the request body into a User struct
3. Checks if a user with the same email already exists
4. Hashes the password using bcrypt
5. Creates a new MongoDB ObjectID for the user
6. Inserts the user into the database
7. Returns a success response with user details (excluding the password)

### 7. User Login (controller/user.go)

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
    }
    var user model.User;
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid json", http.StatusBadRequest)
        return;
    } 
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    userCollection := config.GetCollection("users")
    var existingUser model.User
    if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser); err != nil {
        fmt.Println(err)
        http.Error(w, "User not found", http.StatusNotFound);
        return 
    }
    if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
        http.Error(w, "Password incorrect", http.StatusForbidden)
        return;
    }
    
    token, err := services.CreateToken(existingUser.ID.Hex())
    if err != nil {
        fmt.Println(err)
        http.Error(w, "something went wrong", http.StatusBadGateway)
        return 
    }
    response := map[string]string {
        "token": token,
        "message": "User authentication successful",
    }
    json.NewEncoder(w).Encode(response)
}
```

This function:
1. Validates that the request is a POST request
2. Parses the request body into a User struct
3. Finds the user in the database by email
4. Compares the provided password with the hashed password in the database
5. Creates a JWT token using the user's ID
6. Returns the token in the response

### 8. Add Watchlist Item (controller/watchlist.go)

```go
func AddWatchlistHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only post method allowed", http.StatusMethodNotAllowed)
        return
    }
    var watchlist, existingWatchlist model.Watchlist
    if err := json.NewDecoder(r.Body).Decode(&watchlist); err != nil {
        http.Error(w, "json parsing error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    var watchlistCollection = config.GetCollection("watchlist")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := watchlistCollection.FindOne(ctx, bson.M{"movie_name": watchlist.MovieName}).Decode(&existingWatchlist); err == nil {
        fmt.Println(err)
        http.Error(w, "Movie already in the wartchlist", http.StatusConflict)
        return
    }
    watchlist.ID = primitive.NewObjectID()
    userId := r.Context().Value(middleware.UserIDKey).(string)
    userObjectId, err := primitive.ObjectIDFromHex(userId)
    if err != nil {
        http.Error(w, "Invalid user ID format"+err.Error(), http.StatusBadRequest)
        return
    }
    watchlist.UserID = userObjectId

    // Insert the watchlist into the database
    _, err = watchlistCollection.InsertOne(ctx, watchlist)
    if err != nil {
        http.Error(w, "Error saving watchlist: "+err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]any{
        "message": "Watchlist added successfully",
        "watchlist": map[string]any{
            "id":        watchlist.ID.Hex(),
            "movieName": watchlist.MovieName,
            "watched":   watchlist.Watched,
        },
    }
    json.NewEncoder(w).Encode(response)
}
```

This function:
1. Validates that the request is a POST request
2. Parses the request body into a Watchlist struct
3. Checks if a watchlist with the same movie name already exists
4. Creates a new MongoDB ObjectID for the watchlist
5. Gets the user ID from the request context (set by the authentication middleware)
6. Sets the user ID in the watchlist
7. Inserts the watchlist into the database
8. Returns a success response with watchlist details

## Authentication Flow

1. **User Registration**:
   - User sends a POST request to `/signup` with name, email, and password
   - Server checks if the email is already registered
   - If not, it hashes the password and stores the user in the database
   - Server returns a success response with user details

2. **User Login**:
   - User sends a POST request to `/login` with email and password
   - Server finds the user by email
   - Server compares the provided password with the hashed password
   - If they match, server creates a JWT token with the user ID
   - Server returns the token in the response

3. **Protected Route Access**:
   - User sends a request to a protected route with the JWT token in the Authorization header
   - AuthMiddleware extracts the token from the header
   - AuthMiddleware verifies the token using the VerifyToken function
   - If the token is valid, it adds the user ID to the request context
   - The request is passed to the next handler
   - The handler uses the user ID from the context to perform operations

## Database Operations

### MongoDB Connection

The application connects to MongoDB using the official MongoDB Go driver. The connection is established in the `ConnectDB` function in `config/db.go`:

```go
func ConnectDB() {
    mongoURI := os.Getenv("MONGO_URI")
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("error in connecting to mongodb", err)
    }
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("MongoDB ping failed:", err)
    }
    DB = client
    log.Println("MongoDb connected successfully")
}
```

### Collection Access

Collections are accessed using the `GetCollection` function in `config/db.go`:

```go
func GetCollection(collectionName string) *mongo.Collection {
    return DB.Database("mongodbapi").Collection(collectionName)
}
```

### Common Database Operations

1. **Insert Document**:
   ```go
   _, err = collection.InsertOne(ctx, document)
   ```

2. **Find Document by ID**:
   ```go
   objectId, _ := primitive.ObjectIDFromHex(id)
   err := collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&result)
   ```

3. **Find Document by Field**:
   ```go
   err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
   ```

4. **Find Multiple Documents**:
   ```go
   cursor, err := collection.Find(ctx, bson.M{"user_id": userObjectId})
   var results []Model
   if err = cursor.All(ctx, &results); err != nil {
       // Handle error
   }
   ```

5. **Update Document**:
   ```go
   result, err := collection.UpdateByID(ctx, objectId, bson.M{"$set": updatedDocument})
   ```

6. **Delete Document**:
   ```go
   result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
   ```

## Go Packages Explained

### 1. net/http

The `net/http` package provides HTTP client and server implementations. It's used to:
- Create HTTP handlers
- Handle HTTP requests and responses
- Set response headers and status codes

### 2. encoding/json

The `encoding/json` package implements encoding and decoding of JSON. It's used to:
- Parse JSON request bodies into Go structs
- Convert Go structs to JSON responses

### 3. context

The `context` package defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes. It's used to:
- Set timeouts for database operations
- Store and retrieve values (like user ID) in the request context

### 4. github.com/gorilla/mux

The `gorilla/mux` package implements a request router and dispatcher for matching incoming requests to their respective handler. It's used to:
- Define API routes
- Create subrouters
- Apply middleware to routes
- Extract URL parameters

### 5. go.mongodb.org/mongo-driver

The official MongoDB driver for Go. It's used to:
- Connect to MongoDB
- Perform CRUD operations on MongoDB collections
- Handle MongoDB-specific types like ObjectID

### 6. github.com/golang-jwt/jwt/v5

The `golang-jwt/jwt` package is an implementation of JSON Web Tokens (JWT). It's used to:
- Create JWT tokens with claims
- Verify JWT tokens
- Extract claims from JWT tokens

### 7. golang.org/x/crypto/bcrypt

The `bcrypt` package implements the bcrypt password hashing algorithm. It's used to:
- Hash passwords before storing them in the database
- Compare passwords during login

### 8. github.com/joho/godotenv

The `godotenv` package loads environment variables from a .env file. It's used to:
- Load environment variables like MongoDB URI and JWT secret key

## Conclusion

This Go API project demonstrates a well-structured RESTful API with user authentication and watchlist management functionality. It uses modern Go practices and packages to create a robust and maintainable codebase.

The project showcases:
- API design and implementation
- Authentication with JWT
- MongoDB integration
- Middleware usage
- Error handling
- Request and response processing

This documentation provides a comprehensive overview of the project, its structure, and the key concepts and code used throughout.
