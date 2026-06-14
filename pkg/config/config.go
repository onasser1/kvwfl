package config

// TODO: Add options and config to internal package
import "k8s.io/apiserver/pkg/server"

// Configuration that actually runs the server (holding server itself and configuration)
type Config struct {
	SecureServingInfo *server.SecureServingInfo
}
