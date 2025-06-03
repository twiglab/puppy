package serv

import (
	"os"

	"github.com/imroc/req/v3"
	"github.com/it512/xxl-job-exec"
	"github.com/spf13/cobra"
	"github.com/twiglab/puppy"
	"github.com/twiglab/puppy/gbot"
	"github.com/xen0n/go-workwx/v2"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c()
	},
}

func init() {
}

type XxlJobConf struct {
	ServerAddr  string `yaml:"server-addr" mapstructure:"server-addr"`
	AccessToken string `yaml:"access-token" mapstructure:"access-token"`
	LocalIp     string `yaml:"local-ip" mapstructure:"local-ip"`
	LocalPort   string `yaml:"local-port" mapstructure:"local-port"`
	RegistryKey string `yaml:"registry-key" mapstructure:"registry-key"`
}

type WxBotConf struct {
	BotID string `yaml:"bot-id" mapstructure:"bot-id"`
}

type AmapWeatherConf struct {
	Key string `yaml:"key" mapstructure:"key"`
}

type GBotConf struct {
	Name   string `yaml:"name" mapstructure:"name"`
	AdCode string `yaml:"adcode" mapstructure:"adcode"`
}

type LoggerConf struct {
	Level   string `yaml:"level" mapstructure:"level"`
	LogFile string `yaml:"log-file" mapstructure:"log-file"`
}

type DcpConf struct {
	URL string `yaml:"url" mapstructure:"url"`
}

type AppConf struct {
	ID              string
	LoggerConf      LoggerConf      `yaml:"log" mapstructure:"log"`
	XxlJobConf      XxlJobConf      `yaml:"xxl" mapstructure:"xxl"`
	AmapWeatherConf AmapWeatherConf `yaml:"weather" mapstructure:"weather"`
	GBotConf        GBotConf        `yaml:"gbot" mapstructure:"gbot"`
	WxBotConf       WxBotConf       `yaml:"wxbot" mapstructure:"wxbot"`
	DcpConf         DcpConf         `yaml:"dcp" mapstructure:"dcp"`
}

func buildGBot(conf AppConf) *gbot.GBotJob {
	c := req.C().EnableInsecureSkipVerify()
	whc := workwx.NewWebhookClient(conf.WxBotConf.BotID, workwx.WithHTTPClient(c.GetClient()))
	dcp := puppy.NewDcpServ(conf.DcpConf.URL, c)
	weather := puppy.NewAmapWeather(conf.AmapWeatherConf.Key, c)

	return &gbot.GBotJob{
		JobName: conf.GBotConf.Name,
		AdCode:  conf.GBotConf.AdCode,
		Weater:  weather,
		Dcp:     dcp,
		MsgBot:  whc,
		Tpl:     gbot.GBotTemplate(),
	}
}

func buildLocalExec(conf AppConf) *puppy.LocalExec {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(conf.XxlJobConf.ServerAddr),
		xxl.AccessToken(conf.XxlJobConf.AccessToken),
		xxl.ExecutorIp(conf.XxlJobConf.LocalIp),
		xxl.ExecutorPort(conf.XxlJobConf.LocalPort),
		xxl.RegistryKey(conf.XxlJobConf.RegistryKey),
	)
	return puppy.NewLocalExec(":"+conf.XxlJobConf.LocalPort, exec)
}

func c() {
	conf := AppConf{
		ID: "gbot",
		LoggerConf: LoggerConf{
			Level:   "INFO",
			LogFile: "console",
		},
		XxlJobConf: XxlJobConf{
			ServerAddr:  "http://127.0.0.1:8080/xxl-job-admin",
			AccessToken: "token",
			LocalIp:     "127.0.0.1",
			LocalPort:   "10009",
			RegistryKey: "puppy",
		},
		AmapWeatherConf: AmapWeatherConf{
			Key: "f7e3dbe28def5f2e1028d2ae007e91a7",
		},
		WxBotConf: WxBotConf{
			BotID: "101c622f-b942-42c7-ab92-27bc49c031a8",
		},
		GBotConf: GBotConf{
			Name:   "gbot",
			AdCode: "320100",
		},
		DcpConf: DcpConf{
			URL: "https://127.0.0.1:10005/jsonrpc",
		},
	}
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	enc.Encode(&conf)
}
