package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/onasser1/validating-kontroller/pkg/options"
	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/klog/v2"
)

func main() {
	opts := options.NewDefaultOptions()

	flagSet := pflag.NewFlagSet(options.ValidatingKontroller, pflag.ExitOnError)
	globalflag.AddGlobalFlags(flagSet, options.ValidatingKontroller)
	opts.AddFlagSet(flagSet)

	if err := flagSet.Parse(os.Args); err != nil {
		klog.Fatalf("error reading flags from flagset: %v", err)
	}
	c := opts.Config()

	// mux is the http handler we would use for the https server
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(ServeValidationAdmission))

	ch := server.SetupSignalHandler()
	stoppedCh, listerStoppedCh, err := c.SecureServingInfo.Serve(mux, 30*time.Second, ch)
	if err != nil {
		klog.Fatalf("error from Serve function: %v", err)
	} else {
		<-stoppedCh
		<-listerStoppedCh
	}

}

func ServeValidationAdmission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hi, we are ServeValidationAdmission function, we are called")
}
