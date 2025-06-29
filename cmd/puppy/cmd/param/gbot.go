package param

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/twiglab/puppy/gbot"
)

var GbotParamCmd = &cobra.Command{
	Use:   "gbot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		gbotParam()
	},
}

func init() {
	ParamCmd.AddCommand(GbotParamCmd)
}

func gbotParam() {
	p := gbot.JobParam{}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(p)
}
