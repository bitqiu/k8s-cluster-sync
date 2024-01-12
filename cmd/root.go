package cmd

import (
	"fmt"
	"k8s-cluster-sync/config"
	"k8s-cluster-sync/pkg"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster sync",
	Long:  `cluster sync`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", `使用 `+pkg.Red(`-h`)+` 查看参数`)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.SourceCfgFile, "sourceCfgFile", "s", "", "source cluster kubeconfig config file path $HOME/.kube/source")
	rootCmd.PersistentFlags().StringVarP(&config.TargetCfgFile, "targetCfgFile", "t", "", "target cluster kubeconfig config file path $HOME/.kube/target")
	rootCmd.PersistentFlags().StringVarP(&config.Namespace, "namespace", "n", "default", "cluster namespace")
}
