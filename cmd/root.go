package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	moduleShortName        = "rke"
	configFileName         = "rke-config.json"
	stateFileName          = "state.json"
	inventoryDirectoryName = "inventory"
	inventoryFileName      = "hosts.json"
	envDirectoryName       = "env"
	sshKeyFileName         = "ssh_key"
	cmdlineFileName        = "cmdline"

	defaultSharedDirectoryPath    = "/shared"
	defaultResourcesDirectoryPath = "/resources"
)

var (
	enableDebug       bool
	ansibleDebugLevel int

	Version string

	SharedDirectoryPath    string
	ResourcesDirectoryPath string

	logger zerolog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "m-rke",
	Long: `Module responsible for installation of Kubernetes cluster with RKE`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PersistentPreRun")

		err := viper.BindPFlags(cmd.PersistentFlags())
		if err != nil {
			logger.Fatal().Err(err).Msg("initialization error occurred")
		}

		SharedDirectoryPath = viper.GetString("shared")
		ResourcesDirectoryPath = viper.GetString("resources")

		logger.Trace().Msgf("original ansibleDebugLevel: %d", ansibleDebugLevel)
		if ansibleDebugLevel > 6 {
			ansibleDebugLevel = 6
		}
		if ansibleDebugLevel < 0 {
			ansibleDebugLevel = 0
		}
	},
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Caller().Timestamp().Logger()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&enableDebug, "debug", "d", false, "print debug information")

	rootCmd.PersistentFlags().String("shared", defaultSharedDirectoryPath, "Shared directory path.")
	_ = rootCmd.MarkPersistentFlagDirname("shared")
	rootCmd.PersistentFlags().String("resources", defaultResourcesDirectoryPath, "Resources directory path.")
	_ = rootCmd.MarkPersistentFlagDirname("resources")
	rootCmd.PersistentFlags().IntVarP(&ansibleDebugLevel, "ansible_debug_level", "a", 0, "set ansible debug level 0-6")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if enableDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	viper.AutomaticEnv() // read in environment variables that match
}
