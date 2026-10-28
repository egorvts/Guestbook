# Guestbook

Go + TypeScript guestbook with YDB storage, deployed on Yandex Cloud.

For local deployment, install Node.js, Docker, and the YC CLI, authenticate with
`yc init`, and fill in `scripts/.env` from `scripts/.env.example`.
Use the existing backend container name for `CONTAINER_NAME`.

```sh
./scripts/backend-build-publish.sh
./scripts/frontend-build-publish.sh
```

Infrastructure and API Gateway routing must already exist.

GitHub Actions checks pull requests and deploys both components on pushes to
`main`. To deploy either component manually, use **Actions → CI and Deploy →
Run workflow** on `main`.

Configure these under **Settings → Secrets and variables → Actions**:

| Type | Names |
| --- | --- |
| Secrets | `YC_SA_JSON_CREDENTIALS` (deployment account JSON key), `YDB_CONNECTION_STRING` |
| Variables | `YC_FOLDER_ID`, `REGISTRY_ID`, `SERVICE_ACCOUNT_ID` (backend runtime account), `CONTAINER_NAME` |
| Optional variables | `BUCKET` (default: `guestbook.vetoshkin.com`), `VITE_API_BASE_URL` (default: `/api`) |
