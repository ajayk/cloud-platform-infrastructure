package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds the basic structure of test's YAML file
type Config struct {
	ClusterName             string                  `yaml:"clusterName"`
	Namespaces              map[string]K8SObjects   `yaml:"namespaces"`
	ExternalDNS             ExternalDNS             `yaml:"externalDNS"`
	NginxIngressController  NginxIngressController  `yaml:"nginxIngressController"`
	ModsecIngressController ModsecIngressController `yaml:"modsecIngressController"`
	FilesExist              []string                `yaml:"filesExist"`
}

// K8SObjects are kubernetes objects nested from namespaces, we need to check
// these resources are checked for its existence
type K8SObjects struct {
	Servicemonitors []string `yaml:"servicemonitors"`
	Daemonsets      []string `yaml:"daemonsets"`
	Services        []string `yaml:"services"`
	Secrets         []string `yaml:"secrets"`
}

// ParseConfigFile loads the test file supplied
func ParseConfigFile(f string) (*Config, error) {
	testsFilePath, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	t := Config{}

	err = yaml.Unmarshal(testsFilePath, &t)
	if err != nil {
		return nil, err
	}

	err = t.defaultsFromEnvs()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// defaultsFromEnvs process the mandatory fields in the config. If they are not set,
// it tries to load them from environment variables
func (c Config) defaultsFromEnvs() error {
	if c.ClusterName == "" {
		c.ClusterName = os.Getenv("CP_CLUSTER_NAME")
		if c.ClusterName == "" {
			return errors.New("cluster Name is mandatory - not found it neither in config file nor environment variable")
		}
	}

	return nil
}

// defaultsFromEnvs process the mandatory fields in the config. If they are not set,
// it tries to load them from environment variables
func (c *Config) GetExpectedDaemonSets() map[string][]string {
	r := make(map[string][]string)

	for ns, val := range c.Namespaces {
		var daemonSets []string

		daemonSets = append(daemonSets, val.Daemonsets...)

		if len(daemonSets) > 0 {
			r[ns] = daemonSets
		}
	}

	return r
}

// GetServiceMonitors process the mandatory fields in the config. If they are not set,
// it tries to load them from environment variables
func (c *Config) GetExpectedServiceMonitors() map[string][]string {
	r := make(map[string][]string)

	for ns, val := range c.Namespaces {
		var serviceMonitors []string

		serviceMonitors = append(serviceMonitors, val.Servicemonitors...)

		if len(serviceMonitors) > 0 {
			r[ns] = serviceMonitors
		}
	}

	return r
}

// GetExpectedServices returns a slice of all the services
// that are expected to be in the cluster.
func (c *Config) GetExpectedServices() map[string][]string {
	r := make(map[string][]string)

	for ns, val := range c.Namespaces {
		var services []string

		services = append(services, val.Services...)

		if len(services) > 0 {
			r[ns] = services
		}
	}

	return r
}
