package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	nhentaiAPI "github.com/VTGare/boe-tea-go/nhentai"
	"github.com/VTGare/boe-tea-go/utils"
	"github.com/VTGare/gumi"
	"github.com/bwmarrin/discordgo"
)

func init() {
	nsfwG := CommandFramework.AddGroup("nsfw", gumi.GroupDescription("All NSFW commands dwell here."), gumi.GroupNSFW())
	nhCmd := nsfwG.AddCommand("nhentai", nhentai, gumi.CommandDescription("Posts detailed info about nhentai book"), gumi.WithAliases("nh"))
	nhCmd.Help.ExtendedHelp = []*discordgo.MessageEmbedField{
		{
			Name:  "Usage",
			Value: "bt!nhentai <magic number>",
		},
		{
			Name:  "magic number",
			Value: "Typically, but not always, a 6-digit number only weebs understand.",
		},
	}
}

func nhentai(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		return utils.ErrNotEnoughArguments
	}

	if _, err := strconv.Atoi(args[0]); err != nil {
		return errors.New("invalid nhentai ID")
	}

	book, err := nhentaiAPI.GetNHentai(args[0])
	if err != nil {
		return err
	}

	artists := ""
	tags := ""
	if str := strings.Join(book.Artists, ", "); str != "" {
		artists = str
	} else {
		artists = "-"
	}

	if str := strings.Join(book.Tags, ", "); str != "" {
		tags = str
	} else {
		tags = "-"
	}

	embed := &discordgo.MessageEmbed{
		URL:   book.URL,
		Title: book.Titles.Pretty,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: book.Cover,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Artists",
				Value: artists,
			}, {
				Name:  "Tags",
				Value: tags,
			}, {
				Name:  "Favourites",
				Value: fmt.Sprintf("%v", book.Favourites),
			}, {
				Name:  "Pages",
				Value: fmt.Sprintf("%v", book.Pages),
			},
		},
		Color:     utils.EmbedColor,
		Timestamp: utils.EmbedTimestamp(),
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	return nil
}
