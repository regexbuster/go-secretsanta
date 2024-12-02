package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// define variables via cmd flags that allow bot to start and
var (
	GuildID 		= flag.String("guild", "", "Test Guild ID")
	BotToken		= flag.String("token", "", "Bot access token")
	AppID			= flag.String("app", "", "Application ID")
	Cleanup			= flag.Bool("cleanup", true, "Cleanup of commands")
	ResultsChannel	= flag.String("results", "", "Channel to send survey results to")
)

// define global session before doing anything
var s *discordgo.Session

// send potential error and it panics if error exists
// simplification for writing these all over the place
func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// init functions run before main and is used to set up the rest of the code
func init() {
	// takes in flag configs and assigns them to the above variables
	flag.Parse()
}

func init() {
	// pre-define because s is already globally defined (err is considered undefined if this is not here)
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	// define all the commands that we want to have
	commands = []discordgo.ApplicationCommand{
		{
			Name:			"modals-survey",
			Description:	"Take a survey about modals",
		},
		{
			Name:			"ping",
			Description:	"Returns pong",
		},
		{
			Name:			"buttons",
			Description:	"testing buttons",
		},
		{
			Name:			"embed",
			Description:	"testing embed",
		},
		{
			Name:			"editing",
			Description:	"testing editing",
		},
	}

	// defines what happens when a named command is run
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"modals-survey": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "modals_survey_" + i.Interaction.Member.User.ID,
					Title:	"Modals Survey",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:		"opinion",
									Label:			"What is your opinion on them?",
									Style:			discordgo.TextInputShort,
									Placeholder:	"Don't be shy, share your opinion.",
									Required:		true,
									MaxLength:		300,
									MinLength:		10,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:	"suggestions",
									Label:		"What would you suggest to improve",
									Style:		discordgo.TextInputParagraph,
									Required:	false,
									MaxLength:	2000,
								},
							},
						},
					},
				},
			})
			panicIfError(err)
		},
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					CustomID: "pong_" + i.Interaction.Member.User.ID,
					Title:	"Pong",
					Content: "pong?",
				},
			})
			panicIfError(err)
		},
		"buttons": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					CustomID: "buttons_" + i.Interaction.Member.User.ID,
					Title: "Buttons",
					Content: "button testing",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: "button_primary_" + i.Interaction.Member.User.ID,
									Label: "Click Me",
									Style: discordgo.PrimaryButton,
									Disabled: false,
								},
								discordgo.Button{
									CustomID: "button_second_primary_" + i.Interaction.Member.User.ID,
									Label: "Don't Click Me",
									Style: discordgo.DangerButton,
									Disabled: false,
								},
							},
						},
					},
				},
			})
			panicIfError(err)
		},
		"embed": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					CustomID: "embed" + i.Interaction.Member.User.ID,
					Title: "Embed",
					Content: "embed",
					Embeds: []*discordgo.MessageEmbed{
						{
							Type: "rich",
							Title: "Secret Santa Event Tracker",
							Description: "Keep track of everything for this year's Secret Santa event!",
							Color: 3447003,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name: "Event Date",
									Value: "December 24th, 2024",
									Inline: true,
								},
								{
									Name: "Price Limit",
									Value: "$20 - $50",
									Inline: true,
								},
								{
									Name: "Participants Signed Up",
									Value: "regexbuster, SparkFluxx, etc...",
									Inline: false,
								},
								{
									Name: "Want to join or leave?",
									Value: "Use `/register` to join or `/unregister` to leave.",
									Inline: false,
								},
								{
									Name: "Wishlist Submission",
									Value: "Use `/wishlist` to set your wishlist!",
									Inline: false,
								},
								{
									Name: "Notes",
									Value: "Yap yap yap...",
									Inline: false,
								},
							},
							Footer: &discordgo.MessageEmbedFooter{
								Text: "Built by regexbuster",
							},
						},
					},
				},
			})
			panicIfError(err)
		},
		"editing": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					CustomID: "editing_" + i.Interaction.Member.User.ID,
					Title: "Editing",
					Content: "0",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: "increase_" + i.Interaction.Member.User.ID,
									Label: "Add One",
									Style: discordgo.PrimaryButton,
									Disabled: false,
								},
							},
						},
					},
				},
			})
			panicIfError(err)
		},
	}
)

func main() {
	// event handler that triggers when discordgo is ready (bot is up and connected)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready){
		log.Println("Bot is up!")
	})

	// event handler for when any interaction with bot it created (slash commands and modal submit etc)
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate){
		switch i.Type{
			// when slash command called
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s,i)
			}
			// when modal submitted
		case discordgo.InteractionModalSubmit:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:	"Thank you for taking your time to fill this survey!",
					Flags:		discordgo.MessageFlagsEphemeral,
				},
			})

			panicIfError(err)

			data := i.ModalSubmitData()

			if !strings.HasPrefix(data.CustomID, "modals_survey"){
				return
			}

			userid := strings.Split(data.CustomID, "_")[2]

			_, err = s.ChannelMessageSend(*ResultsChannel, fmt.Sprintf(
				"Feedback received. From <@%s>\n\n**Opinion**:\n%s\n\n**Suggestions**:\n%s",
				userid,
				data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
				data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
			))

			panicIfError(err)
			// interactions with messages (ie button presses)
		case discordgo.InteractionMessageComponent:
			data := i.MessageComponentData()

			if strings.HasPrefix(data.CustomID, "increase"){
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				})
	
				panicIfError(err)

				editedMessage := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)

				numMessage, _ := strconv.Atoi(i.Message.Content)
				numMessage += 1

				editedMessage.SetContent(strconv.Itoa(numMessage));

				s.ChannelMessageEditComplex(editedMessage);

				return
			}			

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:	fmt.Sprintf("%s", data),
					Flags:		discordgo.MessageFlagsEphemeral,
				},
			})

			panicIfError(err)

		}
	})

	cmdIDs := make(map[string]string, len(commands))

	// register slash commands
	for _, cmd := range commands {
		rcmd, err := s.ApplicationCommandCreate(*AppID, *GuildID, &cmd)

		if err != nil {
			log.Fatalf("Cannot create slash commands %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}

	// open session, check for error, defer closing session until program end
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	// create channel watching for interrupt signal from user
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<- stop
	log.Println("Graceful shutdown")

	if !*Cleanup {
		return
	}

	// unregister slash commands
	for id, name := range cmdIDs {
		err := s.ApplicationCommandDelete(*AppID, *GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
}