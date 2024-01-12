/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"k8s-cluster-sync/config"
	"k8s-cluster-sync/logic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "service sync",
	Long:  `service sync`,
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
		return logic.SyncService(sourceCli, targetCli, config.Namespace)
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
