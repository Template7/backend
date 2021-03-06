<p>
  <img align="left" src="resource/readme/logo.png">
</p>

# Template7-Backend

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)

Template of REST API server write by go.

[API Document](./resource/api-documentation.pdf)

## Architecture

<p >
  <img src="resource/readme/architecture.png">
</p>

For clean logic and easy-maintainable purpose, suggest each layer to access its next layer's function / method only,
do not implement cross layer function call.

For example: handler should not access db client directly, suggest to access by correspond component instead.  

| Layer | Function |
| :--- | :--- |
| API / Route | Registered API endpoint. |
| Middle ware | Common / routine functions such like token verification, body check, etc. |
| Handler | Parse necessary variables from URI or body. |
| Component | Core business logic, include third-party client. |
| DB Client | DB manipulation functions. |
| Redis Client | Redis client. |
| Document / Struct | Definition of DB documents / structs, it could be referenced by any layer. |

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
