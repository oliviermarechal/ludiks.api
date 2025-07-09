# Makefile pour Ludiks API
# Usage: make run-subscription-sync

# Charger les variables d'environnement depuis .env
include .env
export

# Variables par défaut
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Variables
REGISTRY = rg.fr-par.scw.cloud
IMAGE_NAME = ludiks/api
TAG = latest
FULL_IMAGE = $(REGISTRY)/$(IMAGE_NAME):$(TAG)

.PHONY: help run-subscription-sync build-subscription-sync test-subscription-sync

# Afficher l'aide
help:
	@echo "Available commands:"
	@echo "  build        - Build Docker image for amd64 architecture"
	@echo "  push         - Push image to Scaleway registry"
	@echo "  build-push   - Build and push in one command"
	@echo "  build-push-tag - Build and push with custom tag"
	@echo "  clean        - Remove local image"
	@echo "  help         - Show this help"

# Lancer le job de synchronisation des souscriptions
run-subscription-sync:
	@echo "🚀 Lancement du job de synchronisation des souscriptions..."
	@echo "📋 Configuration:"
	@echo "  DB_HOST: $(DB_HOST)"
	@echo "  DB_NAME: $(DB_NAME)"
	@echo "  DB_PORT: $(DB_PORT)"
	@echo "  STRIPE_SECRET_KEY: $(shell echo $(STRIPE_SECRET_KEY) | cut -c1-10)..."
	@echo ""
	@go run cmd/subscription-sync/main.go

# Compiler le job de synchronisation
build-subscription-sync:
	@echo "🔨 Compilation du job de synchronisation..."
	@go build -o bin/subscription-sync cmd/subscription-sync/main.go
	@echo "✅ Compilation terminée: bin/subscription-sync"

# Lancer le job compilé
run-subscription-sync-binary: build-subscription-sync
	@echo "🚀 Lancement du job compilé..."
	@./bin/subscription-sync

# Nettoyer les fichiers compilés
clean:
	@echo "🧹 Nettoyage des fichiers compilés..."
	@rm -rf bin/
	@echo "✅ Nettoyage terminé"

check-env:
	@if [ -z "$(DB_HOST)" ]; then \
		echo "❌ DB_HOST n'est pas défini dans .env"; \
		exit 1; \
	fi
	@if [ -z "$(DB_USER)" ]; then \
		echo "❌ DB_USER n'est pas défini dans .env"; \
		exit 1; \
	fi
	@if [ -z "$(DB_PASSWORD)" ]; then \
		echo "❌ DB_PASSWORD n'est pas défini dans .env"; \
		exit 1; \
	fi
	@if [ -z "$(DB_NAME)" ]; then \
		echo "❌ DB_NAME n'est pas défini dans .env"; \
		exit 1; \
	fi
	@if [ -z "$(DB_PORT)" ]; then \
		echo "❌ DB_PORT n'est pas défini dans .env"; \
		exit 1; \
	fi
	@if [ -z "$(STRIPE_SECRET_KEY)" ]; then \
		echo "❌ STRIPE_SECRET_KEY n'est pas défini dans .env"; \
		exit 1; \
	fi
	@echo "✅ Variables d'environnement OK"

run-subscription-sync-safe: check-env run-subscription-sync 

.PHONY: build
build:
	docker build --platform=linux/amd64 -t $(FULL_IMAGE) .

# Push l'image vers le registry Scaleway
.PHONY: push
push:
	docker push $(FULL_IMAGE)

.PHONY: build-push
build-push: build push

.PHONY: build-push-tag
build-push-tag:
	@read -p "Enter tag: " tag; \
	docker build --platform=linux/amd64 -t $(REGISTRY)/$(IMAGE_NAME):$$tag .; \
	docker push $(REGISTRY)/$(IMAGE_NAME):$$tag

.PHONY: clean
clean:
	docker rmi $(FULL_IMAGE) 2>/dev/null || true 