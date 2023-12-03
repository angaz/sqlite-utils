# SQLite Utils

SQLite utils for go, specifically for [modernc / sqlite][modernc-sqlite].

Utils provided:

- URI Builder: Build a connection string for SQLite in a type-safe way.
- Attach database connection hook: The URI builder can encode a parameter for
  adding additional databases to each connection. Running `ATTACH DATABASE`
  will only do it for that one connection. To have it available on all
  connections, use the hook.

[modernc-sqlite]: https://gitlab.com/cznic/sqlite
