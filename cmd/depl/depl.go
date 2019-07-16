package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/depl/pkg/formatter"
	"github.com/softleader/depl/pkg/release"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var (
	// 在包版時會動態指定 version 及 commit
	version, commit string
	metadata        *release.Metadata

	// global flags
	offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	token      = os.Getenv("SL_TOKEN")
)

func main() {
	cobra.OnInitialize(
		initMetadata,
	)
	if err := newRootCmd(os.Args[1:]).Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depl",
		Short: "deploy application to SoftLeader docker swarm ecosystem",
		Long:  "Depl is a command line tool for deploy application to SoftLeader docker swarm ecosystem",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// remove the check if the plugin can run in offline mode
			if offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			logrus.SetOutput(cmd.OutOrStdout())
			logrus.SetFormatter(&formatter.PlainFormatter{})

			// use os.LookupEnv to retrieve the specific value of the environment (e.g. os.LookupEnv("SL_TOKEN"))
			for _, env := range os.Environ() {
				if strings.HasPrefix(env, "SL_") {
					logrus.Println(env)
				}
			}
			return nil
		},
	}

	cmd.AddCommand(
		newVersionCmd(),
	)

	f := cmd.PersistentFlags()
	f.BoolVar(&offline, "offline", offline, "work offline, Overrides $SL_OFFLINE")
	f.BoolVarP(&verbose, "verbose", "v", verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVar(&token, "token", token, "github access token. Overrides $SL_TOKEN")
	f.Parse(args)

	return cmd
}

// initMetadata 準備 app 的 release 資訊
func initMetadata() {
	metadata = release.NewMetadata(version, commit)
}