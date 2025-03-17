## 0.8.2

- Tasks created with `webserver.Task*` now propagate context cancelation to the
  webserver request context (this is done by wrapping the server's Handler
  in a middleware that sets the request context to the task's context).
- `webserver.ListenAndServer` and `webserver.ListenAndServeTLS` now
  propagate context cancelation to the server's request context (this is done
  by wrapping the server's Handler in a middleware that sets the request
  context to the server's context).
- Updated dependencies.

## 0.8.1

- Added functions to `webserver`:
  - `Task` returns a `work.Task` that runs a web server until context cancelation,
    suitable for use in `work.Pool`.
  - `TaskTLS` - same as `Task`, but for TLS servers.
  - `TaskWithShutdownTimeout` - same as `Task`, but with a custom shutdown
    timeout.
  - `TaskWithShutdownTimeoutTLS` - same as `TaskTLS`, but with a custom shutdown
    timeout.

## 0.8.0

- Added `work` package

## 0.7.0

- Added `env` package

## 0.6.0

- Added `webserver.ServeTLS` and `webserver.ListenAndServeTLS` functions.
- Updated dependencies.

## 0.5.3

- Updated dependencies.

## 0.5.2

- Added `service.Main` function, a convenience function for running a service
  with a `Run` method.

## 0.5.1

- Bump version to work around stale go proxy cache :(

## 0.5.0

- Added `boolean` package
- Added `mapindex` package

## 0.4.1

- Maintenance: upgrade github action versions

## 0.4.0

- Added `service` package
- Added `logger.Logger SafeSync()` method

## 0.3.0

- Added `WithStacktrace` option to `logger` package

## 0.2.1

- Added codecov output to GitHub workflow

## 0.2.0

- Added MIT license

## 0.1.7

- Added `jwt` package

## 0.1.6

- Added `middleware` package
- Added `Sync` method to `logger.Logger`

## 0.1.5

- Added `examples/rest-endpoint` example

## 0.1.4

- Added `rest` package

## 0.1.3

- Added `logger` package

## 0.1.2

- BUGFIX: publish.yaml: fetch tags before looking for an existing tag

## 0.1.1

- `signal` package renamed to `signalcontext` and its `Context` method renamed
  to `WithSignals`.
- Added a GitHub workflow

## 0.1.0

- First release
