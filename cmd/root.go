package cmd

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	cfgFile    string
	db         *sqlx.DB
	dataInsert = `INSERT INTO container (data) VALUES (?)`
	dataDelete = `DELETE FROM container WHERE data = (?)`
	dataSelect = `SELECT data FROM container`

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
			_ = viper.WriteConfig()
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
			for _, element := range args {
				tx := db.MustBegin()
				tx.MustExec(dataInsert, element)
				_ = tx.Commit()
			}
		},
	}
	removeCmd = &cobra.Command{
		Use:   "remove [data]",
		Short: "Remove data from container",
		Long:  "Remove additional components from storing data type...",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !checkInitialization() {
				return
			}
			for _, element := range args {
				tx := db.MustBegin()
				tx.MustExec(dataDelete, element)
				_ = tx.Commit()
			}
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
			var container []string
			err := db.Select(&container, dataSelect)
			if err != nil {
				log.Println(err)
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
	connectToDB()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().Bool("viper", true, "use Viper for configuration")
	_ = viper.BindPFlag("useViper", rootCmd.Flags().Lookup("viper"))

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
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func checkInitialization() bool {
	if !viper.GetBool("innited") {
		fmt.Println("Go fuck yourself")
		return false
	}
	return true
}

func connectToDB() {
	db, _ = sqlx.Connect("sqlite3", "testOne.db")
	path := filepath.Join("db", "create_db.sql")
	file, _ := ioutil.ReadFile(path)
	schema := string(file)
	db.MustExec(schema)
}
