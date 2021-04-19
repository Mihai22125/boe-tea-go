package post

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestPost_skipArtworks(t *testing.T) {
	type fields struct {
		indices  map[int]struct{}
		skipMode SkipMode
	}
	tests := []struct {
		name   string
		fields fields
		embeds []*discordgo.MessageEmbed
		want   []*discordgo.MessageEmbed
	}{
		{
			name: "exclude",
			fields: fields{
				indices:  map[int]struct{}{1: {}, 2: {}},
				skipMode: SkipModeExclude,
			},
			embeds: []*discordgo.MessageEmbed{
				{Title: "1"}, {Title: "2"}, {Title: "3"},
			},
			want: []*discordgo.MessageEmbed{
				{Title: "3"},
			},
		},
		{
			name: "include",
			fields: fields{
				indices:  map[int]struct{}{1: {}, 2: {}},
				skipMode: SkipModeInclude,
			},
			embeds: []*discordgo.MessageEmbed{
				{Title: "1"}, {Title: "2"}, {Title: "3"},
			},
			want: []*discordgo.MessageEmbed{
				{Title: "1"}, {Title: "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := Post{indices: tt.fields.indices, skipMode: tt.fields.skipMode}

			res := post.skipArtworks(tt.embeds)
			assert.Equal(t, tt.want, res)
		})
	}
}
