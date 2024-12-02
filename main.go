package main

import (
	"flag"
	// "fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	// "strconv"

	"secretsanta/structs"
	"secretsanta/utils"

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

// define json file to store info
var jsonFile string = "./data/secretsanta.json"

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
			Name:			"start",
			Description:	"Starts secret santa event!",
		},
		{
			Name:			"end",
			Description:	"Ends secret santa event and notifies participants of their person to buy for.",
		},
		{
			Name:			"cancel",
			Description:	"Cancels secret santa event and does not notify participants.",
		},
		{
			Name:			"register",
			Description:	"Join secret santa event.",
		},
		{
			Name:			"unregister",
			Description:	"Leave secret santa event.",
		},
		{
			Name:			"edit",
			Description:	"Edit Secret Santa Event.",
		},
	}

	// defines what happens when a named command is run
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"start": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			eventStarted := utils.IsEventStarted(jsonFile, i.GuildID)

			if eventStarted {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: 	"An event has already been started for this server. You must cancel or end the event before ",
						Flags:		discordgo.MessageFlagsEphemeral,
					},
				})

				panicIfError(err)

				return
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "start_modal_" + i.Interaction.Member.User.ID,
					Title: "Secret Santa Setup",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:		"name",
									Label:			"What is the name of the event?",
									Style:			discordgo.TextInputShort,
									Placeholder:	"Secret Santa Gift Exchange",
									Required:		true,
									MaxLength:		100,
									MinLength:		1,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:		"price",
									Label:			"What is the price range of the event?",
									Style:			discordgo.TextInputShort,
									Placeholder:	"$15-20",
									Required:		true,
									MaxLength:		100,
									MinLength:		1,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:		"register",
									Label:			"When should participants register by?",
									Style:			discordgo.TextInputShort,
									Placeholder:	"December 7th, 2024",
									Required:		true,
									MaxLength:		100,
									MinLength:		1,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:		"deadline",
									Label:			"When is the event?",
									Style:			discordgo.TextInputShort,
									Placeholder:	"December 14th, 2024",
									Required:		true,
									MaxLength:		100,
									MinLength:		1,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:		"notes",
									Label:			"Any extra notes?",
									Style:			discordgo.TextInputParagraph,
									Placeholder:	"",
									Required:		false,
									MaxLength:		500,
									MinLength:		1,
								},
							},
						},
					},
				},
			})

			panicIfError(err)
		},
		"end": func(s *discordgo.Session, i *discordgo.InteractionCreate){

		},
		"cancel": func(s *discordgo.Session, i *discordgo.InteractionCreate){

		},
		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate){

		},
		"unregister": func(s *discordgo.Session, i *discordgo.InteractionCreate){

		},
		"edit": func(s *discordgo.Session, i *discordgo.InteractionCreate){

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
					Content:	"Thank you for setting up a new event!",
					Flags:		discordgo.MessageFlagsEphemeral,
				},
			})

			data := i.ModalSubmitData()

			if strings.HasPrefix(data.CustomID, "start_modal"){
				curGuildData := structs.GuildData{
					Creator: i.Member.User.ID,
					EmbedID: "",
					EmbedData: structs.EmbededData{
						Name: data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Price: data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Register: data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Deadline: data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Notes: data.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
					},
					Responses: make(map[string]string),
					Santas: make(map[string]string),
					Ended: false,
				} 

				msg, embedErr := s.ChannelMessageSendComplex(i.Interaction.ChannelID, &discordgo.MessageSend{
					Embeds: []*discordgo.MessageEmbed{
						&discordgo.MessageEmbed{
							Type: "rich",
							Title: curGuildData.EmbedData.Name,
							Description: "Keep teack of everything for this year's Secret Santa event!",
							Color: 3447003,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name: "Registration Date",
									Value: curGuildData.EmbedData.Register,
									Inline: true,
								},
								{
									Name: "Event Date",
									Value: curGuildData.EmbedData.Deadline,
									Inline: true,
								},
								{
									Name: "Price Limit",
									Value: curGuildData.EmbedData.Price,
									Inline: true,
								},
								{
									Name: "Participants Signed Up",
									Value: "Nobody yet!",
									Inline: false,
								},
								{
									Name: "Want to join or leave?",
									Value: "Use `/register` to join or `/unregister` to leave. Or press the buttons below!",
									Inline: false,
								},
								{
									Name: "Wishlist Submission",
									Value: "Use `/wishlist` to set your wishlist!",
									Inline: false,
								},
								{
									Name: "Notes",
									Value: curGuildData.EmbedData.Notes,
									Inline: false,
								},
							},
							Footer: &discordgo.MessageEmbedFooter{
								Text: "Built by regexbuster",
							},
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: "register_" + i.Interaction.Member.User.ID,
									Label: "Register",
									Style: discordgo.SuccessButton,
									Disabled: false,
								},
								discordgo.Button{
									CustomID: "unregister_" + i.Interaction.Member.User.ID,
									Label: "Unregister",
									Style: discordgo.DangerButton,
									Disabled: false,
								},
								discordgo.Button{
									CustomID: "wishlist_" + i.Interaction.Member.User.ID,
									Label: "Wishlist",
									Style: discordgo.PrimaryButton,
									Disabled: false,
								},
							},
						},
					},
				})

				panicIfError(embedErr)

				curGuildData.EmbedID = msg.ID

				var jsonData map[string]structs.GuildData

				utils.ReadJSONFile(jsonFile, &jsonData)

				jsonData[i.GuildID] = curGuildData

				utils.WriteJSONFile(jsonFile, &jsonData)				
			}

			panicIfError(err)
			// interactions with messages (ie button presses)
		case discordgo.InteractionMessageComponent:
			// data := i.MessageComponentData()
			return;
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