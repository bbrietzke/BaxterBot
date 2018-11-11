package cmd

import (
	"strings"

	"github.com/bbrietzke/BaxterBot/pkg/store"
	"github.com/spf13/cobra"
)

func storeFlagParse(cmd *cobra.Command, args []string) []store.Argument {
	storeOptions := []store.Argument{}

	if cmd.Flag("repl").Changed {
		storeOptions = append(storeOptions, store.Port(cmd.Flag("repl").Value.String()))
	}

	if cmd.Flag("join").Changed {
		s := strings.TrimPrefix(cmd.Flag("join").Value.String(), "[\"")
		s = strings.TrimSuffix(s, "\"]")

		for _, a := range strings.Split(s, ",") {
			storeOptions = append(storeOptions, store.Join(a))
		}
	}

	if cmd.Flag("name").Changed {
		storeOptions = append(storeOptions, store.Name(cmd.Flag("name").Value.String()))
	}

	return storeOptions
}
