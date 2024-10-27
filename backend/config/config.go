
package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

// Config struct pour stocker la configuration de l'application
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
}

// LoadConfig charge les variables d'environnement et retourne une structure de configuration
func LoadConfig() (*Config, error) {
    // Charger les variables d'environnement depuis un fichier .env
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Aucun fichier .env trouvé, chargement des variables d'environnement système")
    }

    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "mydb"),
    }

    return config, nil
}

// getEnv récupère une variable d'environnement ou retourne une valeur par défaut si elle est absente
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}