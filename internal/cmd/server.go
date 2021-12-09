package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ocr-pub/internal/server"
)

type ServerCmd struct {
	command *cobra.Command
}

type Config struct{
	redisHost      string
	redisDB        int8
	socketPort     int32
	webSocketPort  int32
}

func NewServerCommand() *ServerCmd{
	command := ServerCmd{

	}
	config := Config{}
	command.command = &cobra.Command{
		Use: "1580 ocr result publish system",
		Short: "1580 ocr result publish system",
		Long: "使用两种方式socket长连接和websocket，通过参数选择使用或者全部使用，消息来源于redis订阅频道",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(config)
			a, _ := server.NewServer(config.redisHost, config.redisDB,config.socketPort, config.webSocketPort)
			// fmt.Println(a)
			a.Start()
			return nil
		},
	}
	command.command.Flags().StringVarP(&config.redisHost, "redisHost", "r", ":6379", "the host of redis")
	command.command.Flags().Int8VarP(&config.redisDB, "redisDB", "d", 0, "the db of redis")
	command.command.Flags().Int32VarP(&config.socketPort, "socketPort", "s", 9999, "the port of socket")
	command.command.Flags().Int32VarP(&config.webSocketPort, "websocketPort", "w", 8082, "the port of websocket")
	return &command
}

func (cmd *ServerCmd) GetCommand() *cobra.Command {
	return cmd.command
}