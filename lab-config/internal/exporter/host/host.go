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

package host

import (
	"fmt"
	"regexp"
	"strings"

	api "github.com/the78mole/jumpstarter-mono/lab-config/api/v1alpha1"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/config"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/exporter/ssh"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/exporter/template"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/templating"
)

type ExporterHostSyncer struct {
	cfg                  *config.Config
	tapplier             *templating.TemplateApplier
	serviceParametersMap map[string]template.ServiceParameters
	dryRun               bool
	debugConfigs         bool
	exporterFilter       *regexp.Regexp
}

func NewExporterHostSyncer(cfg *config.Config,
	tapplier *templating.TemplateApplier,
	serviceParametersMap map[string]template.ServiceParameters,
	dryRun, debugConfigs bool,
	exporterFilter *regexp.Regexp) *ExporterHostSyncer {
	return &ExporterHostSyncer{
		cfg:                  cfg,
		tapplier:             tapplier,
		serviceParametersMap: serviceParametersMap,
		dryRun:               dryRun,
		debugConfigs:         debugConfigs,
		exporterFilter:       exporterFilter,
	}
}

// SyncExporterHosts synchronizes exporter hosts via SSH
func (e *ExporterHostSyncer) SyncExporterHosts() error {
	fmt.Print("\n🔄 Syncing exporter hosts via SSH ===========================\n")
	for _, host := range e.cfg.Loaded.ExporterHosts {
		exporterInstances := e.cfg.Loaded.GetExporterInstancesByExporterHost(host.Name)
		// Apply filter to exporter instances if provided
		if e.exporterFilter != nil {
			filteredInstances := []*api.ExporterInstance{}
			for _, exporterInstance := range exporterInstances {
				if e.exporterFilter.MatchString(exporterInstance.Name) {
					filteredInstances = append(filteredInstances, exporterInstance)
				}
			}
			exporterInstances = filteredInstances
		}
		// Skip the host if no exporter instances match the filter
		if len(exporterInstances) == 0 {
			continue
		}

		hostCopy := host.DeepCopy()
		err := e.tapplier.Apply(hostCopy)
		if err != nil {
			return fmt.Errorf("error applying template for %s: %w", host.Name, err)
		}
		fmt.Printf("\n💻  Exporter host: %s\n", hostCopy.Spec.Addresses[0])
		hostSsh, err := ssh.NewSSHHostManager(hostCopy)
		if err != nil {
			return fmt.Errorf("error creating SSH host manager for %s: %w", host.Name, err)
		}
		status, err := hostSsh.Status()
		if err != nil {
			return fmt.Errorf("error getting status for %s: %w", host.Name, err)
		}
		fmt.Printf("    Connection status: %s\n", status)

		for _, exporterInstance := range exporterInstances {
			fmt.Printf("    📟 Exporter instance: %s\n", exporterInstance.Name)
			errName := "ExporterInstance:" + exporterInstance.Name
			et, err := template.NewExporterInstanceTemplater(e.cfg, exporterInstance)
			serviceParameters, ok := e.serviceParametersMap[exporterInstance.Name]
			if !ok {
				return fmt.Errorf("service parameters not found for %s", exporterInstance.Name)
			}
			et.SetServiceParameters(serviceParameters)
			if err != nil {
				return fmt.Errorf("error creating ExporterInstanceTemplater for %s : %w", errName, err)
			}

			_, err = et.RenderTemplateLabels()
			if err != nil {
				return fmt.Errorf("error creating ExporterInstanceTemplater for %s : %w", errName, err)
			}
			tcfg, err := et.RenderTemplateConfig()
			if err != nil {
				return fmt.Errorf("error rendering template config for %s : %w", errName, err)
			}

			if e.debugConfigs {
				fmt.Printf("--- 📄 Config Template %s\n", strings.Repeat("─", 40))
				fmt.Printf("%s\n", tcfg.Spec.ConfigTemplate)
				if tcfg.Spec.SystemdContainerTemplate != "" {
					fmt.Printf("  - ⚙️  Systemd Container Template %s\n", strings.Repeat("─", 30))
					fmt.Printf("%s\n", tcfg.Spec.SystemdContainerTemplate)
				}
				if tcfg.Spec.SystemdServiceTemplate != "" {
					fmt.Printf("  - 🔧 Systemd Service Template %s\n", strings.Repeat("─", 31))
					fmt.Printf("%s\n", tcfg.Spec.SystemdServiceTemplate)
				}
				fmt.Println(strings.Repeat("─", 60))
			}

			if err := hostSsh.Apply(tcfg, e.dryRun); err != nil {
				return fmt.Errorf("error applying exporter config for %s: %w", exporterInstance.Name, err)
			}

		}
	}
	return nil
}
