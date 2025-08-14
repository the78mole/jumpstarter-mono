/*
Copyright 2025. The Jumpstarter Authors.

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

package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/the78mole/jumpstarter-mono/lab-config/internal/config"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/config_lint"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/exporter/host"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/exporter/template"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/instance"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/templating"
)

var applyCmd = &cobra.Command{
	Use:   "apply [config-file]",
	Short: "Apply configuration changes",
	Long:  `Apply configuration changes to the jumpstarter controllers. Use --dry-run to verify changes before applying.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		prune, _ := cmd.Flags().GetBool("prune")
		vaultPassFile, _ := cmd.Flags().GetString("vault-password-file")
		debugConfigs, _ := cmd.Flags().GetBool("debug-configs")
		filterClients, _ := cmd.Flags().GetString("filter-clients")
		filterExporters, _ := cmd.Flags().GetString("filter-exporters")

		// Determine config file path
		configFilePath := "jumpstarter-lab.yaml" // default
		if len(args) > 0 {
			configFilePath = args[0]
		}

		// Load the configuration file
		cfg, err := config.LoadConfig(configFilePath, vaultPassFile)
		if err != nil {
			return fmt.Errorf("error loading config file %s: %w", configFilePath, err)
		}

		config_lint.Validate(cfg)
		tapplier, err := templating.NewTemplateApplier(cfg, nil)
		if err != nil {
			return fmt.Errorf("error creating template applier %w", err)
		}

		// Compile client filter regexp if provided
		var clientFilter *regexp.Regexp
		if filterClients != "" {
			clientFilter, err = regexp.Compile(filterClients)
			if err != nil {
				return fmt.Errorf("invalid client filter regexp '%s': %w", filterClients, err)
			}
		}

		// Compile exporter filter regexp if provided
		var exporterFilter *regexp.Regexp
		if filterExporters != "" {
			exporterFilter, err = regexp.Compile(filterExporters)
			if err != nil {
				return fmt.Errorf("invalid exporter filter regexp '%s': %w", filterExporters, err)
			}
		}

		if dryRun {
			fmt.Println("Dry run: Would apply changes to:")
			fmt.Println()
		} else {
			fmt.Println("Applying changes:")
			fmt.Println()
		}
		var serviceParametersMap map[string]template.ServiceParameters

		for _, inst := range cfg.Loaded.JumpstarterInstances {
			instanceCopy := inst.DeepCopy()
			err = tapplier.Apply(instanceCopy)
			if err != nil {
				return fmt.Errorf("error applying template for %s: %w", inst.Name, err)
			}
			instanceClient, err := instance.NewInstance(instanceCopy, instanceCopy.Spec.Kubeconfig, dryRun, prune)
			if err != nil {
				return fmt.Errorf("error creating instance for %s: %w", inst.Name, err)
			}

			err = instanceClient.SyncClients(context.Background(), cfg, clientFilter)
			if err != nil {
				return fmt.Errorf("error syncing clients for %s: %w", inst.Name, err)
			}

			serviceParametersMap, err = instanceClient.SyncExporters(context.Background(), cfg, exporterFilter)
			if err != nil {
				return fmt.Errorf("error syncing exporters for %s: %w", inst.Name, err)
			}
		}

		exporterHostSyncer := host.NewExporterHostSyncer(cfg, tapplier, serviceParametersMap, dryRun, debugConfigs,
			exporterFilter)

		err = exporterHostSyncer.SyncExporterHosts()
		if err != nil {
			return fmt.Errorf("error syncing exporter hosts: %w", err)
		}

		return nil
	},
}

func init() {
	// Add flags to apply command
	applyCmd.Flags().Bool("dry-run", false, "Show what would be applied without making changes")
	applyCmd.Flags().Bool("prune", false, "Delete resources that are no longer defined in configuration")
	applyCmd.Flags().String("vault-password-file", "", "Path to the vault password file for decrypting variables")
	applyCmd.Flags().Bool("debug-configs", false, "Show debug configs")
	applyCmd.Flags().String("filter-clients", "", "Regexp pattern to filter clients by name")
	applyCmd.Flags().String("filter-exporters", "", "Regexp pattern to filter exporters by name")

	rootCmd.AddCommand(applyCmd)
}
