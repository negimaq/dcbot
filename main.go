package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/negimaq/dcbot/handler"
)

var (
	commands = []*discordgo.ApplicationCommand{
		handler.TeamCommand,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"team": handler.TeamCommandHandler,
	}
)

func main() {
	// Discord token の取得
	token, ok := os.LookupEnv("DCBOT_TOKEN")
	if !ok {
		log.Fatal("環境変数 DCBOT_TOKEN をセットしてください")
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("セッションの作成に失敗しました:", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		log.Fatalln("Discord への接続に失敗しました:", err)
	}
	defer s.Close()

	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Fatalf("%v コマンドの作成に失敗しました: %v", v.Name, err)
		}
	}

	log.Print("Bot を起動しました（CTRL-C で終了）")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	log.Print("Bot をシャットダウンしています")
}
