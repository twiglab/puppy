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
	p := gbot.JobParam{
		Date:  "now",
		Proj:  "项目名称",
		Entry: []string{"abcd", "1234"},
		Areas: []gbot.Area{
			{ID: "A", Name: "BF1", Cameras: []string{"5678", "wxyz"}},
			{ID: "B", Name: "F1", Cameras: []string{"4567", "hijk"}},
			{ID: "C", Name: "F2", Cameras: []string{"7890", "opqr"}},
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(p)
}
