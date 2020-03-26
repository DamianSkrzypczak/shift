# TODO (probably in far future) featueres:

- Inter-domain communication mechanism (to keep domains code-spearated)

  - with abstraction on-top of memory/http/ipc/rpc communication

- More AutoAPIs:

  - with custom validation middleware layer
  - using gojay instead of standard json package
  - using graphql
  - using grpc
  - for direct communication with database via custom adapter

- Remove external dependencies

- Create live-reload:

  - if possible/without many requirements from user, make code reload itself via `RunDebug`
  - if not possible/too complex, make external watcher via `shift` command

- Provide some kind of API documentation generation tool

  - first target output format is OpenAPI

- [Low priority] Create web UI with features:
  - for API documentation
  - for API & APP testing
  - for App debugging & bechmarking
