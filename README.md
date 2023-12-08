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

## Run

```
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
