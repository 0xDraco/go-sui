package grpc

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Option configures the client.
type Option func(*config)

type config struct {
	dialOptions          []grpc.DialOption
	transportCredentials credentials.TransportCredentials
	tlsConfig            *tls.Config
}

func defaultConfig() *config {
	return &config{}
}

func WithDialOption(opt grpc.DialOption) Option {
	return func(cfg *config) {
		cfg.dialOptions = append(cfg.dialOptions, opt)
	}
}

// overrides the transport credentials used when dialing.
func WithTransportCredentials(creds credentials.TransportCredentials) Option {
	return func(cfg *config) {
		cfg.transportCredentials = creds
	}
}

// supplies a custom TLS configuration when TLS is enabled by the endpoint.
func WithTLSConfig(tlsCfg *tls.Config) Option {
	return func(cfg *config) {
		if tlsCfg == nil {
			cfg.tlsConfig = nil
			return
		}
		cfg.tlsConfig = tlsCfg.Clone()
	}
}

// forces the client to dial without TLS regardless of the endpoint scheme.
func WithInsecure() Option {
	return func(cfg *config) {
		cfg.transportCredentials = insecure.NewCredentials()
		cfg.tlsConfig = nil
	}
}
