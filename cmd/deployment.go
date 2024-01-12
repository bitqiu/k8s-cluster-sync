package cmd

import (
	"k8s-cluster-sync/config"
	"k8s-cluster-sync/logic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

// deploymentCmd represents the deployment command
var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "deployment sync",
	Long:  `deployment sync`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceCfg, err := clientcmd.BuildConfigFromFlags("", config.SourceCfgFile)
		if err != nil {
			return err
		}
		sourceCli, err := kubernetes.NewForConfig(sourceCfg)
		if err != nil {
			return err
		}

		targetCfg, err := clientcmd.BuildConfigFromFlags("", config.TargetCfgFile)
		if err != nil {
			return err
		}
		targetCli, err := kubernetes.NewForConfig(targetCfg)
		if err != nil {
			return err
		}

		return logic.SyncDeployments(sourceCli, targetCli, config.Namespace)
	},
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
}
