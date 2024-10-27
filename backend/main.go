package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v4"
    "my_project/backend/config"
)

var dbConn *pgx.Conn

func connectDB() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Erreur de chargement de la configuration:", err)
    }

    // Modifiez cette ligne pour utiliser le nom du service 'db'
    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", 
        cfg.DBUser, cfg.DBPassword, "db", cfg.DBPort, cfg.DBName)

    dbConn, err = pgx.Connect(context.Background(), connStr)
    if err != nil {
        log.Fatal("Unable to connect to database:", err)
    }
    fmt.Println("Connected to PostgreSQL!")
}

// Route pour récupérer des données
func getUsers(c *gin.Context) {
	rows, err := dbConn.Query(context.Background(), "SELECT id, name,age FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		users = append(users, map[string]interface{}{"id": id, "name": name, "age":age})
	}

	c.JSON(http.StatusOK, users)
}

// Route pour ajouter un nouvel utilisateur
func addUser(c *gin.Context) {
    var newUser struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    _, err := dbConn.Exec(context.Background(), "INSERT INTO users (name, age) VALUES ($1, $2)", newUser.Name, newUser.Age)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
}


// Route pour mettre à jour un utilisateur
func updateUser(c *gin.Context) {
    id := c.Param("id")
    var updatedUser struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    if err := c.BindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    _, err := dbConn.Exec(context.Background(), "UPDATE users SET name=$1, age=$2 WHERE id=$3", updatedUser.Name, updatedUser.Age, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Route pour supprimer un utilisateur
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := dbConn.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func main() {
    // Connexion à PostgreSQL
    connectDB()
    defer dbConn.Close(context.Background())

    // Initialisation du routeur Gin
    r := gin.Default()

    // Gestionnaire pour la route /
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API!"})
    })

    // Configuration du middleware CORS
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    })

    // Routes
    r.GET("/users", getUsers)
    r.POST("/users", addUser)
    r.PUT("/users/:id", updateUser)
    r.DELETE("/users/:id", deleteUser)

    // Lancer le serveur
    r.Run(":8080") // Lancer sur localhost:8080
}
