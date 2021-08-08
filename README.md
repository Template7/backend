<p>
  <img align="left" src="resource/readme/logo.png">
</p>

# Template7-Backend

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)

<br />
<br />

## Architecture

<p >
  <img src="resource/readme/architecture.png">
</p>

| Layer | Function |
| :--- | :--- |
| API/Route | Registered API endpoint. |
| Middle ware | Common/routine functions such like token verification, body check, etc. |
| Handler | Parse necessary variables from URI or body. |
| Component | Core business logic, include third-party client. |
| DB Client | DB access functions. |
| Redis Client | Redis client. |
| Document/Struct | Definition of DB documents/structs. |
