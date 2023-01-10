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
go generate ./... > ace-code-completer.ts
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