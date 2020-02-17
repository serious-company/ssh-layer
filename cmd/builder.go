package main

import (
	"fmt"
	"io/ioutil"
	"os"

	builder "github.com/serious-company/ssh-layer/pkg/builder"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type builderFlags struct {
	sshKey     string `mapstructure:"ssh-key"`
	sshKeyPath string `mapstructure:"ssh-key-path"`
}

func (f *builderFlags) getKeyPath(homePath string) string {
	if f.sshKeyPath != "" {
		return f.sshKeyPath
	}

	return homePath + "/.ssh/" + f.sshKey
}

var builderFlagsInstance builderFlags

var builderCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a docker image with the file id_rsa inside",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(builderCmd)

	builderCmd.Flags().StringVar(&builderFlagsInstance.sshKey, "ssh-key", "id_rsa", "SSH Key Name")
	builderCmd.Flags().StringVar(&builderFlagsInstance.sshKeyPath, "ssh-key-path", "", "SSH Key path")
}

func run() (err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := builderFlagsInstance.getKeyPath(home)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// @TODO: Turn logs on with verbose flag
	builder.SetLogger(zap.New(func(o *zap.Options) {
	}))

	b := builder.Config{
		IDRsa:     string(file),
		ImageName: "ssh-layer",
	}

	return b.Build()
}
