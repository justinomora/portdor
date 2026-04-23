# portdor

**portdor** is a local dev service registry and manager — the `localhost` port manager that rules them all. You start your processes yourself; portdor lets you name them, track their ports and projects, check health on demand, and stop, kill, or restart them from the CLI or a Web UI at `localhost:4242`. It keeps a persistent registry across sessions so you never have to remember what's running where.

---

## Installation

```sh
go install github.com/justinomora/portdor/cmd/portdor@latest
```

To build from source:

```sh
git clone https://github.com/justinomora/portdor
cd portdor
make build       # produces ./portdor in the repo root
make install     # copies it to /usr/local/bin so portdor works from anywhere
```

You can also run `./portdor` directly from the repo root without installing.

---

## Quick Start

The portdor server starts automatically the first time you run any command, so there's no mandatory `serve` step. A typical workflow looks like:

```sh
# Register a service you've already started
portdor register api --cmd "npm run dev" --port 3000 --cwd ./api --project myapp

# List all registered services
portdor list

# Open the Web UI
portdor ui
```

If you have a `portdor.toml` in your project, register everything at once:

```sh
portdor up    # register all services from portdor.toml
portdor down  # stop and unregister all services from portdor.toml
```

To start the server explicitly in the foreground:

```sh
portdor serve
```

---

## portdor.toml

Place a `portdor.toml` in your project root to define services as a group. `portdor up` reads this file and registers all services with portdor; `portdor down` stops and unregisters them. portdor does **not** start the processes — you do that yourself.

### Field Reference

#### `[project]`

| Field  | Required | Description                                                   |
| ------ | -------- | ------------------------------------------------------------- |
| `name` | optional | Project label used to group services in `list` and the Web UI |

#### `[[services]]`

Each service is declared as a `[[services]]` table (one per service).

| Field     | Required | Description                                                                            |
| --------- | -------- | -------------------------------------------------------------------------------------- |
| `name`    | required | Unique service name used in all CLI commands                                           |
| `command` | required | The command that starts the service (used by `restart`)                                |
| `port`    | optional | Port the service listens on; shown in `list` and the Web UI                            |
| `cwd`     | optional | Working directory for the command; defaults to the directory containing `portdor.toml` |

### First-time workflow

1. Add a `portdor.toml` to your project root (see example below)
2. Start your services as you normally would (portdor does not start them for you)
3. From the project root, run:

```sh
portdor up
```

That registers all services defined in the file. From then on, use `portdor list`, `portdor ui`, or the CLI commands to manage them. When you're done, `portdor down` stops and unregisters them all.

### Example

```toml
[project]
name = "myapp"

[[services]]
name    = "api"
port    = 3000
command = "npm run dev"
cwd     = "./api"

[[services]]
name    = "frontend"
port    = 3001
command = "npm start"
cwd     = "./web"

[[services]]
name    = "worker"
command = "python manage.py worker"
cwd     = "./worker"
```

---

## CLI Reference

### `portdor serve`

Start the portdor server in the foreground on `:4242`. Not required — any other command auto-starts the server if it isn't already running.

### `portdor register <name> --cmd <cmd> [--port N] [--cwd path] [--project name]`

Register a service with portdor.

| Flag        | Required | Description                                       |
| ----------- | -------- | ------------------------------------------------- |
| `--cmd`     | required | Command used to start/restart the service         |
| `--port`    | optional | Port the service listens on                       |
| `--cwd`     | optional | Working directory (defaults to current directory) |
| `--project` | optional | Project label for grouping                        |

### `portdor unregister <name>`

Remove a service from the registry.

### `portdor update <name> [--name x] [--project x] [--port N] [--cmd x] [--cwd x]`

Update one or more fields of a registered service. Only the flags you pass are changed.

### `portdor list`

List all registered services with their project, port, and status.

### `portdor status <name>`

Show full JSON detail for a single service.

### `portdor stop <name>`

Gracefully stop a service (SIGTERM).

### `portdor kill <name>`

Force kill a service (SIGKILL).

### `portdor restart <name>`

Restart a service using its registered command and working directory.

### `portdor up`

Register all services defined in `portdor.toml` in the current directory. Services are not started — only added to the registry.

### `portdor down`

Stop and unregister all services defined in `portdor.toml` in the current directory.

### `portdor ui`

Open the portdor Web UI at `http://localhost:4242` in your default browser.

---

## Web UI

The Web UI is available at `http://localhost:4242` (or run `portdor ui` to open it directly).

It shows all registered services grouped by project. For each service you can:

- **Stop** — send SIGTERM
- **Force Kill** — send SIGKILL
- **Restart** — restart using the registered command

Click a service name to edit its fields (name, project, port, command, working directory) inline without leaving the page.
