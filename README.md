<p>
  <img align="left" src="resource/readme/logo.png">
</p>

# Template7-Backend

[![Licence](https://img.shields.io/github/license/Ileriayo/markdown-badges?style=for-the-badge)](./resource/readme/LICENSE)

Template of REST API server write by go.

[API Document](./resource/api-documentation.pdf)

## Architecture

Three layer architecture design: handler -> service -> repo

<p >
  <img src="resource/readme/architecture.png">
</p>

| Block             | Layer | Function                                                                   |
|:------------------|:------|:---------------------------------------------------------------------------|
| API / Route       | 0     | Registered API endpoint.                                                   |
| Middle ware       | 0     | Common / routine functions such like token verify/generate, logging, etc.  |
| Handler           | 1     | Parse necessary variables from URI and body,                               |
| Service           | 2     | Core business logic, include third-party client.                           |
| DB Client         | 3     | Repo layer.                                                                |
| Redis Client      | 3     | Redis client.                                                              |

## Configuration

Settings are read from `config/config.yaml`. Any key can be overridden with an env var:
prefix `T7_`, upper-case, `.` replaced by `_`. Secrets are **not** stored in the config file
and must be provided this way:

| Env var               | Config key           | Note                         |
|:----------------------|:---------------------|:-----------------------------|
| `T7_AUTH_JWTSECRET`   | `auth.jwtSecret`     | required, at least 32 chars  |
| `T7_DB_SQL_PASSWORD`  | `db.sql.password`    |                              |
| `T7_DB_NOSQL_PASSWORD`| `db.nosql.password`  |                              |
| `T7_CACHE_PASSWORD`   | `cache.password`     |                              |

Generate a JWT secret with `openssl rand -hex 32`.

## Run

```
$ export T7_AUTH_JWTSECRET=$(openssl rand -hex 32)
$ make run
```

## Build

### All (Swagger -> Binary -> Run)

```
$ make all 
```

### Binary Only

```
$ make build
```

### Swagger Document

```
$ make swagger
```

### Docker

```
$ docker-compose build
```

## Docker Image

```
$ docker push allensyk/template7-backend:latest
```
