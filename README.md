# Setlist PWA 🎸

Bienvenue sur le projet **Setlist PWA** ! Cette application est une *Progressive Web App* conçue pour aider les groupes de musique à créer, gérer et partager leurs setlists.

## ✨ Fonctionnalités

*   **Gestion de Bibliothèque** : Ajoutez et organisez toutes vos chansons et interludes.
*   **Création de Setlists** : Créez des setlists personnalisées par simple glisser-déposer.
*   **Calcul de Durée** : Calcule automatiquement la durée totale de vos concerts.
*   **Mode PWA** : Installez l'application sur votre téléphone pour un accès rapide.
*   **Export PDF** : Générez un PDF propre de vos setlists pour l'impression ou le partage.
*   **Collaboratif** : Partagez l'accès à votre groupe pour que tout le monde soit synchronisé.

## 🛠️ Stack Technique

Ce projet est un monorepo qui combine un backend en Go avec un frontend moderne en SvelteKit.

*   **Frontend** (`/frontend`):
    *   **Framework** : [SvelteKit](https://kit.svelte.dev/) (avec Svelte 5 & Runes)
    *   **Langage** : TypeScript
    *   **Styling** : [Tailwind CSS](https://tailwindcss.com/)
    *   **PWA** : `@vite-pwa/sveltekit`
    *   **Tests** : Vitest (Unitaires) & Playwright (End-to-End)

*   **Backend** (`/backend`):
    *   **Langage** : [Go](https://go.dev/) (API REST native sans framework)
    *   **Base de données** : PostgreSQL
    *   **Authentification** : JWT (JSON Web Tokens)
    *   **Migrations** : `golang-migrate`

*   **Environnement** :
    *   **Orchestration** : [Docker](https://www.docker.com/) & Docker Compose
    *   **Automatisation** : `Makefile` pour des commandes simplifiées
    *   **Intégration Continue** : GitHub Actions

## 🚀 Démarrage Rapide

Le projet est entièrement conteneurisé avec Docker, ce qui simplifie grandement l'installation.

### Prérequis

Avant de commencer, assurez-vous d'avoir installé les outils suivants sur votre machine :

1.  **[Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)** : Pour lancer l'environnement de développement.
2.  **[Make](https://www.gnu.org/software/make/)** : Pour utiliser les commandes simplifiées du `Makefile`. (généralement pré-installé sur Linux/macOS).
3.  **[Go](https://go.dev/doc/install)** (v1.21+) : Uniquement si vous devez gérer les migrations manuellement.
4.  **[Node.js](https://nodejs.org/)** (v22+) : Uniquement si vous souhaitez développer le frontend en dehors de Docker.

### 1. Configuration Initiale

Commencez par cloner le projet et configurer vos variables d'environnement.

```bash
# Clonez le dépôt
git clone <URL_DU_REPO>
cd <NOM_DU_REPO>

# Copiez le fichier d'environnement d'exemple
cp .env.example .env
```

> **Important** : Le fichier `.env` contient des secrets (comme `JWT_SECRET`). Il est ignoré par Git et ne doit jamais être partagé.

### 2. Lancer l'environnement de développement

Grâce au `Makefile`, lancer tous les services est un jeu d'enfant.

```bash
# Construit les images Docker et démarre tous les services (frontend, backend, DB)
make up
```

Une fois la commande terminée :
*   🚀 **Le Frontend** sera accessible sur [http://localhost:4000](http://localhost:4000)
*   ⚙️ **L'API Backend** sera accessible sur [http://localhost:8089](http://localhost:8089)

Les migrations de la base de données sont appliquées automatiquement au démarrage grâce au service `migrator`.

### Commandes `Makefile` utiles

Voici les commandes que vous utiliserez le plus souvent :

*   `make up` : Démarre tous les services en arrière-plan.
*   `make down` : Arrête et supprime tous les conteneurs.
*   `make logs` : Affiche les logs de tous les services en temps réel.
*   `make shell-back` : Ouvre un terminal `sh` dans le conteneur du backend.
*   `make shell-front` : Ouvre un terminal `sh` dans le conteneur du frontend.

Pour voir toutes les commandes, consultez le fichier `Makefile`.

## 🧪 Lancer les Tests

Le projet est équipé d'une suite de tests complète.

### Tests Unitaires (Frontend)

Ces tests vérifient de petites parties isolées du code frontend (fonctions utilitaires, etc.).

```bash
# Lance les tests unitaires une seule fois
make test-unit

# Lance les tests unitaires en mode "watch" pour le développement
make test-unit-watch
```

### Tests End-to-End (Playwright)

Ces tests simulent une interaction utilisateur complète dans un navigateur. Ils nécessitent un environnement de test dédié.

```bash
# Lance toute la suite de tests E2E (crée un environnement de test complet)
make test
```

> Après l'exécution, un rapport de test HTML détaillé est généré dans `frontend/playwright-report/index.html`.

Vous pouvez également lancer des sous-ensembles de tests (ex: `make test-song-list`) pour un débogage plus rapide.

## 🤝 Comment Contribuer

Ce projet est ouvert aux contributions ! Pour garantir la qualité et la cohérence du code, veuillez suivre ces quelques règles :

1.  **Ne jamais pusher sur `main`** : La branche `main` est protégée. Toutes les modifications doivent passer par une Pull Request.

2.  **Workflow de contribution** :
    *   **Créez une branche** : Partez de la branche `main` la plus à jour. Choisissez un nom de branche explicite (ex: `feat/add-dark-mode` ou `fix/login-bug`).
        ```bash
        git switch main
        git pull origin main
        git switch -c feat/ma-nouvelle-fonctionnalite
        ```
    *   **Développez** : Codez votre fonctionnalité ou votre correctif.
    *   **Testez** : Assurez-vous que vos modifications ne cassent rien en lançant les tests (`make test-all`). Idéalement, ajoutez de nouveaux tests pour couvrir votre code.
    *   **Ouvrez une Pull Request (PR)** : Poussez votre branche sur le dépôt distant et ouvrez une PR vers `main`. Décrivez clairement vos changements.
    *   **Revue de code** : Un mainteneur examinera votre code. Une fois la PR approuvée et les tests de l'intégration continue (CI) au vert, elle sera fusionnée.

3.  **Standards de code** :
    *   Le code est formaté avec **Prettier** et vérifié avec **ESLint**. Avant de commiter, vous pouvez lancer `npm run lint` et `npm run format` dans le dossier `frontend` pour vous assurer que tout est en ordre.
