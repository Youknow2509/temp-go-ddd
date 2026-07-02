package http

import (
	"crypto/tls"
	"fmt"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

// buildTLSConfig converts HttpTlsSetting into a *tls.Config.
// Returns (nil, nil) when TLS is disabled.
func buildTLSConfig(cfg domain_config.HttpTlsSetting) (*tls.Config, error) {
	if !cfg.IsEnabled {
		return nil, nil
	}

	minVersion, err := parseTLSVersion(cfg.MinVersion)
	if err != nil {
		return nil, err
	}

	tlsCfg := &tls.Config{
		MinVersion: minVersion,
	}

	if cfg.RequireClientCert {
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return tlsCfg, nil
}

func parseTLSVersion(v string) (uint16, error) {
	switch v {
	case "", "1.2", "TLS1.2", "TLSv1.2":
		return tls.VersionTLS12, nil
	case "1.3", "TLS1.3", "TLSv1.3":
		return tls.VersionTLS13, nil
	case "1.1", "TLS1.1", "TLSv1.1":
		return tls.VersionTLS11, nil
	case "1.0", "TLS1.0", "TLSv1.0":
		return tls.VersionTLS10, nil
	default:
		return 0, fmt.Errorf("unsupported tls min_version: %q", v)
	}
}
