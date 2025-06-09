package serv

import (
	"log"
	"net/http"
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

type AmapWeatherConf struct {
	Key string `yaml:"key" mapstructure:"key"`
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
	WxAppConf       WxAppConf       `yaml:"wxapp" mapstructure:"wxapp"`
	DcpConf         DcpConf         `yaml:"dcp" mapstructure:"dcp"`
}

type WxAppConf struct {
	CorpID           string           `yaml:"corp-id" mapstructure:"corp-id"`
	CorpSecret       string           `yaml:"corp-secret" mapstructure:"corp-secret"`
	AgentID          int64            `yaml:"agent-id" mapstructure:"agent-id"`
	WxAppReceiveConf WxAppReceiveConf `yaml:"receive" mapstructure:"receive"`
}

type WxAppReceiveConf struct {
	Token  string `yaml:"token" mapstructure:"token"`
	AESKey string `yaml:"aes-key" mapstructure:"aes-key"`
}

func buildGBotApp(conf AppConf) *gbot.GBotApp {
	c := req.C().EnableInsecureSkipVerify()
	dcp := puppy.NewDcpServ(conf.DcpConf.URL, c)
	weather := puppy.NewAmapWeather(conf.AmapWeatherConf.Key, c)

	wx := workwx.New(conf.WxAppConf.CorpID)
	app := wx.WithApp(conf.WxAppConf.CorpSecret, conf.WxAppConf.AgentID)
	app.SpawnAccessTokenRefresher()

	return &gbot.GBotApp{
		Weater: weather,
		Dcp:    dcp,
		Tpl:    gbot.GBotTemplate(),

		App: app,

		AI:puppy.NewAI("http://119.8.32.64:8000"),
	}
}

func buildGBot(conf AppConf) *gbot.GBot {
	c := req.C().EnableInsecureSkipVerify()
	dcp := puppy.NewDcpServ(conf.DcpConf.URL, c)
	weather := puppy.NewAmapWeather(conf.AmapWeatherConf.Key, c)

	return &gbot.GBot{
		Weater: weather,
		Dcp:    dcp,
		Tpl:    gbot.GBotTemplate(),
	}
}

func buildHandle(conf AppConf, h workwx.RxMessageHandler) http.Handler {
	x, err := workwx.NewHTTPHandler(conf.WxAppConf.WxAppReceiveConf.Token,
		conf.WxAppConf.WxAppReceiveConf.AESKey, h)
	if err != nil {
		log.Fatal(err)
	}
	return x
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
		WxAppConf: WxAppConf{
			CorpID:     "12345",
			CorpSecret: "45679",
			AgentID:    1000013,
			WxAppReceiveConf: WxAppReceiveConf{
				Token:  "abcde",
				AESKey: "xxxxx",
			},
		},
		DcpConf: DcpConf{
			URL: "https://127.0.0.1:10005/jsonrpc",
		},
	}
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	enc.Encode(&conf)
}
