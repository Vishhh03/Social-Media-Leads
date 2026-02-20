---
name: Docker Compose Rebuilds
description: How to ensure code and CSS changes effectively apply inside a Docker container workflow.
---
# Docker Compose Rebuilds

## Key Findings & Caveats
1. **`docker compose restart` is Not Enough:** Restarting a container only re-runs the existing built image. For Vite frontend apps or compiled Go backend APIs, a `restart` command will NOT compile the new code changes.
2. **Aggressive Layer Caching:** Docker will often cache the `npm install` and `npm run build` layers to save time. If you apply CSS or JSX fixes locally, Docker may ignore them and fail to push them to the NGINX serving layer because it cached the build step globally.
3. **The Solution:** Force a clean rebuild using the `--no-cache` flag to bypass the cached layers entirely:
   ```bash
   docker compose build --no-cache frontend
   docker compose up -d frontend
   ```
4. **Go Backend Rebuilds:** If a `.go` file is changed, the backend must be rebuilt. `docker compose up -d --build backend` will re-compile the binary and restart the container cleanly. If errors persist, try appending `--no-cache` to force a clean `go mod download` and `go build`.
