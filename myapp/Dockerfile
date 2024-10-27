# frontend/Dockerfile

# Étape de build
FROM node:18-alpine AS build

# Définir le répertoire de travail
WORKDIR /app

# Installer les dépendances
COPY myapp/package*.json ./
RUN npm install

# Copier le code source
COPY myapp/ .

# Build de l'application
RUN npm run build

# Étape de production
FROM nginx:alpine

# Copier les fichiers de build de Vue.js dans Nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Exposer le port 80 pour Nginx
EXPOSE 80

# Démarrer Nginx
CMD ["nginx", "-g", "daemon off;"]