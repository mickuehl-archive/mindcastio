package environment

import (
	"os"
	"strconv"
	"strings"
)

const (
	LISTEN_PORT            string = "LISTEN_PORT"
	STATS_HOST             string = "STATS_HOST"
	STATS_PORT             string = "STATS_PORT"
	STATS_ENABLED          string = "STATS_ENABLED"
	BACKEND_HOSTS          string = "BACKEND_HOSTS"
	BACKEND_MESSAGING_PORT string = "BACKEND_MESSAGING_PORT"
	BACKEND_SEARCH_PORT    string = "BACKEND_SEARCH_PORT"

	// defaults
	DEFAULT_LISTEN_PORT            string = ":42001"
	DEFAULT_STATS_HOST             string = "127.0.0.1"
	DEFAULT_STATS_PORT             string = "8125"
	DEFAULT_BACKEND_HOSTS          string = "127.0.0.1"
	DEFAULT_BACKEND_MESSAGING_PORT string = "4222"
	DEFAULT_BACKEND_SEARCH_PORT    string = "9200"
)

var _environment *Environment

type Environment struct {
	port                 string
	statsServiceHost     string
	statsServicePort     string
	statsServiceEnabled  bool
	backendServiceHosts  []string
	messagingServicePort string
	searchServicePort    string
}

func (e *Environment) ListenPort() string {
	return e.port
}

func (e *Environment) StatsServiceHost() string {
	return e.statsServiceHost
}

func (e *Environment) StatsServicePort() string {
	return e.statsServicePort
}

func (e *Environment) StatsServiceEnabled() bool {
	return e.statsServiceEnabled
}

func (e *Environment) BackendServiceHosts() []string {
	return e.backendServiceHosts
}

func (e *Environment) MessagingServicePort() string {
	return e.messagingServicePort
}

func (e *Environment) SearchServicePort() string {
	return e.searchServicePort
}

func (e *Environment) MessagingServiceUrls() []string {
	u := make([]string, len(e.backendServiceHosts))
	for i := range e.backendServiceHosts {
		u[i] = strings.Join([]string{"nats://", e.backendServiceHosts[i], ":", e.messagingServicePort}, "")
	}
	return u
}

func (e *Environment) SearchServiceUrl() string {
	return strings.Join([]string{"http://", e.backendServiceHosts[0], ":", e.searchServicePort,"/"}, "")
}

func GetEnvironment() *Environment {
	if _environment == nil {
		e := Environment{
			getEnvOrDefault(LISTEN_PORT, DEFAULT_LISTEN_PORT),
			getEnvOrDefault(STATS_HOST, DEFAULT_STATS_HOST),
			getEnvOrDefault(STATS_PORT, DEFAULT_STATS_PORT),
			getEnvOrDefaultBool(STATS_ENABLED, true),
			getEnvOrDefaultN(BACKEND_HOSTS, DEFAULT_BACKEND_HOSTS),
			getEnvOrDefault(BACKEND_MESSAGING_PORT, DEFAULT_BACKEND_MESSAGING_PORT),
			getEnvOrDefault(BACKEND_SEARCH_PORT, DEFAULT_BACKEND_SEARCH_PORT),
		}
		_environment = &e
	}
	return _environment
}

func getEnvOrDefault(env string, defaultValue string) string {
	envVar := os.Getenv(env)
	if envVar == "" {
		envVar = defaultValue
	}
	return envVar
}

func getEnvOrDefaultBool(env string, defaultValue bool) bool {
	envVar := os.Getenv(env)
	if envVar == "" {
		return defaultValue
	} else {
		b, err := strconv.ParseBool(envVar)
		if err != nil {
			return defaultValue
		}
		return b
	}
}

func getEnvOrDefaultN(env string, defaultValue string) []string {
	envVar := os.Getenv(env)
	if envVar == "" {
		envVar = defaultValue
	}

	// split into hosts/IPs
	h := strings.Split(envVar, ",")
	for i, s := range h {
		h[i] = strings.Trim(s, " ")
	}

	return h

}
