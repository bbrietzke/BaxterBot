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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// vegetaCmd represents the vegeta command
var vegetaCmd = &cobra.Command{
	Use:   "vegeta",
	Short: "A brief description of your command",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		host := cmd.Flag("addr").Value.String()
		fileNameKey := time.Now().UnixNano()
		rand.Seed(fileNameKey)
		c := rand.Intn(500) + 50
		jsonFile := fmt.Sprintf("/tmp/%d.json", fileNameKey)
		data := map[string]string{
			"Value": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			"Key":   "1",
		}
		b, _ := json.Marshal(data)
		ioutil.WriteFile(jsonFile, b, 0644)

		keys := []string{"january", "feburary", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december", "octember"}

		statements := []string{
			"POST %s/api/%s\nContent-Type: application/json\n@" + jsonFile + "\n",
			"GET %s/api/%s\nContent-Type: application/json\n",
			"POST %s/api/%s\nContent-Type: application/json\n@" + jsonFile + "\n",
			"GET %s/api/%s\nContent-Type: application/json\n",
			"DELETE %s/api/%s\nContent-Type: application/json\n",
			"GET %s/api/%s\nContent-Type: application/json\n",
		}

		for i := 0; i < c; i++ {
			key := keys[rand.Intn(len(keys))]
			statement := statements[rand.Intn(len(statements))]
			fmt.Printf(statement, host, key)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(vegetaCmd)
	vegetaCmd.Flags().StringP("addr", "a", "http://localhost:8080/", "http://hostname:port/")
}
