package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "database/sql"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v4"
    "my_project/backend/config"
    "my_project/backend/metier"
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


func getTokens(c *gin.Context) {
	rows, err := dbConn.Query(context.Background(), "SELECT id, name,price,supply FROM token")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch token"})
		return
	}
	defer rows.Close()

	var tokens []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
        var price float64
		var supply float64
		if err := rows.Scan(&id, &name, &price,&supply); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		tokens = append(tokens, map[string]interface{}{"id": id, "name": name, "price":price, "supply":supply})
	}

	c.JSON(http.StatusOK, tokens)
}

func getOneToken(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr) // Convertir l'ID en entier
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
        return
    }

    log.Printf("Fetching token with ID: %d", id)

    // Préparer la structure du token
    var token struct {
        ID     int     `json:"id"`
        Name   string  `json:"name"`
        Price  float64 `json:"price"`
        Supply float64 `json:"supply"`
    }

    // Exécuter la requête pour récupérer le token
    err = dbConn.QueryRow(context.Background(), "SELECT id, name, price, supply FROM tokens WHERE id=$1", id).Scan(&token.ID, &token.Name, &token.Price, &token.Supply)

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("No token found for ID: %d", id)
            c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
            return
        }
        log.Printf("Error fetching token: %v", err) // Log l'erreur exacte
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch one token"})
        return
    }

    // Retourner le token en format JSON
    c.JSON(http.StatusOK, token)
}



//Pour être exporté en go la fonction doit commencer par une majsucule
/*
func GetEntite(c *gin.Context) {
	rows, err := dbConn.Query(context.Background(), "SELECT id, name, image FROM entite")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entite"})
		return
	}
	defer rows.Close()

	var entite []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var image string

		if err := rows.Scan(&id, &name, &image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		entite = append(entite, map[string]interface{}{"id": id, "name": name, "image": image})
	}

	c.JSON(http.StatusOK, entite)
}


func GetOneEntite(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr) // Convertir l'ID en entier
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token entite"})
        return
    }

    log.Printf("Fetching token with ID: %d", id)

    // Préparer la structure du token
    var entite struct {
        ID     int     `json:"id"`
        Name   string  `json:"name"`
        Image  string  `json:"image"`
      
    }

    // Exécuter la requête pour récupérer le token
    err = dbConn.QueryRow(context.Background(), "SELECT id, name, image FROM entite WHERE id=$1", id).Scan(&entite.ID, &entite.Name, &entite.Image)

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("No entite found for ID: %d", id)
            c.JSON(http.StatusNotFound, gin.H{"error": "Entité not found"})
            return
        }
        log.Printf("Error fetching entite: %v", err) // Log l'erreur exacte
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch one entite"})
        return
    }

    // Retourner l'entite' en format JSON
    c.JSON(http.StatusOK, entite)
}
*/

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
    metier.ConnectDB()

    // Vérifiez que la connexion est bien initialisée
    if metier.DBConn() == nil {
        log.Fatal("Database connection failed to initialize")
    }

    // Fermeture de la connexion à la base de données lorsque le programme se termine
    defer metier.DBConn().Close(context.Background())

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
    r.GET("/tokens",getTokens)
    //r.GET("/tokens/:id", getOneToken)
   r.GET("/entites", metier.GetEntite)
   r.GET("/entites/:id", metier.GetOneEntite)
 // r.GET("/entites", GetEntite)
    //r.GET("/entites/:id", GetOneEntite)
    
    r.POST("/users", addUser)
    r.PUT("/users/:id", updateUser)
    r.DELETE("/users/:id", deleteUser)

    // Lancer le serveur
    r.Run(":8080") // Lancer sur localhost:8080
}
