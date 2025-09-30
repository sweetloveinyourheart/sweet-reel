package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var ServiceRegistry = make(map[*cobra.Command]Service)

type Service struct {
	Command          *cobra.Command `json:"-" yaml:"-"`
	Ports            []ServicePort  `json:"ports" yaml:"ports"`
	HealthCheckPorts []ServicePort  `json:"healthCheckPorts" yaml:"healthCheckPorts"`
	DefaultDatabase  string         `json:"defaultDatabaseName" yaml:"defaultDatabaseName"`
}

type ServicePort struct {
	Name          string   `json:"name" yaml:"name"`
	WireProtocol  string   `json:"wireProtocol" yaml:"wireProtocol"`
	Protocol      string   `json:"protocol" yaml:"protocol"`
	Public        bool     `json:"public" yaml:"public"`
	DefaultRoutes []string `json:"defaultRoutes" yaml:"defaultRoutes"`

	FlagName     string   `json:"flagName" yaml:"flagName"`
	FlagUsage    string   `json:"flagUsage" yaml:"flagUsage"`
	ConfigName   string   `json:"configName" yaml:"configName"`
	DefaultValue int      `json:"defaultValue" yaml:"defaultValue"`
	Env          []string `json:"env" yaml:"env"`
}

func RegisterService(command *cobra.Command, service Service) {
	ServiceRegistry[command] = service
}

func GetService(command *cobra.Command) (Service, bool) {
	service, ok := ServiceRegistry[command]
	return service, ok
}

func GetServices() map[*cobra.Command]Service {
	return ServiceRegistry
}

func AddServicePort(command *cobra.Command, flag *pflag.Flag, port ServicePort) {
	if command == nil {
		panic(errors.WithStack(errors.New("command is nil")))
	}
	if flag == nil {
		panic(errors.WithStack(fmt.Errorf("flag is nil for command %s", command.Name())))
	}

	service, _ := GetService(command)

	def, _ := strconv.ParseInt(flag.Value.String(), 10, 64)
	port.DefaultValue = int(def)
	port.FlagName = flag.Name
	port.FlagUsage = flag.Usage
	binding, ok := Registry[flag]
	if ok {
		port.ConfigName = binding.Aliases[0]
		port.Env = binding.Env
	}
	service.Ports = append(service.Ports, port)
	RegisterService(command, service)
}

func AddDefaultServicePorts(cmd *cobra.Command, rootCmd *cobra.Command) {
	addServiceHealthCheckPort(cmd, rootCmd.PersistentFlags().Lookup("healthcheck-port"), ServicePort{
		Name:         "healthcheck-grpc",
		WireProtocol: "tcp",
		Protocol:     "grpc",
	})
	addServiceHealthCheckPort(cmd, rootCmd.PersistentFlags().Lookup("healthcheck-web-port"), ServicePort{
		Name:          "healthcheck-web",
		WireProtocol:  "tcp",
		Protocol:      "http",
		DefaultRoutes: []string{"/healthz", "/readyz"},
	})
}

func AddDefaultDatabase(cmd *cobra.Command, dbName string) {
	if cmd == nil {
		panic(errors.WithStack(errors.New("command is nil")))
	}
	service, _ := GetService(cmd)
	dbName = strings.ReplaceAll(dbName, "-", "_")
	dbName = strings.ReplaceAll(dbName, ".", "_")
	service.DefaultDatabase = dbName
	RegisterService(cmd, service)
}

func addServiceHealthCheckPort(command *cobra.Command, flag *pflag.Flag, port ServicePort) {
	if command == nil {
		panic(errors.WithStack(errors.New("command is nil")))
	}
	if flag == nil {
		panic(errors.WithStack(fmt.Errorf("flag is nil for command %s", command.Name())))
	}

	service, _ := GetService(command)

	def, _ := strconv.ParseInt(flag.Value.String(), 10, 64)
	port.DefaultValue = int(def)
	port.FlagName = flag.Name
	port.FlagUsage = flag.Usage
	binding, ok := Registry[flag]
	if ok {
		port.ConfigName = binding.Aliases[0]
		port.Env = binding.Env
	}
	service.HealthCheckPorts = append(service.HealthCheckPorts, port)
	RegisterService(command, service)
}
