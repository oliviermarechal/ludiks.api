# Jobs Cron Serverless - Synchronisation des Souscriptions

## Vue d'ensemble

Solution simple pour exécuter des tâches cron dans une architecture serverless. Chaque job s'exécute une fois et se termine, parfait pour les hébergeurs cloud.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Server    │    │  Cron Job       │    │   Database      │
│   (Serverless)  │    │  (Serverless)   │    │   (Shared)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌──────────────────┐
                    │     Stripe       │
                    │   (External)     │
                    └──────────────────┘
```

## Déploiement sur Scaleway

### 1. Build de l'image

```bash
# Build de l'image
docker build -f Dockerfile.job -t registry.scw.cloud/ludiks/subscription-sync .

# Push vers le registry Scaleway
docker push registry.scw.cloud/ludiks/subscription-sync
```

### 2. Configuration du job cron dans Scaleway

#### Option A : Scaleway Containers (Recommandé)

1. **Aller dans la console Scaleway**
   - Container → Containers → Create Container

2. **Configuration du container**
   ```
   Name: subscription-sync
   Image: registry.scw.cloud/ludiks/subscription-sync
   CPU: 100m
   Memory: 128Mi
   Timeout: 5 minutes
   ```

3. **Variables d'environnement**
   ```
   DATABASE_URL=postgres://user:password@host:5432/db
   STRIPE_SECRET_KEY=sk_test_...
   ```

4. **Configuration du cron**
   ```
   Schedule: 0 2 * * *  # Tous les jours à 2h du matin
   ```

#### Option B : Scaleway Functions

```bash
# Build pour les fonctions
GOOS=linux GOARCH=amd64 go build -o subscription-sync ./cmd/subscription-sync

# Déployer
scw function deploy \
  --name subscription-sync \
  --runtime go120 \
  --handler subscription-sync \
  --env DATABASE_URL=$DATABASE_URL \
  --env STRIPE_SECRET_KEY=$STRIPE_SECRET_KEY
```

### 3. Configuration avancée Scaleway

#### Variables d'environnement recommandées
```bash
# Base de données
DATABASE_URL=postgres://user:password@host:5432/db?sslmode=require

# Stripe
STRIPE_SECRET_KEY=sk_live_...  # ou sk_test_... pour les tests

# Logs
LOG_LEVEL=info
```

#### Ressources recommandées
```bash
# Pour un job simple
CPU: 100m
Memory: 128Mi
Timeout: 5 minutes

# Pour un job plus complexe
CPU: 200m
Memory: 256Mi
Timeout: 10 minutes
```

## Configuration du cron

### Fréquences recommandées

```bash
# Tous les jours à 2h du matin (recommandé)
0 2 * * *

# Toutes les 12h
0 */12 * * *

# Tous les lundis à 6h
0 6 * * 1

# Toutes les heures (pour les tests)
0 * * * *
```

### Monitoring et logs

Le job retourne :
- **Exit code 0** : Succès
- **Exit code 1** : Erreur

**Logs disponibles dans :**
- Console Scaleway → Containers → Logs
- Ou via CLI : `scw container logs subscription-sync`

## Compatibilité Serverless

### ✅ Compatible
- **Gin framework** : Parfait pour serverless
- **Variables d'environnement** : Bien configuré
- **Architecture hexagonale** : Excellente séparation

### ⚠️ Points d'attention
- **Migrations** : Ne pas exécuter automatiquement en serverless
- **Connexions DB** : Ajouter un pool de connexions
- **Timeouts** : Configurer des timeouts appropriés

### 🔧 Adaptations recommandées

```go
// Dans main.go de l'API
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    ConnPool: &gorm.ConnPool{
        MaxIdleConns: 5,
        MaxOpenConns: 10,
        ConnMaxLifetime: time.Hour,
    },
})

// Ajouter des timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

## Avantages de cette approche

✅ **Simple** : Un seul fichier main.go par job  
✅ **Serverless** : Pas de service qui tourne 24h/24  
✅ **Économique** : Payez seulement quand le job s'exécute  
✅ **Scalable** : Géré par l'hébergeur  
✅ **Maintenable** : Logique métier séparée de l'infrastructure  

## Ajouter d'autres jobs

Pour ajouter un nouveau job cron :

1. **Créer un nouveau dossier** : `cmd/autre-job/`
2. **Créer main.go** avec la logique spécifique
3. **Créer un Dockerfile** : `Dockerfile.autre-job`
4. **Configurer le cron** dans Scaleway

Exemple :
```bash
cmd/
├── subscription-sync/
│   └── main.go
├── cleanup-logs/
│   └── main.go
└── backup-database/
    └── main.go
```

## Troubleshooting

### Job ne s'exécute pas
- Vérifier la configuration cron dans Scaleway
- Vérifier les variables d'environnement
- Vérifier les logs dans la console

### Job échoue
- Vérifier la connexion à la base de données
- Vérifier les permissions
- Vérifier les logs détaillés

### Performance
- Augmenter CPU/Memory si nécessaire
- Optimiser les requêtes SQL
- Ajouter des index si besoin

### Erreurs courantes
```bash
# Connexion DB échoue
# Solution : Vérifier DATABASE_URL et SSL

# Timeout
# Solution : Augmenter le timeout dans Scaleway

# Memory insuffisante
# Solution : Augmenter la mémoire allouée
``` 