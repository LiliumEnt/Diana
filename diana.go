// made by amy (amyadzuki@gmail.com)
// this file is released into the public domain

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	///////////////////////////////////////////////////////////////////////////////
	////////////////////////  BEGIN CUSTOMIZATION SECTION  ////////////////////////
	///////////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////////
	// FROM THE BOT PAGE IN "MY APPS"; PLEASE DON'T PUT IT IN QUOTATION MARKS  ////
	ClientID = 440391423083282432

	///////////////////////////////////////////////////////////////////////////////
	// FROM THE BOT PAGE IN "MY APPS"; PLEASE DO PUT IT IN QUOTATION MARKS  ///////
	Token = "YOUR TOKEN GOES HERE"

	///////////////////////////////////////////////////////////////////////////////
	// ABOUT YOU AND YOUR SERVER; PLEASE DON'T PUT THEM IN QUOTATION MARKS  ///////
	ServerID  = 378599231583289346
	ChannelID = 440163315176701953
	OwnerID   = 77256980288253952
	BotUserID = 440391423083282432

	///////////////////////////////////////////////////////////////////////////////
	// CUSTOMIZE HERE /////////////////////////////////////////////////////////////
	PlayingStatus = "Tutoring Akko"
	//BannedText      = `discord.gg/`
	BannedWordWords = `n[iy](?:bb|gg)(?:a|er)|fag(?:got)?|cunt|hitler|r[ei]tard(?:ed)?`
	BannedWords     = `k[my]s`
	AutoDeleteAfter = 10 * time.Second

	///////////////////////////////////////////////////////////////////////////////
	/////////////////////////  END CUSTOMIZATION SECTION  /////////////////////////
	///////////////////////////////////////////////////////////////////////////////
	ReBannedWords = `(?i)\b(?:(?:` + BannedWordWords + `)s?|(?:` + BannedWords + `))\b`

//  ReBannedWords = `(?i)(?:(?:` + BannedText + `)|\b(?:(?:` +
//          BannedWordWords + `)s?|(?:` + BannedWords + `))\b)`
)

var BotID, CliID, GldID, ChnID, UsrID string
var REM, RELT, REBan *regexp.Regexp

func init() {
	BotID = strconv.FormatInt(BotUserID, 10)
	CliID = strconv.FormatInt(ClientID, 10)
	GldID = strconv.FormatInt(ServerID, 10)
	ChnID = strconv.FormatInt(ChannelID, 10)
	UsrID = strconv.FormatInt(OwnerID, 10)
	REM = regexp.MustCompile(`^<@!?` + BotID + `>\s*(.*)$`)
	RELT = regexp.MustCompile(`(^|[^\\])<`)
	REBan = regexp.MustCompile(ReBannedWords)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	t := *flag.String("t", Token, "Bot token")
	fmt.Println("This bot is ©2018 amy (amyadzuki@gmail.com).")
	dg, err := discordgo.New("Bot " + t)
	check(err)
	dg.AddHandler(onMessageCreate)
	check(dg.Open())
	defer dg.Close()
	fmt.Println("======= BOT UP (type Ctrl-C to exit)")
	defer fmt.Println("\n....... BOT DOWN")
	dg.UpdateStatus(0, PlayingStatus)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func onMessageCreate(session *discordgo.Session, arg *discordgo.MessageCreate) {
	message := arg.Message
	//  // Don't scan the bot's own messages:
	//  if message.Author.ID == BotID {
	//      return
	//  }
	// Don't scan messages by logging bots:
	if message.Author.Bot {
		return
	}
	text := message.Content
	if REBan.MatchString(text) {
		// Delete their message immediately:
		session.ChannelMessageDelete(message.ChannelID, message.ID)

		// Log it to the log channel:
		embed := new(discordgo.MessageEmbed)
		embed.Color = 0xff00bf
		embed.Title = "Member said a banned word"
		embed.Description = "<@!" + message.Author.ID + "> " +
			message.Author.Username + "\nUser ID: " + message.Author.ID +
			"\nChannel: <#" + message.ChannelID + "> (ID " + message.ChannelID + ")"
		embed.Fields = make([]*discordgo.MessageEmbedField, 1, 1)
		embed.Fields[0] = new(discordgo.MessageEmbedField)
		embed.Fields[0].Name = "They said:"
		embed.Fields[0].Value = text
		embed.Author = new(discordgo.MessageEmbedAuthor)
		embed.Author.Name = message.Author.Username
		embed.Author.IconURL = message.Author.AvatarURL("512")
		session.ChannelMessageSendEmbed(ChnID, embed)

		channel, err := session.UserChannelCreate(message.Author.ID)
		if err == nil {
			// Scold them in DMs
			session.ChannelMessageSend(channel.ID,
				"Greetings "+message.Author.Username+", I must inform you "+
					"that your message on :fleur_de_lis: **Lilium Ent.** has been "+
					"deleted as it contained a banned word.  We would be grateful "+
					"if you would please make sure to not repeat it.  Thank you.")
		} else {
			// Scold them in main chat cause they had DMs disabled
			msg, err := session.ChannelMessageSend(message.ChannelID,
				"<@!"+message.Author.ID+"> Excuse my interruption, however "+
					"I must ask of you to not make use of that word as it's banned "+
					"within our server, we hope you have the heart to understand.  "+
					"Thank you very much.")
			if err == nil {
				go AutoDelete(session, msg)
			}
		}
	}
}

func AutoDelete(session *discordgo.Session, message *discordgo.Message) {
	time.Sleep(AutoDeleteAfter)
	session.ChannelMessageDelete(message.ChannelID, message.ID)
}
