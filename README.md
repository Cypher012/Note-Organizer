# OrganizeNote API

A RESTful API for organizing notes within folders, built with Go and the Fiber web framework. This project demonstrates clean architecture principles, JWT-based authentication, and concurrent load testing capabilities.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Load Testing](#load-testing)
- [Problems Tackled](#problems-tackled)
- [Lessons Learned](#lessons-learned)
- [Future Improvements](#future-improvements)

## Overview

OrganizeNote API provides a hierarchical note organization system where users can create folders and store notes within them. Each user has isolated access to their own folders and notes, ensuring data privacy and security through JWT-based authentication.

## Features

- User registration and authentication with JWT tokens
- HTTP-only cookie-based session management
- CRUD operations for folders
- CRUD operations for notes within folders
- Slug-based URL-friendly resource identification
- Input validation using go-playground/validator
- Concurrent automation testing with configurable load
- SQLite database with GORM ORM

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| Web Framework | Fiber v2 |
| ORM | GORM |
| Database | SQLite |
| Authentication | JWT (golang-jwt/jwt) |
| Validation | go-playground/validator |
| Password Hashing | bcrypt (golang.org/x/crypto) |
| Slug Generation | gosimple/slug |
| Fake Data Generation | brianvoe/gofakeit |

## Project Structure

```
.
├── cmd/
│   ├── server/          # Main application entry point
│   │   └── main.go
│   └── automate/        # Load testing entry point
│       └── main.go
├── internal/
│   ├── automate/        # Automation and load testing logic
│   │   ├── auth.go
│   │   ├── automation.go
│   │   ├── config.go
│   │   ├── cookiejar.go
│   │   ├── folder.go
│   │   ├── note.go
│   │   ├── url.go
│   │   └── userflow.go
│   ├── config/          # Application configuration
│   │   ├── auth-middleware.go
│   │   └── config.go
│   ├── db/              # Database initialization
│   │   └── db.go
│   ├── handlers/        # HTTP request handlers
│   │   ├── auth.go
│   │   ├── folder.go
│   │   ├── note.go
│   │   └── user.go
│   ├── models/          # Data models and DTOs
│   │   ├── auth.go
│   │   ├── folder.go
│   │   ├── note.go
│   │   └── user.go
│   ├── routes/          # Route definitions
│   │   ├── auth.go
│   │   ├── folder.go
│   │   ├── note.go
│   │   ├── routes.go
│   │   └── user.go
│   ├── services/        # Business logic layer
│   │   ├── auth.go
│   │   ├── folder.go
│   │   └── note.go
│   └── utils/           # Utility functions
│       └── auth.go
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Cypher012/OrganizeNoteAPi.git
cd OrganizeNoteAPi
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the project root:
```env
PORT=8080
JWT_SECRET=your-secret-key-here
```

4. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

### Using Air for Hot Reload (Development)

For development with hot reload:
```bash
# Install Air
go install github.com/air-verse/air@latest

# Run with Air
air
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register a new user |
| POST | `/api/auth/login` | Login and receive auth cookies |

### Folders (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/folders` | Get all folders for the authenticated user |
| GET | `/api/folders/:slug` | Get a specific folder by slug |
| POST | `/api/folders` | Create a new folder |
| PUT | `/api/folders/:slug` | Update a folder |
| DELETE | `/api/folders/:slug` | Delete a folder |

### Notes (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/notes` | Get all notes for the authenticated user |
| GET | `/api/folders/:folderSlug/notes` | Get all notes in a folder |
| GET | `/api/folders/:folderSlug/notes/:noteSlug` | Get a specific note |
| POST | `/api/folders/:folderSlug/notes` | Create a note in a folder |
| PUT | `/api/folders/:folderSlug/notes/:noteSlug` | Update a note |
| DELETE | `/api/folders/:folderSlug/notes/:noteSlug` | Delete a note |

### Request/Response Examples

**Register User:**
```json
POST /api/auth/register
{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePassword123"
}
```

**Create Folder:**
```json
POST /api/folders
{
    "name": "Work Projects"
}
```

**Create Note:**
```json
POST /api/folders/work-projects/notes
{
    "name": "Meeting Notes",
    "content": "Discussion points from today's meeting..."
}
```

## Authentication

The API uses JWT tokens stored in HTTP-only cookies for authentication:

- **Access Token**: Short-lived token for API requests
- **Refresh Token**: Long-lived token for obtaining new access tokens

Protected routes require the `access_token` cookie to be present in the request.

## Load Testing

The project includes a built-in load testing system that simulates concurrent user flows:

```bash
go run cmd/automate/main.go
```

### Configuration

The automation system can be configured in `internal/automate/automation.go`:

```go
userCount := 1_000    // Total number of simulated users
concurrency := 100    // Maximum concurrent operations
```

### User Flow

Each simulated user performs the following operations:
1. Register a new account
2. Login
3. Create multiple folders
4. Verify folder retrieval
5. Update folders
6. Create notes in each folder
7. Verify note retrieval
8. Update notes
9. Delete all notes
10. Delete all folders

The semaphore-based concurrency control prevents overwhelming the server while maintaining high throughput.

## Problems Tackled

### 1. Hierarchical Resource Ownership
**Challenge:** Ensuring notes are correctly associated with folders and users across all CRUD operations.

**Solution:** Implemented a folder-slug-based relationship where notes reference their parent folder by slug rather than ID. This provides cleaner URLs and maintains referential integrity through service-layer validation.

### 2. Concurrent Load Testing
**Challenge:** Testing API performance under realistic multi-user load without external tools.

**Solution:** Built an internal automation system using Go's goroutines with semaphore-based rate limiting. This allows controlled concurrent execution of complete user flows.

### 3. Cookie-Based JWT Authentication
**Challenge:** Securely managing authentication tokens across requests.

**Solution:** Implemented HTTP-only cookies for token storage, preventing XSS attacks while maintaining seamless authentication across API calls.

### 4. Slug Generation and Uniqueness
**Challenge:** Creating URL-friendly identifiers while preventing collisions within user scope.

**Solution:** Used the gosimple/slug library with custom configuration, combined with database-level unique constraints scoped to user and parent resource.

### 5. Clean Architecture Separation
**Challenge:** Maintaining testable, maintainable code as the project grew.

**Solution:** Adopted a layered architecture with clear separation between handlers (HTTP), services (business logic), and models (data structures).

## Lessons Learned

### Technical Insights

1. **Fiber Framework Proficiency**: Gained hands-on experience with Fiber's middleware system, context handling, and cookie management.

2. **GORM Relationships**: Learned effective patterns for handling associations, particularly using slugs instead of foreign key IDs for cleaner API design.

3. **Concurrency Patterns**: Implemented semaphore-based rate limiting using channels, understanding the trade-offs between goroutine spawning and controlled concurrency.

4. **JWT Best Practices**: Understood the importance of HTTP-only cookies for token storage and the dual-token (access/refresh) pattern for security.

5. **Input Validation**: Leveraged struct tags with go-playground/validator for declarative validation, reducing boilerplate code.

### Architectural Decisions

1. **Service Layer Pattern**: Separating business logic from HTTP handlers improves testability and allows for easier refactoring.

2. **Error Handling Strategy**: Defined domain-specific errors (ErrFolderNotFound, ErrNoteExists) that map cleanly to HTTP status codes.

3. **URL Structure Design**: RESTful nested resources (`/folders/:slug/notes/:slug`) provide intuitive API navigation.

### Development Workflow

1. **Hot Reload with Air**: Significantly improved development velocity during iteration.

2. **Internal Testing Tools**: Building automation tests alongside the API ensures continuous validation during development.

## Future Improvements

- Add refresh token rotation for enhanced security
- Implement rate limiting middleware
- Add pagination for list endpoints
- Create comprehensive unit and integration tests
- Add API documentation with Swagger/OpenAPI
- Implement note search functionality
- Add folder sharing between users
- Support for note attachments
- WebSocket support for real-time collaboration

## License

This project is for educational purposes.

## Author

Built as part of learning Go web development with Fiber framework.
