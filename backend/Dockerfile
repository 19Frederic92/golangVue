# backend/Dockerfile

# Étape de build
FROM golang:1.20-alpine AS build

# Installer les dépendances
RUN apk add --no-cache git

# Définir le répertoire de travail
WORKDIR /app

# Copier le module Go et télécharger les dépendances
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copier le code source
COPY backend/ .

# Copier le fichier .env
COPY .env .env

# Compiler l'application
RUN go build -o server main.go

# Étape de production
FROM alpine:latest

# Copier le binaire de l'étape de build
COPY --from=build /app/server /server

# Exposer le port 8080
EXPOSE 8080

# Démarrer l'application
CMD ["/server"]