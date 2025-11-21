# ThynGo
**API REST en Go (Gin) — backend modulaire complet avec MongoDB.**

[![CI/CD Status](https://github.com/RaphaelMailhiot/thyngo/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/RaphaelMailhiot/thyngo/actions/workflows/ci-cd.yml)

ThynGo est une API REST développée en Go et conçue pour offrir une architecture solide, testable et facile à maintenir. L’architecture est organisée en modules clairs (handlers, services, repositories) et le projet est prêt pour un déploiement en production grâce à Docker et GitHub Actions.

## ✨ Fonctionnalités principales

- **Endpoints REST complets (CRUD)**
- **Architecture modulaire et testable** (handlers, services, repositories)
- **Connexion configurable à MongoDB** via variables d’environnement
- **Tests unitaires et intégration continue (CI/CD)** avec GitHub Actions
- **Conteneurisation** avec Docker & Docker Compose

## 🔧 Prérequis

- **Go 1.25+**
- **Docker & Docker Compose**
- **MongoDB** (si exécution sans Docker Compose)

## 🚀 Lancer le projet en local

### 1. Cloner le dépôt

```sh
git clone https://github.com/RaphaelMailhiot/thyngo.git
cd thyngo
```

### 2. Installer les dépendances

```sh
go mod download
```

### 3. Configurer les variables d’environnement

Configurer les variables nécessaires (par exemple) :

- `APP_PORT`
- `MONGO_URI`
- `etc.`

### 4. Démarrer l’application

#### Avec Docker Compose (recommandé)

```sh
docker compose up --build
```

#### Sans Docker Compose

```sh
go run ./cmd/api
```

L’API sera accessible à l’adresse suivante : `http://localhost:8080`

## 🧪 Exécuter les tests

Lancer tous les tests unitaires du projet :

```sh
go test ./... -v
```

## 🏗️ Architecture & Stack technique

- **Langage :** Go
- **Framework web :** Gin
- **Base de données :** MongoDB
- **Conteneurisation :** Docker
- **CI/CD :** GitHub Actions