# Social Media Lead Automation (SaaS)

A multi-tenant lead management platform integrating Meta (WhatsApp/Instagram/Facebook) APIs with a custom CRM. Built for high concurrency and real-time communication.

## Tech Stack & Architecture

### Backend (Go)
- **Framework**: [Gin](https://gin-gonic.com/) for high-performance HTTP routing.
- **Database**: [PostgreSQL](https://www.postgresql.org/) with `pgx` driver and connection pooling.
- **Cache/Queue**: [Redis](https://redis.io/) for session management, rate-limiting (leaky bucket), and job deduplication.
- **Design Pattern**: Domain-Driven Design (DDD) with Clean Architecture layers:
  - `handlers` (HTTP transport)
  - `store` (Data access / Repository pattern)
  - `models` (Domain entities)
  - `meta` (External API integration)

### Frontend (React)
- **Build Tool**: [Vite](https://vitejs.dev/) for instant HMR and optimized builds.
- **State**: React Context API + LocalStorage for auth persistence.
- **Styling**: Vanilla CSS with a custom variable-based design system (no heavy UI frameworks).
- **Communication**: REST API with JWT authentication (stateless).

### Infrastructure
- **Containerization**: Docker & Docker Compose for full-stack orchestration.
- **Reverse Proxy**: Nginx (production) / Vite Proxy (development).

## Key Features

- **Multi-channel Inbox**: Unified chat interface for WhatsApp, Instagram, and Facebook Messenger.
- **Automated Workflows**: Keyword-based triggers with regex support for auto-replies.
- **Broadcast System**: Bulk messaging with Redis-backed deduplication and rate limiting (preventing Meta policy violations).
- **Authentication**:
  - Email/Password (Bcrypt hashing)
  - Google OAuth 2.0 (OpenID Connect)
  - JWT Session management

## Local Development

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (optional, if running outside Docker)
- Node.js 20+ (optional, if running outside Docker)

### Quick Start (Windows)
Run the included dev script to spin up the entire stack locally:
```bash
./dev.bat
```
*Starts Backend (:8080) and Frontend (:3000) concurrently.*

### Manual Setup
1. **Database & Redis**:
   ```bash
   docker compose up -d db redis
   ```
2. **Backend**:
   ```bash
   cd backend
   go run ./cmd/api
   ```
3. **Frontend**:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Configuration
Copy `.env.example` to `.env` and populate:
- `DB_*`: Database credentials (default: `leadbot`/`leadautomation`)
- `JWT_SECRET`: Secure random string for token signing
- `META_*`: App credentials from Meta Developer Portal
- `GOOGLE_*`: OAuth client ID/Secret from Google Cloud console

## Deployment
Production builds use a multi-stage Dockerfile for the frontend (Node build -> Nginx alpine) and a scratch-based image for the Go binary to keep image sizes minimal (<50MB).

```bash
docker compose up -d --build
```
