package commands

import (
	"strings"

	"github.com/VTGare/boe-tea-go/bot"
	"github.com/VTGare/boe-tea-go/messages"
	"github.com/VTGare/embeds"
	"github.com/VTGare/gumi"
)

func ownerGroup(b *bot.Bot) {
	group := "owner"

	b.Router.RegisterCmd(&gumi.Command{
		Name:        "reply",
		Group:       group,
		Description: "Owner's command to reply to feedback",
		Usage:       "bt!reply <wall of text>",
		Example:     "bt!reply You know who else is shit? Yrou'e mom in bed :^)",
		AuthorOnly:  true,
		Exec:        reply(b),
	})
}

func reply(b *bot.Bot) func(ctx *gumi.Ctx) error {
	return func(ctx *gumi.Ctx) error {
		if ctx.Args.Len() < 2 {
			return messages.ErrIncorrectCmd(ctx.Command)
		}

		userID := strings.Trim(ctx.Args.Get(0).Raw, "<@!>")

		s := b.ShardManager.SessionForDM()
		ch, err := s.UserChannelCreate(userID)
		if err != nil {
			return err
		}

		eb := embeds.NewBuilder()

		eb.Author(
			"Feedback reply",
			"",
			ctx.Session.State.User.AvatarURL(""),
		).Description(
			strings.TrimPrefix(
				strings.TrimSpace(ctx.Args.Raw),
				ctx.Args.Get(0).Raw,
			),
		)

		if attachments := ctx.Event.Attachments; len(attachments) >= 1 {
			if strings.HasSuffix(attachments[0].Filename, "png") || strings.HasSuffix(attachments[0].Filename, "jpg") || strings.HasSuffix(attachments[0].Filename, "gif") {
				eb.Image(attachments[0].URL)
			}
		}

		_, err = s.ChannelMessageSendEmbed(ch.ID, eb.Finalize())
		if err != nil {
			return err
		}

		eb.Clear()
		ctx.ReplyEmbed(eb.SuccessTemplate("Reply has been sent.").Finalize())
		return nil
	}
}
