package cmd

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
	"testOne/internal"
)

var (
	cfgFile string
	port    string

	rootCmd = &cobra.Command{
		Use:   "testOne",
		Short: "Some test exercises to learn new libs.",
		Long:  "A long time ago in a galaxy far, far away...",
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Service initialization",
		Long:  "Init DB connection, sync time, check the weather...",
		Run: func(cmd *cobra.Command, args []string) {
			internal.Start(port)
		},
	}
	portCmd = &cobra.Command{
		Use:   "port",
		Short: "Set port",
		Long:  "Set port where server will listen",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			port = args[0]
			i, err := strconv.Atoi(args[0])
			if err != nil || i > 9999 || i < 1 {
				fmt.Println("Incorrect input: port")
				log.Println(err)
				return
			}
			viper.Set("port", port)
			_ = viper.WriteConfig()
		},
	}
	setAdminCmd = &cobra.Command{
		Use:   "admin",
		Short: "Set admin",
		Long:  "Set administrator",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			password := args[1]
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}

			user := internal.UserInternal{
				Name:    username,
				Email:   username,
				Secret:  string(hashedPassword),
				IsAdmin: true,
			}
			internal.AddUserToDB(&user)
		},
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	internal.ConnectToDB()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	startCmd.PersistentFlags().StringVarP(&port, "port", "p", "1232", "port where server will listen")
	_ = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(portCmd)
	rootCmd.AddCommand(setAdminCmd)
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
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
