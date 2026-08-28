<div align="center">
  <img src="client/icons/blendkit_logo.png" alt="Logo" width="100" height="100"/>
  <h3 align="center">Blendkit-Client</h3>

  Local bridge between Blendkit DCC add-ons and the Blendkit service.

  [![Project license](https://img.shields.io/github/license/blenderkit/blenderkit.svg?color=orange)](LICENSE)
</div>

## About

This repository contains the **Blendkit-Client** (formerly *daemon*) — a small
local HTTP server written in Go. It runs on the user's machine and bridges the
Blendkit DCC add-ons with the [Blendkit service](https://www.blendkit.com/):
it handles search, downloads, uploads, login/OAuth2, ratings, comments and more,
and reports progress back to the connected software.

It is designed to be used in two ways:

- **As a subrepository** embedded into the Blendkit add-ons:
  - [`blenderkit_addon`](https://github.com/BlenderKit/blenderkit) (Blender)
  - `bk_maya` (Maya)
  - `bk_godot` (Godot)
  - `blenderkit_rhino` (Rhino)
- **As a standalone app** with its own settings and a small set of features.

The Client is built for Windows, macOS and Linux on both x86_64 and arm64. Multiple
Client versions can coexist on a single machine, and several can run at the same
time (each on its own port), so different DCCs/add-on versions are supported
side by side.

## Repository layout

```
client/                 Go source for the Blendkit-Client
  *.go                  HTTP handlers, networking, login, download, wrappers…
  internal/apispec/     Route registry — single source of truth for the API
  internal/config/      Standalone on-disk Client settings (no secrets)
  cmd/apidocgen/        Generator: renders docs/openapi.json and docs/API.md
  docs/                 Generated API documentation (OpenAPI 3.1 + Markdown)
  tools/                Bundled Python "recipes" run under headless Blender
  icons/                Branding assets + embedded tray icon (blendkit.ico)
  VERSION               Current Client version
dev.py                  Developer helper (build / run / verify / test / lint / docs)
pyproject.toml          Python tooling config (ruff, pydoclint) for tools/
```

## Developer quick start

Requirements: a recent **Go** toolchain (see [`client/go.mod`](client/go.mod)) and
**Python 3.10+** with the dev tools (`ruff`, `pydoclint`) for linting the recipes.

### Git hooks (pre-commit)

This repo uses [pre-commit](https://pre-commit.com) to run `ruff`, `pydoclint`
and `bandit` before each commit. The hook config lives in
[`.pre-commit-config.yaml`](.pre-commit-config.yaml); the generated hook under
`.git/hooks/` is local-only and not versioned, so each contributor installs it
once:

```sh
pip install pre-commit          # or: uv tool install pre-commit
pre-commit install              # generate .git/hooks/pre-commit from the config
pre-commit run --all-files      # optional: run the hooks against the whole repo
```

All commands go through [`dev.py`](dev.py):

```sh
python dev.py build            # cross-compile the Client for all platforms -> ./out/vX.Y.Z/
python dev.py test             # run Go unit tests + lint the Python recipes
python dev.py lint             # ruff + pydoclint on client/tools (no changes)
python dev.py format           # ruff format + auto-fix on client/tools
python dev.py docs             # regenerate the API documentation
python dev.py verify <path>    # verify code-signing/notarization of built binaries
```

You can also work directly in the `client/` directory with the Go toolchain:

```sh
cd client
go test ./...                       # run all Go tests (incl. the API drift test)
go test -benchmem -run=^$ -bench=.  # run performance benchmarks
go generate ./...                   # regenerate docs/openapi.json and docs/API.md
go build .                          # build a Client for the current platform
```

## API

The Client exposes a local HTTP API (default port **62485**). It is documented from
a single source of truth and published in two formats:

- [`client/docs/openapi.json`](client/docs/openapi.json) — OpenAPI 3.1 (Postman/Insomnia, Swagger UI/Redoc, SDK generation).
- [`client/docs/API.md`](client/docs/API.md) — human-readable reference.

See [`client/README.md`](client/README.md) for how the Client is started and how the
docs stay in sync with the code.

## Roadmap

Planned work for this repository (tracked incrementally):

- **Local Client configuration** stored next to the binary (`blendkit-client-config.json`),
  excluding the API key, which must be stored with maximum available security.
  Foundation in [`client/internal/config`](client/internal/config/config.go); wiring it
  into Client startup and a settings UI is next.
- Settings change notifications pushed to connected DCCs.
- A taskbar/menubar tray icon to open settings, view logs, etc. *(Windows
  implemented: shown automatically for standalone runs; macOS/Linux pending —
  they need a cgo build.)*
- Automated, signed builds on Git, consumable live by developers or as packaged
  binaries by end users.

## License

GPL-2.0-or-later. See [LICENSE](LICENSE).
