# Unlimited Texas Hold'em Poker

## Service Frontends

### Start the poker service

Build the web assets first if you want `pokerd` to serve the browser UI itself:

```powershell
cd D:\Git\go\src\poker\web
npm install
npm run build
```

Then start the service:

```powershell
cd D:\Git\go\src\poker
go run ./cmd/pokerd -addr 127.0.0.1:8080 -web-dist web/dist
```

This serves:

- `http://127.0.0.1:8080/` for the bundled web console
- `http://127.0.0.1:8080/api/rooms` for the room API
- `ws://127.0.0.1:8080/ws/rooms/{roomId}` for live room updates

### Start both frontend and backend together

For day-to-day development, use the repo script:

```powershell
cd D:\Git\go\src\poker
.\scripts\dev-up.ps1
```

It will:

- stop stale repo frontend/backend dev processes first
- stop listeners on `:8080` and `:4173`
- install `web/` dependencies if `node_modules` is missing
- open one PowerShell window for `pokerd`
- open one PowerShell window for the Vite frontend

### Start the web frontend in dev mode

```powershell
cd D:\Git\go\src\poker\web
npm install
npm run dev -- --host 127.0.0.1 --port 4173
```

The Vite dev server proxies `/api` and `/ws` to `http://127.0.0.1:8080`.

### Use the terminal frontend

List rooms:

```powershell
cd D:\Git\go\src\poker
go run ./cmd/pokerctl -server http://127.0.0.1:8080
```

Watch a room as a spectator:

```powershell
go run ./cmd/pokerctl -server http://127.0.0.1:8080 -room room-001 -watch
```

Take the human seat, start a hand, or submit an action:

```powershell
go run ./cmd/pokerctl -server http://127.0.0.1:8080 -room room-001 -viewer-seat 5 -take-seat
go run ./cmd/pokerctl -server http://127.0.0.1:8080 -room room-001 -start
go run ./cmd/pokerctl -server http://127.0.0.1:8080 -room room-001 -token turn-1 -action FOLD
```
