# Frontend

Run the backend on port `8080`, then start the frontend:

```sh
npm install
npm run dev
```

Vite forwards `/api/messages` to `http://localhost:8080/api/messages`, preserving the path. To use a different local backend address:

```sh
BACKEND_URL=http://localhost:3000 npm run dev
```

The frontend loads messages with `GET /api/messages` and publishes `{ name, content }` with `POST /api/messages`. After a successful publish, it fetches the saved feed again. Empty names are sent as `Anonymous`.

Requests use `/api` by default. Set `VITE_API_BASE_URL` when building if a different API base is needed. The local proxy is for development only.

In production, route `/api/*` to the backend without rewriting the path and serve the static frontend for other paths.

```sh
npm run build
```
