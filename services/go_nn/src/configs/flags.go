package configs

import (
	"flag"
)

func ParseFlags(settings *Settings) {

	var host, port, protocol string

	flag.StringVar(&host, "h", "", "Host service")
	flag.StringVar(&port, "p", "", "Port")
	flag.StringVar(&protocol, "m", "", "Protocol: grpc | http")
	flag.Parse()

	// Update settings
	if host != "" {
		settings.ServiceHost = host
	}
	if port != "" {
		settings.ServicePort = port
	}
	if protocol != "" {
		settings.ServiceProtocol = protocol
	}
}
