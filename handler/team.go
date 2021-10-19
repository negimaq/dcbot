package handler

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	teamEmoji = [...]string{
		":red_circle:",
		":blue_circle:",
		":green_circle:",
		":orange_circle:",
		":purple_circle:",
		":yellow_circle:",
		":brown_circle:",
		":white_circle:",
		":black_circle:",
	}

	TeamCommand = &discordgo.ApplicationCommand{
		Name:        "team",
		Description: "チーム分けします。",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "チーム数",
				Description: fmt.Sprintf("作成したいチーム数（2以上%d以下）", len(teamEmoji)),
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "参加者",
				Description: "参加者（スペース区切り）",
				Required:    true,
			},
		},
	}
)

func TeamCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg := ""
	n := int(i.ApplicationCommandData().Options[0].IntValue())                   // チーム数
	p := strings.Split(i.ApplicationCommandData().Options[1].StringValue(), " ") // 参加者

	if n < 2 { // チーム数が2未満の場合
		msg = fmt.Sprintf("チーム数は2以上%d以下を指定してください", len(teamEmoji))
	} else if len(p) < 2 { // 参加者が2人未満の場合
		msg = "参加者が2人以上必要です"
	} else {
		if len(teamEmoji) < n { // チーム数がチーム絵文字の数より大きい場合
			n = len(teamEmoji)
		}
		if len(p) < n { // チーム数が参加者数よりも大きい場合
			n = len(p)
		}

		// 参加者をシャッフル
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })

		// チーム分けのインデックスを計算
		size, rem := len(p)/n, len(p)%n
		divs := make([]int, 1, n+1)
		for i := 0; i < n; i++ {
			if i < rem {
				divs = append(divs, size+1)
			} else {
				divs = append(divs, size)
			}
		}
		for i := 1; i < len(divs); i++ {
			divs[i] += divs[i-1]
		}

		for i := 0; i < n; i++ {
			start, end := divs[i], divs[i+1]
			msg += fmt.Sprintln(teamEmoji[i], strings.Join(p[start:end], " "))
			if i < n-1 {
				msg += fmt.Sprintln(" vs.")
			}
		}
	}

	// レスポンスを返す
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: "メッセージの送信に失敗しました",
		})
	}
}
