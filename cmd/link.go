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

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "link all current dependencies into .bin/",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := repository.Load()
		if err != nil {
			return err
		}

		return repo.Link()
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
