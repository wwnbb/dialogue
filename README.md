# Dialogue: A Go Web Framework Based on Functions

#### Dialogue is a sleek web framework for Go that revolves around the concept of utilizing straightforward functions to process web requests. By simplifying the routing and handler function definitions, Dialogue enables developers to concentrate on building the logic of their applications.

### Getting Started

### Prerequisites

```
    Go 1.15 or higher
```

## Basic Usage

The following snippet is a basic example of how to set up a Dialogue application, showcasing the creation of a routing function, starting the server, and implementing graceful shutdown on receiving an OS signal.

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	D "github.com/wwnbb/dialogue"
)

func main() {
    Application := D.Switch(ApplicationRoutes)  // Define your routes

    // Start the server
    server, err := D.ListenAndServe(
        ":8089",
        D.Chain(
            Application,
            D.NotFoundHandler(),
        ),
    )
    if err != nil {
        log.Fatalf("Could not start the server: %v", err)
        return
    }

    // Listen for OS signals
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

    // Wait for an OS signal
    <-signals

    // Shutdown the server gracefully with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    server.Shutdown(ctx)

    log.Println("Exiting gracefully")
}

```

##### In this example:

1. Creating a New Dialogue Application:
A new Dialogue application is initialized using the D.Switch function with ApplicationRoutes as the argument. ApplicationRoutes should be a map where the keys are path strings and the values are handler functions.

2. Starting the Server:
The server is started on port 8089 using the D.ListenAndServe function. The D.Chain function is used to chain the main application handler and a not-found handler together.

3. Listening for OS Signals:
The application listens for SIGINT and SIGTERM signals to initiate a graceful shutdown.

4. Graceful Shutdown:
Upon receiving a shutdown signal, the server is shut down gracefully with a timeout of 5 seconds.

## Further Exploration

For a deeper understanding and exploration of Dialogue's features, visit the documentation and browse through the examples provided in the repository.
Contributing

We appreciate your contributions! Please see the CONTRIBUTING.md file for more information.

## License

Dialogue is distributed under the GNU License. See the LICENSE file for more details.
