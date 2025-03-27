package metier

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/jackc/pgx/v4"
    "my_project/backend/config"
)

var dbConn *pgx.Conn
var clients = make(map[*websocket.Conn]bool) // Map pour suivre les clients connectés
var broadcast = make(chan Message)           // Canal pour diffuser les messages

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Autoriser toutes les origines (à adapter selon vos besoins)
    },
}

type Message struct {
    Type string `json:"type"`
    Data string `json:"data"`
}

// ConnectDB initialise la connexion à la base de données
func ConnectDB() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Erreur de chargement de la configuration:", err)
    }

    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

    dbConn, err = pgx.Connect(context.Background(), connStr)
    if err != nil {
        log.Fatal("Unable to connect to database for entite:", err)
    }
    fmt.Println("Connected to PostgreSQL!")
}

// DBConn retourne la connexion à la base de données
func DBConn() *pgx.Conn {
    return dbConn
}

// GetEntite récupère toutes les entités de la base de données
func GetEntite(c *gin.Context) {
    if dbConn == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
        return
    }

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

// GetOneEntite récupère une entité spécifique
func GetOneEntite(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr) // Convertir l'ID en entier
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token entite"})
        return
    }

    log.Printf("Fetching token with ID: %d", id)

    var entite struct {
        ID    int    `json:"id"`
        Name  string `json:"name"`
        Image string `json:"image"`
    }

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

    c.JSON(http.StatusOK, entite)
}

// WebSocket endpoint
func serveWs(c *gin.Context) {
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    clients[ws] = true

    for {
        _, _, err := ws.ReadMessage()
        if err != nil {
            delete(clients, ws) // Supprimer le client déconnecté
            break
        }
    }
}

// Broadcast updates to all clients
func broadcastUpdates() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

// Notify clients about database changes
func notifyClients(msg Message) {
    broadcast <- msg
}

// Example of updating the database and notifying clients
func updateEntite(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    // Exemple de mise à jour de la base de données
    _, err = dbConn.Exec(context.Background(), "UPDATE entite SET name=$1 WHERE id=$2", "Nouveau Nom", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entite"})
        return
    }

    // Notifier les clients après la mise à jour
    notifyClients(Message{Type: "update", Data: "Entité mise à jour"})
    c.JSON(http.StatusOK, gin.H{"message": "Entité mise à jour"})
}

// Assurez-vous de démarrer la diffusion des mises à jour
func init() {
    go broadcastUpdates()
}
