package main

import (
	"context"
	"encoding/json"
	"github.com/corner4world/plugin-gbclient/internal/conf"
	"io/ioutil"

	"github.com/corner4world/plugin-gbclient/internal/useragent"
	"github.com/corner4world/plugin-gbclient/internal/version"
	cli "github.com/jawher/mow.cli"
	"github.com/qiniu/x/xlog"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
	"os"
)

var sim_conf = `{
	"localSipPort":5061,
	"serverID":"41010500002000000001",
	"realm":"4101050000",
	"serverAddr":"192.168.1.14:5060",
	"userName":"admin",
	"password":"admin",
	"regExpire":3600,
	"keepaliveInterval":100,
	"maxKeepaliveRetry":3,
	"transport": "udp",
	"gbId":"31011500991320000532",
    	"devices":[
        {
            "deviceID":"32011500991320000040",
            "name":"test001",
            "manufacturer":"simulatorFactory",
            "model":"Mars",
            "CivilCode":"civilCode",
            "address":"192.18.1.1",
            "parental":"0",
            "safeWay":"1",
            "registerWay":"1",
            "secrecy":"1",
            "status":"ON"
        },
        {
            "deviceID":"32011500991320000041",
            "name":"test002",
            "manufacturer":"simulatorFactory",
            "model":"Mars",
            "CivilCode":"civilCode",
            "address":"192.18.1.2",
            "parental":"0",
            "safeWay":"1",
            "registerWay":"1",
            "secrecy":"1",
            "status":"ON"
        },
        {
            "deviceID":"32011500991320000042",
            "name":"test003",
            "manufacturer":"simulatorFactory",
            "model":"Mars",
            "CivilCode":"civilCode",
            "address":"192.18.1.3",
            "parental":"0",
            "safeWay":"1",
            "registerWay":"1",
            "secrecy":"1",
            "status":"ON"
        },
        {
            "deviceID":"32011500991320000043",
            "name":"test004",
            "manufacturer":"simulatorFactory",
            "model":"Mars",
            "CivilCode":"civilCode",
            "address":"192.18.1.4",
            "parental":"0",
            "safeWay":"1",
            "registerWay":"1",
            "secrecy":"1",
            "status":"OFF"
        }
    ]
}`

type GBClientConfig struct {
	config.HTTP
	config.Publish
	config.Subscribe

	LocalSipPort      int    `default:"5061" desc:"本地端口"`
	ServerID          string `default:"41010500002000000001" desc:"服务器编码"`
	Realm             string `default:"4101050000" desc:"区域编码"`
	ServerAddr        string `default:"192.168.1.14:5060" desc:"远端端口"`
	UserName          string `default:"admin" desc:"用户"`
	Password          string `default:"admin" desc:"密码"`
	RegExpire         int    `default:"3600" desc:"有效期"`
	KeepaliveInterval int    `default:"100" desc:"保活时间间隔"`
	MaxKeepaliveRetry int    `default:"3" desc:"最大保活重试次数"`
	Transport         string `default:"udp" desc:"最大保活重试次数"`
	GBID              string
	Devices           []DeviceInfo
}

type DeviceInfo struct {
	Text         string `xml:",chardata"`
	DeviceID     string `xml:"DeviceID" json:"deviceID"`
	Name         string `xml:"Name" json:"name""`
	Manufacturer string `xml:"Manufacturer" json:"manufacturer"`
	Model        string `xml:"Model" json:"model"`
	Owner        string `xml:"Owner" json:"owner"`
	CivilCode    string `xml:"CivilCode" json:"civilCode"`
	Address      string `xml:"Address" json:"address"`
	Parental     string `xml:"Parental" json:"parental"`
	SafetyWay    string `xml:"SafetyWay" json:"safeWay"`
	RegisterWay  string `xml:"RegisterWay" json:"registerWay"`
	Secrecy      string `xml:"Secrecy" json:"secrecy"`
	Status       string `xml:"Status" json:"status"`
}

var clientConfig GBClientConfig

// 安装插件
var GBClientPlugin = InstallPlugin(&clientConfig)

// 插件事件回调，来自事件总线
func (conf *GBClientConfig) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		xlog.SetOutputLevel(0)
		xlog.SetFlags(xlog.Llevel | xlog.Llongfile | xlog.Ltime)
		xlog := xlog.NewWith(context.Background())
		app := cli.App("gb28181Client", "Runs the gb28181 client.")
		app.Spec = "[ -c=<configuration path> ] "
		confPath := app.StringOpt("c config", "sim.conf", "Specifies the configuration path (file) to use for the client.")
		app.Action = func() { run(xlog, app, confPath) }

		// Register sub-commands
		app.Command("version", "Prints the version of the executable.", version.Print)
		app.Run(os.Args)
	}
}

func run(xlog *xlog.Logger, app *cli.Cli, conf *string) {
	xlog.Infof("gb28181 client is running...")
	cfg, err := ParseJsonConfig(conf)
	if err != nil {
		xlog.Errorf("load config file failed, err = ", err)
	}
	xlog.Infof("config file = %#v", cfg)
	srv, err := useragent.NewService(xlog, cfg)
	if err != nil {
		xlog.Infof("new service failed err = %#v", err)
		return
	}
	srv.HandleIncommingMsg()
}

func ParseJsonConfig(f *string) (*conf.Config, error) {
	jsonFile, err := os.Open(*f)
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var cfg conf.Config
	err = json.Unmarshal(b, &cfg)
	return &cfg, err
}
