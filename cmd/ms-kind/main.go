package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/NorskHelsenett/ror-ms-kind/internal/httpserver"
	"github.com/NorskHelsenett/ror-ms-kind/internal/kindconfig"
	"github.com/NorskHelsenett/ror-ms-kind/internal/mskindconnections"
	"github.com/NorskHelsenett/ror-ms-kind/internal/rabbitmq/mskindrabbitmqdefinitions"
	"github.com/NorskHelsenett/ror-ms-kind/internal/rabbitmq/mskindrabbitmqhandler"
	"github.com/NorskHelsenett/ror-ms-kind/internal/rorclient"

	"syscall"

	"github.com/NorskHelsenett/ror/pkg/config/configconsts"

	"github.com/NorskHelsenett/ror/pkg/telemetry/trace"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	"github.com/spf13/viper"

	// https://blog.devgenius.io/know-gomaxprocs-before-deploying-your-go-app-to-kubernetes-7a458fb63af1
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(rlog.Infof))
}
func main() {
	cancelChan := make(chan os.Signal, 1)
	stop := make(chan struct{})
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	rlog.Info("Kind micro service starting")
	kindconfig.Load()

	mskindconnections.InitConnections()

	mskindrabbitmqdefinitions.InitOrDie()
	rorclient.SetupRORClient()

	go func() {
		httpserver.InitHttpServer()
	}()

	if viper.GetBool(configconsts.ENABLE_TRACING) {
		go func() {
			trace.ConnectTracer(stop, viper.GetString(configconsts.ROLE), viper.GetString(configconsts.OPENTELEMETRY_COLLECTOR_ENDPOINT))
			sig := <-cancelChan
			_, _ = fmt.Println()
			_, _ = fmt.Println(sig)
			stop <- struct{}{}
		}()
	}

	mskindrabbitmqhandler.StartListening()

	sig := <-cancelChan
	rlog.Info("Caught signal", rlog.Any("singal", sig))
	// shutdown other goroutines gracefully
	// close other resources
}
