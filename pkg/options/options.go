package options

// TODO: Add options and config to internal package
import (
	"github.com/onasser1/validating-kontroller/pkg/config"
	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/klog/v2"
)

const ValidatingKontroller string = "validating-kontroller"

// Options is needed for setting up the HTTPS server for Admission Webhook (options and flags)
type Options struct {
	SecureServingOptions options.SecureServingOptions
}

func NewDefaultOptions() *Options {
	opts := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	opts.SecureServingOptions.BindPort = 8443
	opts.SecureServingOptions.ServerCert.PairName = ValidatingKontroller
	return opts
}

func (o *Options) AddFlagSet(fs *pflag.FlagSet) {
	o.SecureServingOptions.AddFlags(fs)
}

func (o *Options) Config() *config.Config {
	if err := o.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		klog.Fatalf("error creating HTTPS server config: %v", err)
	}
	config := &config.Config{}
	o.SecureServingOptions.ApplyTo(&config.SecureServingInfo)
	return config
}
