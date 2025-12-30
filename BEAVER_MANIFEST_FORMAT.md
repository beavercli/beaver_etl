# Beaver Manifest Format

This document defines the new repository manifest format used by Beaver ingestion. It replaces the
old single-snippet `beaver.json` layout with a single repo-level manifest that can describe many
files via explicit paths or glob patterns.

The manifest is authored as YAML:
- `beaver.yaml`

## Location and naming
- File name: `beaver.yaml`
- Location: repository root
- Paths and globs are evaluated relative to the repository root using `/` separators.

## Goals
- One manifest per repository.
- Generic: ingest any file type (code, markdown, docs, configs, etc.).
- Simple defaults with predictable fallbacks.
- Optional per-source metadata (tags, link, language, contributors).

---

## Schema (overview)

Required:
- `version` (integer)
- `sources` (array)

Optional:
- `defaults`
- `ignore`

---

## Full YAML example

```yaml
version: 1

defaults:
  tags: []
  link: null
  language: auto
  contributors: auto

ignore:
  - "**/node_modules/**"
  - "**/.git/**"

sources:
  - type: file
    path: README.md
    title: Beaver CLI Overview
    tags: [docs, entry]
    link: https://github.com/owner/repo#readme

  - type: pattern
    include:
      - "docs/**/*.md"
    tags: [docs]

  - type: pattern
    include:
      - "src/**/*.go"
      - "src/**/*.py"
    tags: [code]
```

## Field reference

### `version` (required, integer)
Schema version. Start at `1` and increment when breaking changes are introduced.

### `defaults` (optional, object)
Defaults applied to all sources unless overridden at the source level.

- `tags` (array of strings)
  - Default tags to apply to all ingested files.
- `link` (string or null)
  - Default URL to associate with ingested files. Use `null` to omit.
- `language` (`"auto"` or string)
  - `auto` detects language from file extension.
  - A string forces a language identifier for all matched files. Unknown extensions are treated as `text`.
- `contributors` (`"auto"` or array of contributor objects)
  - `auto` means infer contributors from `git blame` per file.
  - Array entries follow the `contributors` format below.

### `ignore` (optional, array of strings)
Glob patterns that are excluded from all sources. Ignored paths are removed after include
matching and before per-source `exclude` is applied.

### `sources` (required, array)
Each entry describes one ingestion source. There are two types: `file` and `pattern`.

Common fields for all sources:
- `type` (string, required)
  - Either `file` or `pattern`.
- `tags` (array of strings, optional)
- `link` (string or null, optional)
- `language` (`"auto"` or string, optional)
- `contributors` (`"auto"` or array of contributor objects, optional)

#### Source type: `file`
Required fields:
- `path` (string)
  - Path to a single file, relative to repo root.

Optional fields:
- `title` (string)
  - Explicit title for this file. If omitted, title is derived from the filename.

#### Source type: `pattern`
Required fields:
- `include` (array of strings)
  - Glob patterns that match multiple files.

Optional fields:
- `exclude` (array of strings)
  - Glob patterns to exclude for this source.

Notes:
- `title` is ignored for `pattern` sources. Titles are derived from filenames.

### `contributors` object format
Contributor entries contain:
- `name` (string)
- `last_name` (string)
- `email` (string)

---

## Implicit metadata (derived from git/CI)

These values are inferred unless explicitly overridden in future schema versions:
- `repo.url` from `git remote origin` or CI-provided repository URL.
- `repo.name` from `repo.url`.
- `git_version` from `git rev-parse HEAD`.
- `git_path` from the matched file path.

---

## Matching and conflict rules

- `ignore` is applied globally before per-source `exclude`.
- `exclude` applies only to the source where it is defined.
- If a file is matched by multiple sources, the ingestion is deterministic:
  - Tags are unioned across matches.
  - Duplicate tag names are removed.
  - `link`, `language`, and `contributors` use the last matching source in manifest order.
  - `title` comes from a `file` source if present, otherwise derived from filename.
- Unmatched sources are ignored.

---

## Content handling

- Files are read as text. If a file is detected as binary or exceeds the ingestion size limit, it
  may be skipped (implementation-defined).

---

## Resolution and fallback rules

For each ingested file, resolve fields in this order:

- `language`:
  1. Source `language`
  2. `defaults.language`
  3. `auto` detection by file extension

- `tags`:
  1. `defaults.tags` (base)
  2. Source `tags` (added on top)

- `link`:
  1. Source `link`
  2. `defaults.link`
  3. Repository URL (implicit)

- `contributors`:
  1. Source `contributors`
  2. `defaults.contributors`
  3. `auto` via `git blame` (per file). If blame is unavailable, contributors are omitted.

- `title`:
  1. Source `title` (only for `file` sources)
  2. Filename-derived title (for all files)

---

## Ingestion mapping (manifest -> Beaver API)

For each matched file:
- `title` -> `title`
- `link` -> `project_url`
- `language` -> `language.name`
- `tags[]` -> `tags[].name`
- `contributors[]` -> `contributors[]`
- file path -> `git_path`
- file content -> `code`
- repo URL (implicit) -> `git_repo_url.name`
- git SHA (implicit) -> `git_version`

---

## Validation guidance

- `version` must be present and supported.
- `sources` must be non-empty.
- `file` sources must point to an existing file.
- `pattern` sources must have non-empty `include` globs.
- `contributors` entries must include `name`, `last_name`, and `email` when provided.
- Paths must remain within the repository root (no `..` escapes).
- When globs match directories, those entries are ignored.
