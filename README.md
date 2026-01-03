# Beaver ETL Context

## FILE STRUCTURE
  - cmd/etl – CLI entrypoint, flags/env config, wiring.
  - internal/manifest – manifest loading + validation; keep your existing parser logic here or wrap it.
  - internal/matcher – glob expansion, ignore/exclude handling, conflict resolution (last-write wins for link/language/contributors, tag union).
  - internal/resolver – apply defaults + derive title/language/link/contributors for each matched file.
  - internal/gitmeta – repo URL, commit SHA, blame contributors; keep shelling to git in one place.
  - internal/content – file reading, binary detection, size limits.
  - internal/api – Beaver API client (POST /api/v1/snippets), request structs, auth header handling.
  - internal/pipeline – orchestrates: load manifest → match → resolve → read content → build request → send.
  - internal/logging (optional) – structured logging + progress output.


This repository ingests files described by a repository-level manifest (`beaver.yaml`) and creates
snippets in the Beaver API (`../beaver_api`). The API accepts any snippet content (code, text, docs).
The links below point to the main files needed to build that pipeline.

## Manifest format
- Spec: `BEAVER_MANIFEST_FORMAT.md`
- File name: `beaver.yaml` (repo root)

## beaver_api (snippet ingestion target)
- API overview: `../beaver_api/README.md`
- Routes: `../beaver_api/internal/router/router.go`
- Snippet endpoints: `../beaver_api/internal/router/snippets.go`
- Request/response models: `../beaver_api/internal/router/models.go`
- Request parsing + mapping: `../beaver_api/internal/router/utils.go`
- Snippet service logic: `../beaver_api/internal/service/snippets.go`
- Auth middleware: `../beaver_api/internal/router/middleware.go`
- Service access tokens: `../beaver_api/internal/router/service-access-tokens.go`

### Snippet ingestion endpoint
`POST /api/v1/snippets` requires `Authorization: Bearer <token>`.
The request body matches `IngestSnippetRequest`:
```json
{
  "title": "Binary Search",
  "code": "…",
  "project_url": "https://en.wikipedia.org/wiki/Binary_search_algorithm",
  "git_repo_url": { "name": "https://github.com/your-org/algorithms" },
  "git_path": "c/search/binary_search/main.c",
  "git_version": "commit-sha",
  "language": { "name": "c" },
  "tags": [{ "name": "search" }, { "name": "divide_and_conquer" }],
  "contributors": [
    { "first_name": "Krishna", "last_name": "Vedala", "email": "7001608+kvedala@users.noreply.github.com" }
  ]
}
```

Notes:
- `git_repo_url` is an object; `CreateGit` uses `json:"name"` for the URL field, so send `{ "name": "…" }`.
- `project_url` is optional (`omitempty`) but should come from the manifest `link` or repo URL fallback.
- `git_path` should be the matched file path from `beaver.yaml`.
- `git_version` should be the commit SHA (or other version string) from the source repo.

### Authentication options
Snippet ingestion is protected by `authMiddleware` and accepts bearer tokens.
If the pipeline needs a long-lived token, use:
- `POST /api/v1/service-access-tokens` (requires normal auth) to create one.

## Ingestion pipeline mapping (manifest -> beaver_api)
The pipeline can ingest any repository with files described in `beaver.yaml`.
1. Locate `beaver.yaml` in the repo root and resolve `sources`, `ignore`, and defaults.
2. For each matched file:
   - Read file content (code or text).
   - Build `IngestSnippetRequest`:
     - `title` <- manifest title (file source) or filename-derived
     - `project_url` <- manifest link or repo URL fallback
     - `language.name` <- resolved language (`auto` or explicit)
     - `tags[].name` <- resolved tags
     - `contributors[]` <- resolved contributors (or `git blame`)
     - `git_repo_url.name` <- source repo URL
     - `git_path` <- matched file path
     - `git_version` <- current commit SHA for source repo
     - `code` <- contents of the matched file
3. `POST /api/v1/snippets` with `Authorization: Bearer <token>`.
