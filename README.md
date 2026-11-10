# Guestbook

Go + TypeScript guestbook with YDB storage, deployed on Yandex Cloud.

For local deployment, install Go (the version in `backend/go.mod`), Node.js,
Docker, and the YC CLI, authenticate with `yc init`, and fill in `scripts/.env`
from `scripts/.env.example`.
Use the existing backend container name for `CONTAINER_NAME`.

```sh
./scripts/schema-create.sh
./scripts/backend-build-publish.sh
./scripts/frontend-build-publish.sh
```

Infrastructure and API Gateway routing must already exist.

The schema script obtains an IAM token using the current YC CLI account and
runs the backend with `--migrate`. This mode applies the existing GORM migrations
and exits without starting the HTTP server or requiring `PORT`. The script also
accepts `YDB_CONNECTION_STRING` from the environment.

Normal backend startup does not apply migrations. Run `./scripts/schema-create.sh`
before deploying backend changes that modify the schema, including changes
deployed through GitHub Actions.

GitHub Actions checks pull requests and deploys both components on pushes to
`main`. To deploy either component manually, use **Actions → CI and Deploy →
Run workflow** on `main`.

Configure these under **Settings → Secrets and variables → Actions**:

| Type | Names |
| --- | --- |
| Secrets | `YC_SA_JSON_CREDENTIALS` (deployment account JSON key), `YDB_CONNECTION_STRING` |
| Variables | `YC_FOLDER_ID`, `REGISTRY_ID`, `SERVICE_ACCOUNT_ID` (backend runtime account), `CONTAINER_NAME`, `BUCKET` |
| Optional variables | (default: `guestbook.vetoshkin.com`), `VITE_API_BASE_URL` (default: `/api`) |
