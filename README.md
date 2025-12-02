# ByteFreezer Goodies

Shared Go libraries and utilities for ByteFreezer services. This module provides common functionality used across all ByteFreezer components.

## Packages

### `log`

Structured logging wrapper with level support.

```go
import "github.com/bytefreezer/goodies/log"

log.Info("Server started")
log.Infof("Listening on port %d", 8080)
log.Errorf("Failed to connect: %v", err)
log.Debug("Processing request")
```

Supported log levels: `debug`, `info`, `warn`, `error`, `fatal`

### `control-client`

HTTP client for interacting with ByteFreezer Control Service API.

```go
import controlclient "github.com/bytefreezer/goodies/control-client"

client := controlclient.NewClient(controlclient.Config{
    BaseURL:        "http://control:8082",
    APIKey:         "your-api-key",
    TimeoutSeconds: 30,
})

// Get tenants
tenants, err := client.GetTenants(ctx)

// Get datasets for a tenant
datasets, err := client.GetDatasets(ctx, tenantID)

// Report service health
err := client.ReportHealth(ctx, healthReport)
```

## Installation

```bash
go get github.com/bytefreezer/goodies
```

## Usage in ByteFreezer Services

All ByteFreezer Go services import these packages:

```go
import (
    "github.com/bytefreezer/goodies/log"
    controlclient "github.com/bytefreezer/goodies/control-client"
)
```

## License

ByteFreezer is licensed under the [Elastic License 2.0](LICENSE.txt).

You're free to use, modify, and self-host. You cannot offer it as a managed service.
