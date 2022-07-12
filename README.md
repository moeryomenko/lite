# Lite [![Go Reference](https://pkg.go.dev/badge/github.com/moeryomenko/lite.svg)](https://pkg.go.dev/github.com/moeryomenko/lite)

Lite is package contains compact version of [squad](http://github.com/moeryomenko/squad) with auto liveness and readiness checks. 

## Usage

```go
package main

import (
	"context"
	"time"

	"github.com/moeryomenko/healing"
	"github.com/moeryomenko/healing/decorators/pgx"
	"github.com/moeryomenko/lite"
)

func main() {
	h := lite.New(8081 // health controller port.
		healing.WithCheckPeriod(3 * time.Second),
		healing.WithReadinessTimeout(time.Second),
		healing.WithReadyEndpoint("/readz"),
	)

	// create postgresql pool.
	pool, err := pgx.New(ctx, pgx.Config{
		Host:     pgHost,
		Port:     pgPort,
		User:     pgUser,
		Password: pgPassword,
		DBName:   pgName,
	}, pgx.WithHealthCheckPeriod(100 * time.Millisecond)) // sets the duration between checks of the health of idle conn.
	if err != nil {
		// error handling.
		...
	}

	// add pool readiness controller to readiness group.
	l.AddReadyChecker("pgx", pool.CheckReadinessProber)

    // Run your service.
	err = l.Run(svc.Run)
	// log err if service stopped with error.
}
```

## License

Lite is primarily distributed under the terms of both the MIT license and Apache License (Version 2.0).

See [LICENSE-APACHE](LICENSE-APACHE) and/or [LICENSE-MIT](LICENSE-MIT) for details.
