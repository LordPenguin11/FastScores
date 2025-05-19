# LiveScore Backend

A Go backend for real-time sports scores using Fiber, PostgreSQL, and Redis.

## Features
- Fast HTTP API with [Fiber](https://gofiber.io/)
- WebSocket endpoint for live updates
- PostgreSQL for persistent storage
- Redis for real-time pub/sub
- Automated data fetching from Sportmonks API

## Project Structure
```
/cmd/server         # main entrypoint
/internal/api       # HTTP/WebSocket handlers
/internal/db        # PostgreSQL logic
/internal/realtime  # Redis logic
/internal/scheduler # Automated jobs
/admin              # React Admin dashboard
```

## Setup

### Local (manual)
1. **Clone the repo**
2. **Set environment variables:**
   - `REDIS_URL` (e.g. `localhost:6379`)
   - `DATABASE_URL` (e.g. `postgres://user:password@localhost:5432/livescore?sslmode=disable`)
   - `API_BASE_URL` (e.g. `http://localhost:3000`)
   - `SPORTMONKS_CRON_SCHEDULE` (e.g. `0 0 0 1 */3 *` - first day of every 3 months at midnight)
   - `RUN_SPORTMONKS_IMMEDIATELY` (e.g. `true` or `false`)
3. **Install dependencies:**
   ```sh
   go mod tidy
   ```
4. **Run the server:**
   ```sh
   go run cmd/server/main.go
   ```

### Docker (recommended for local dev)
1. **Build and start all services:**
   ```sh
   docker-compose up --build
   ```
2. The app will be available at [http://localhost:3000](http://localhost:3000)
   - **Admin dashboard:** [http://localhost:3001](http://localhost:3001)
   - PostgreSQL: `localhost:5432` (user: `postgres`, password: `postgres`, db: `livescore`)
   - Redis: `localhost:6379`

## How to Connect to the Database (Docker)

1. **Start your containers (if not already running):**
   ```sh
   docker-compose up -d
   ```
2. **Connect to the PostgreSQL container:**
   ```sh
   docker-compose exec db psql -U postgres -d livescore
   ```
3. **Check your database/tables:**
   - List all tables:
     ```sql
     \dt
     ```
   - Show table schema:
     ```sql
     \d+ leagues
     ```
   - Run a query:
     ```sql
     SELECT * FROM leagues LIMIT 5;
     ```
4. **Exit the database prompt:**
   ```
   \q
   ```

## Endpoints
- `GET /health` — Health check
- `GET /leagues` — List all leagues
- `GET /leagues/fetch` — Fetch leagues from Sportmonks API
- `GET /ws` — WebSocket for live scores

## Automated Jobs
The application includes a scheduler that periodically fetches data from Sportmonks API.

### Configuration
You can configure the automated jobs through environment variables:
- `SPORTMONKS_CRON_SCHEDULE`: Cron expression for when to run the job (default: `0 0 0 1 */3 *`, meaning first day of every 3 months at midnight)
- `RUN_SPORTMONKS_IMMEDIATELY`: Whether to run the job immediately on startup (default: `true`)
- `API_BASE_URL`: Base URL for API calls (default: `http://localhost:3000`)

### Manual Trigger
You can manually trigger the data fetch by calling the endpoint:
```
GET /leagues/fetch
```

---

This is a starter scaffold. Extend with your own models, handlers, and business logic! 
