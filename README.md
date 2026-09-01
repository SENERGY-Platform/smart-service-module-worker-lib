# Script Env
worker may receive scripts in the inputs named "prescript" and "postscript". 
those scripts may be chunked by appending suffixes to the input which. 
all inputs/scripts with the "prescript" prefix are sort by name, concatenated and executed before the worker code is run.
all inputs/scipts with the "postscript" prefix are sort by name, concatenated and executed after the worker code is run.

the script environment allows access to multiple apis to 
- handle smart-service variables
- read inputs
- write outputs
- access the device-repository
- ...

to allow the web-ui (https://github.com/SENERGY-Platform/web-ui) code completion in https://github.com/SENERGY-Platform/web-ui/tree/master/src/app/modules/smart-services/designer/dialog/edit-smart-service-task-dialog,
this repository provides a code generator, that creates an ace completer by calling
```
go generate ./...
```

# Use Example

```
import (
	"context"
	"flag"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-process/pkg"
	"github.com/SENERGY-Platform/smart-service-module-worker-process/pkg/processdeployment"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	configLocation := flag.String("config", "config.json", "configuration file")
	flag.Parse()

	libConfig, err := configuration.LoadLibConfig(*configLocation)
	if err != nil {
		log.Fatal(err)
	}
	config, err := configuration.Load[processdeployment.Config](*configLocation)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	err = pkg.Start(ctx, wg, config, libConfig)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		sig := <-shutdown
		log.Println("received shutdown signal", sig)
		cancel()
	}()

	<-ctx.Done() //waiting for context end; may happen by shutdown signal
	wg.Wait()
}

```

```
func Start(ctx context.Context, wg *sync.WaitGroup, config processdeployment.Config, libConfig configuration.Config) error {
	handlerFactory := func(auth *auth.Auth, smartServiceRepo *smartservicerepository.SmartServiceRepository) (camunda.Handler, error) {
		return processdeployment.New(config, libConfig, auth, smartServiceRepo), nil
	}
	return lib.Start(ctx, wg, libConfig, handlerFactory)
}
```
# OpenTelemetry

Traces are sent to an OTLP collector over gRPC. The endpoint is set with the config value
`otel_endpoint` (env `OTEL_ENDPOINT`); if it is empty, the default
`jaeger.logging.svc.cluster.local:4317` is used.

The service name is not a constant of this library. Every worker using it is its own service and
has to show up in jaeger under its own name, so `configuration.ServiceName()` derives it from the
main module of the binary (e.g. `github.com/SENERGY-Platform/smart-service-module-worker-analytics`
becomes `smart-service-module-worker-analytics`). Workers that create spans of their own should use
the same function.

What is instrumented:

- one span per fetched camunda task, named after the worker topic, with the task-id, the
  process-instance-id and the activity-id as attributes. The poll (`fetchAndLock`) itself gets no
  span; it runs continuously and mostly returns nothing.
- the context of a task is derived with `context.WithoutCancel`: a shutdown does not cut off a task
  that is already locked and running.
- outgoing requests to the smart-service-repository and the camunda calls that belong to a task
  (`complete`, error and stop). They carry trace-context and baggage on, so a task can be followed
  across services.
- log records, through the open-telemetry handler of `struct-logger`. It writes the baggage into the
  record and the record into the current span, but **only for the context based log methods**
  (`ErrorContext`, `InfoContext`, ...). A `logger.Error()` without a context stays untraced.

`middleware.Do` adds the id of the smart-service-instance to the baggage as
`smart_service_instance_id` before the worker handler runs, so every following request and log
record of that task carries it. The camunda business-key is not used for this: for maintenance
procedures it holds a maintenance-id, only the smart-service-repository resolves the
process-instance-id to the instance.

Not propagated: the keycloak token-exchange and the device-repository client, which has no context
aware methods.
