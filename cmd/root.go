// Copyright © 2018 Brian Brietzke
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bbrietzke/BaxterBot/pkg/swarm"

	"github.com/bbrietzke/BaxterBot/pkg/web"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile          string
	terminateChannel chan error
)

var rootCmd = &cobra.Command{
	Use:   "BaxterBot",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		terminateChannel = make(chan error, 5)
		webOptions := []web.Option{}
		swarmOptions := []swarm.Option{}

		if cmd.Flag("swarm").Changed {
			swarmOptions = append(swarmOptions, swarm.Port(cmd.Flag("swarm").Value.String()))
		}

		if cmd.Flag("join").Changed {
			swarmOptions = append(swarmOptions, swarm.Join(cmd.Flag("join").Value.String()))
		}

		if cmd.Flag("name").Changed {
			swarmOptions = append(swarmOptions, swarm.Name(cmd.Flag("name").Value.String()))
		}

		if cmd.Flag("http").Changed {
			v := cmd.Flag("http").Value.String()
			webOptions = append(webOptions, web.Port(v))
			swarmOptions = append(swarmOptions, swarm.HTTP(v))
		}

		if cmd.Flag("rps").Changed {
			v, _ := strconv.ParseInt(cmd.Flag("rps").Value.String(), 10, 64)
			webOptions = append(webOptions, web.RequestsPerSecond(v))
		}

		if cmd.Flag("wait").Changed {
			if v, _ := strconv.ParseBool(cmd.Flag("wait").Value.String()); v {
				webOptions = append(webOptions, web.Wait())
			}
		}

		go delegateSwarmStart(terminateChannel, swarmOptions...)
		go delegateWebStart(terminateChannel, webOptions...)

		select {
		case err := <-terminateChannel:
			fmt.Println(err)
		}
	},
}

func delegateWebStart(stop chan error, options ...web.Option) {
	stop <- web.Start(options...)
}

func delegateSwarmStart(stop chan error, options ...swarm.Option) {
	if err := swarm.Start(options...); err != nil {
		stop <- err
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.baxter_bot.yaml)")
	rootCmd.PersistentFlags().Bool("wait", false, "use wait protocol for rate limiting")
	rootCmd.PersistentFlags().Int64("rps", 0, "requests per second")
	rootCmd.PersistentFlags().String("swarm", ":21000", "port to host the swarm on")
	rootCmd.PersistentFlags().String("http", ":8080", "port to host the http server on")
	rootCmd.PersistentFlags().String("name", "", "name to take as part of the swarm")
	rootCmd.PersistentFlags().String("join", "", "swarm host to join to")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".baxter_bot")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
