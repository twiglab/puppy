package serv

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twiglab/puppy/gbot"
	"gopkg.in/yaml.v3"
)

var ServCmd = &cobra.Command{
	Use:   "serv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}
var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	ServCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func printConf(conf AppConf) {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	enc.Encode(conf)
	enc.Close()
}

func run() {
	var conf AppConf
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	printConf(conf)

	exec := buildLocalExec(conf)
	exec.Init()

	app := buildGBotApp(conf)
	bot := buildGBot(conf)
	exec.RegJob(app)
	exec.RegJob(bot)

	h := buildHandle(conf, app)

	mux := chi.NewMux()
	mux.Mount("/", exec)
	mux.Handle("/gbot", h)
	mux.Handle("/mcp", gbot.McpHandle())

	if err := http.ListenAndServe(exec.LocalAddr, mux); err != nil {
		log.Fatal(err)
	}
}
