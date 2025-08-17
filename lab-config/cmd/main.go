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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	metav1alpha1 "github.com/the78mole/jumpstarter-mono/lab-config/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// Version information set by build
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(metav1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

var rootCmd = &cobra.Command{
	Use:   "jumpstarter-lab-config",
	Short: "A tool for managing labs and environments with jumpstarter controllers",
	Long: `A gitops/configuration system and tool for managing labs and environments
	with jumpstarter controllers, exporter hosts (sidekicks), exporters, and clients
	in enterprise environments.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "version" {
			fmt.Printf("jumpstarter-lab-config version %s (commit %s, built %s)\n", Version, Commit, BuildTime)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("jumpstarter-lab-config version %s (commit %s, built %s)\n", Version, Commit, BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// nolint:gocyclo
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
