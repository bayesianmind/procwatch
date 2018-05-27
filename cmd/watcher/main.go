package main

import (
	"fmt"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bayesianmind/procwatch"
)

var rootCmd = &cobra.Command{
	Use:   "watcher",
	Short: "watcher watches a process and starts it if idle",
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan bool, 1)
		w := watch.Watcher{
			CheckInterval: viper.GetDuration("chkint"),
			Args:          viper.GetStringSlice("args"),
			Command:       viper.GetString("cmd"),
			DedupeCmd:     viper.GetString("dedupecmd"),
			IdleInterval:  viper.GetDuration("idleint"),
			Cwd:           viper.GetString("cwd"),
		}

		fmt.Printf("Using settings %+v\n", w)

		w.Start()
		<-done
	},
}

func init() {
	// disable Cobra forcing run from cmd.exe on windows
	cobra.MousetrapHelpText = ""

	rootCmd.Flags().Duration("chkint", 30*time.Second, "How often to check idle time/processes")
	rootCmd.Flags().Duration("idleint", 1*time.Hour, "Idle time before process started if not already running")
	rootCmd.Flags().String("cmd", "", "Command to run")
	rootCmd.Flags().StringSlice("args", nil, "Command args")
	rootCmd.Flags().String("dedupecmd", "", "Command to look for in proc table to skip creating")
	rootCmd.Flags().String("cwd", "", "Working dir")

	viper.BindPFlag("chkint", rootCmd.Flags().Lookup("chkint"))
	viper.BindPFlag("idleint", rootCmd.Flags().Lookup("idleint"))
	viper.BindPFlag("cmd", rootCmd.Flags().Lookup("cmd"))
	viper.BindPFlag("args", rootCmd.Flags().Lookup("args"))
	viper.BindPFlag("dedupecmd", rootCmd.Flags().Lookup("dedupecmd"))
	viper.BindPFlag("cwd", rootCmd.Flags().Lookup("cwd"))

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("can't get homedir:", err)
		os.Exit(1)
	}

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName("watcher")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Can't read config: %v\n", err)
			os.Exit(1)
		}
	}

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
