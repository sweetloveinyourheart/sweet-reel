package config

import (
	"fmt"
	"slices"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ServerId            = "server.id"
	ServerReplicaNumber = "server.replica_number"
	ServerReplicaCount  = "server.replica_count"
	ServerNamespace     = "server.namespace"
)

var Registry = make(map[*pflag.Flag]Binding)
var securePaths = make(map[string]struct{})

func SecureFields(fields ...string) {
	for _, field := range fields {
		securePaths[field] = struct{}{}
	}
}

func GetSecureFields() []string {
	var fields []string
	for field := range securePaths {
		fields = append(fields, field)
	}
	return fields
}

type Binding struct {
	Flag       *pflag.Flag
	Aliases    []string
	Value      any
	Usage      string
	Env        []string
	HasDefault bool
}

func Bind(flag *pflag.Flag, name string, env ...string) {
	if flag == nil {
		panic(fmt.Sprintf("cannot bind pFlag '%s' no such flag was registered", name))
	}
	_ = viper.BindPFlag(name, flag)
	_ = viper.BindEnv(append([]string{name}, env...)...)
	registerBinding(flag, name, flag.Usage, nil, false, env...)
}

func BindWithDefault(flag *pflag.Flag, name string, defaultValue any, env ...string) {
	if flag == nil {
		panic(fmt.Sprintf("cannot bind pFlag '%s' no such flag was registered", name))
	}
	_ = viper.BindPFlag(name, flag)
	viper.SetDefault(name, defaultValue)
	_ = viper.BindEnv(append([]string{name}, env...)...)
	registerBinding(flag, name, flag.Usage, defaultValue, true, env...)
}

func Int64(command *cobra.Command, name string, flag string, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().Int(flag, 0, usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, 0, false, env...)
	}
	Bind(command.PersistentFlags().Lookup(flag), name, env...)
}

func Int64Default(command *cobra.Command, name string, flag string, defaultValue int64, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().Int64(flag, defaultValue, usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, defaultValue, true, env...)
	}
	BindWithDefault(command.PersistentFlags().Lookup(flag), name, defaultValue, env...)
}

func Int32(command *cobra.Command, name string, flag string, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().Int32(flag, 0, usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, 0, false, env...)
	}
	Bind(command.PersistentFlags().Lookup(flag), name, env...)
}

func Int32Default(command *cobra.Command, name string, flag string, defaultValue int32, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().Int32(flag, defaultValue, usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, defaultValue, true, env...)
	}
	BindWithDefault(command.PersistentFlags().Lookup(flag), name, defaultValue, env...)
}

func Bool(command *cobra.Command, name string, flag string, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().Bool(flag, false, usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, false, false, env...)
	}
	Bind(command.PersistentFlags().Lookup(flag), name, env...)
}

func BoolDefault(command *cobra.Command, name string, flag string, defaultValue bool, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().Bool(flag, defaultValue, usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, defaultValue, true, env...)
	}
	BindWithDefault(command.PersistentFlags().Lookup(flag), name, defaultValue, env...)
}

func String(command *cobra.Command, name string, flag string, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().String(flag, "", usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, "", false, env...)
	}
	Bind(command.PersistentFlags().Lookup(flag), name, env...)
}

func StringDefault(command *cobra.Command, name string, flag string, defaultValue string, usage string, env ...string) {
	if command.PersistentFlags().Lookup(flag) == nil {
		command.PersistentFlags().String(flag, defaultValue, usage)
		registerBinding(command.PersistentFlags().Lookup(flag), name, usage, defaultValue, true, env...)
	}
	BindWithDefault(command.PersistentFlags().Lookup(flag), name, defaultValue, env...)
}

func registerBinding(flag *pflag.Flag, name string, usage string, defaultValue any, hasDefault bool, env ...string) {
	if Registry == nil {
		Registry = make(map[*pflag.Flag]Binding)
	}
	if existing, ok := Registry[flag]; !ok {
		b := Binding{
			Flag:       flag,
			Aliases:    []string{name},
			Value:      defaultValue,
			Usage:      usage,
			Env:        env,
			HasDefault: hasDefault,
		}

		Registry[flag] = b
	} else {
		if !lo.Contains(existing.Aliases, name) {
			existing.Aliases = append(existing.Aliases, name)
		}
		if existing.Usage != usage {
			existing.Usage = usage
		}
		if existing.Value != defaultValue && defaultValue != nil {
			existing.Value = defaultValue
		}
		if existing.HasDefault != hasDefault {
			existing.HasDefault = hasDefault
		}
		if len(existing.Env) != len(env) {
			existing.Env = env
		}

		Registry[flag] = existing
	}
}

func ResolveRequireFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		lookup := Registry[f]
		if slices.ContainsFunc(lookup.Aliases, viper.IsSet) {
			_ = cmd.Flags().SetAnnotation(f.Name, cobra.BashCompOneRequiredFlag, []string{"false"})
			return
		}
	})
}
