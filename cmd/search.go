/*
Copyright © 2022 Isan Rivkin isanrivkin@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"path/filepath"

	search "github.com/isan-rivkin/search-unified-recusive-fast/lib/search/vaultsearch"
	"github.com/isan-rivkin/search-unified-recusive-fast/lib/vault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	query        *string
	parallel     *int
	mount        *string
	prefix       *string
	outputWebURL *bool
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "pattern match again storage in valut",
	Long: `
$surf search -q aws -m backend-secrets/prod  -t 15
	`,
	Run: func(cmd *cobra.Command, args []string) {

		mount := getEnvOrOverride(mount, "default_mount")
		prefix := getEnvOrOverride(prefix, "default_prefix")

		basePath := filepath.Join(*mount, *prefix)

		client := runDefaultAuth()

		log.WithFields(log.Fields{
			"address":   client.GetVaultAddr(),
			"base_path": basePath,
			"query":     *query,
		}).Info("starting search")

		m := search.NewDefaultRegexMatcher()
		s := search.NewRecursiveSearcher[search.VC, search.Matcher](client, m)
		output, err := s.Search(search.NewSearchInput(*query, basePath, *parallel))

		if err != nil {
			panic(err)
		}

		if output != nil {
			for _, i := range output.Matches {
				path := i.GetFullPath()
				if *outputWebURL {
					fmt.Println(vault.PathToWebURL(client.GetVaultAddr(), path))
				} else {
					fmt.Println(path)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")
	query = searchCmd.PersistentFlags().StringP("query", "q", "", "search query regex supported")
	mount = searchCmd.PersistentFlags().StringP("mount", "m", "", "mount to start the search at the root")
	prefix = searchCmd.PersistentFlags().StringP("prefix", "p", "", "$mount/prefix inside the mount to search in")
	parallel = searchCmd.PersistentFlags().IntP("threads", "t", 10, "parallel search number")

	outputWebURL = searchCmd.PersistentFlags().Bool("output-url", true, "defaullt output is web urls to click on and go to the browser UI")

	//searchCmd.MarkPersistentFlagRequired("query")
	//searchCmd.MarkPersistentFlagRequired("mount")
}