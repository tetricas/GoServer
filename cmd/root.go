package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"sort"
	"strings"
)

var (
	cfgFile   string
	container []string

	rootCmd = &cobra.Command{
		Use:   "testOne",
		Short: "Some test exercises to learn new libs.",
		Long:  "A long time ago in a galaxy far, far away...",
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Service initialization",
		Long:  "Init DB connection, sync time, check the weather...",
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set("innited", true)
			viper.WriteConfig()
		},
	}
	addCmd = &cobra.Command{
		Use:   "add [data]",
		Short: "Add data to container",
		Long:  "Add additional components to storing data type...",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !checkInitialization() {
				return
			}
			container = append(container, args...)
		},
	}
	removeCmd = &cobra.Command{
		Use:   "remove [data]",
		Short: "Remove data from container",
		Long:  "Remove additional components from storing data type...",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !checkInitialization() {
				return
			}
			idx := sort.Search(len(container), func(i int) bool {
				return container[i] == args[0]
			})
			container[idx] = container[len(container)-1]
			container[len(container)-1] = ""
			container = container[:len(container)-1]
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all container data",
		Long:  "List all additional data from storing data type...",
		Run: func(cmd *cobra.Command, args []string) {
			if !checkInitialization() {
				return
			}
			fmt.Println("Data: " + strings.Join(container, " "))
		},
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", rootCmd.Flags().Lookup("viper"))

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".testOne")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func checkInitialization() bool {
	if !viper.GetBool("innited") {
		fmt.Println("Go fuck yourself")
		return false
	}
	return true
}
