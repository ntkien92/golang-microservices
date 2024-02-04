package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	cadenceclient "go.uber.org/cadence/client"
)

const (
	ServerModeRelease = "release"
	ServerModeDebug   = "debug"

	cadenceService = "cadence-frontend"
)

var (
	logger        *zap.Logger
	cadenceClient cadenceclient.Client
)

func mode() string {
	return viper.GetString("config.mode")
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "./config.yaml", "[optional] path of configuration file")
	flag.Parse()

	loadConfig(configFile)

	// Init logger
	initLog()

	// Start worker
	domain := viper.GetString("cadence.domain")
	taskList := viper.GetString("cadence.tasklist")
	cadenceServer := viper.GetString("cadence.server")
	workerClient, err := buildWorkerClient(cadenceServer)
	if nil != err {
		logger.Panic("Fail to init cadence client", zap.Error(err))
	}

	worker, err := startWorker(logger, workerClient, domain, taskList)
	if nil != err {
		logger.Panic("Fail to start worker", zap.Error(err))
	}

	// Start cadence client
	cadenceService := buildCadenceService()
	cadenceClient = cadenceclient.NewClient(cadenceService, viper.GetString("cadence.domain"), nil)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info("Server is preparing to shutdown")
	worker.Stop()
}

func loadConfig(file string) {
	// Config from file
	viper.SetConfigType("yaml")
	if file != "" {
		viper.SetConfigFile(file)
	}

	viper.AddConfigPath("/.config/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No config file. Read config from env.")
		viper.AllowEmptyEnv(false)
	}

	// Config from env if possible
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GOLANG")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initLog() {
	var config zap.Config
	if mode() == ServerModeRelease {
		config = zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.InfoLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.StacktraceKey = ""
		config.EncoderConfig.TimeKey = ""
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level.SetLevel(zapcore.DebugLevel)
	}

	var err error
	l, err := config.Build()
	if err != nil {
		panic("Failed to setup logger")
	}

	logger = l
}
