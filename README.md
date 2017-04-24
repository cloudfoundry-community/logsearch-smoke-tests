Logsearch-for-CloudFoundry Smoke Tests
======================================

## Running the tests

### Set up your `go` environment

Set up your golang development environment, [per golang.org](http://golang.org/doc/install).

See [Go CLI](https://github.com/cloudfoundry/cli) for instructions on
installing the go version of `cf`.

Make sure that [curl](http://curl.haxx.se/) is installed on your system.

Make sure that the go version of `cf` is accessible in your `$PATH`.

All `go` dependencies required by the smoke tests are vendored in
`logsearch-for-cloufoundry-smoke-tests/vendor`.

### Test Setup

To run the Logsearch-for-CloudFoundry Smoke tests, you will need:
- a running CF instance
- an environment variable `$CONFIG` which points to a `.json` file that
contains the CF settings

Below is an example `integration_config.json`:
```json
{
  "api":                              "api.example.com",
  "apps_domain":                      "example.com",
  "system_domain":                    "example.com",
  "admin_user":                       "admin",
  "admin_password":                   "admin",
  "skip_ssl_validation":              true
}
```

### Test Execution

To execute the tests, run:

```bash
./bin/test
```

Internally the `bin/test` script runs tests using [ginkgo](https://github.com/onsi/ginkgo).
