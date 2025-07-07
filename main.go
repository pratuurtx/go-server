package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// User represents a sample data model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory database
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())  // Logs all HTTP requests
	e.Use(middleware.Recover()) // Recovers from panics

	// Routes
	e.GET("/", hello)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.POST("/users", createUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler for the root route
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// Handler to get all users
func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

// Handler to get a specific user by ID
func getUser(c echo.Context) error {
	id := c.Param("id")
	for _, user := range users {
		if strconv.Itoa(user.ID) == id {
			return c.JSON(http.StatusOK, user)
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "User not found")
}

// Handler to create a new user
func createUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Generate a new ID (in a real app, use a database auto-increment)
	user.ID = len(users) + 1
	users = append(users, *user)
	return c.JSON(http.StatusCreated, user)
}

// Handler to update an existing user
func updateUser(c echo.Context) error {
	id := c.Param("id")
	
	// Find the user
	for i := range users {
		if strconv.Itoa(users[i].ID) == id {
			updatedUser := new(User)
			if err := c.Bind(updatedUser); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			
			// Update user fields
			users[i].Name = updatedUser.Name
			users[i].Email = updatedUser.Email
			return c.JSON(http.StatusOK, users[i])
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "User not found")
}

// Handler to delete a user
func deleteUser(c echo.Context) error {
	id := c.Param("id")
	for i := range users {
		if strconv.Itoa(users[i].ID) == id {
			// Remove the user from the slice
			users = append(users[:i], users[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "User not found")
}