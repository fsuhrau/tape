// Copyright © 2020 Fabian Suhrau fabian.suhrau@me.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/fsuhrau/tape/repository"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add name url",
	Short: "add a new dependency",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return repository.ErrMissingParameter
		}
		repo, err := repository.Load()
		if err != nil {
			return err
		}
		if err := repo.Add(args[0], args[1]); err != nil {
			return err
		}

		return repo.Save()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
