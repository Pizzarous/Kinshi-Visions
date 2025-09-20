package discord_bot

import (
	"errors"
	"fmt"
	"kinshi_vision_bot/entities"
	"kinshi_vision_bot/invision_queue"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type botImpl struct {
	developmentMode    bool
	botSession         *discordgo.Session
	guildID            string
	invisionQueue      invision_queue.Queue
	registeredCommands []*discordgo.ApplicationCommand
	invisionCommand    string
	removeCommands     bool
}

type Config struct {
	DevelopmentMode bool
	BotToken        string
	GuildID         string
	InvisionQueue   invision_queue.Queue
	InvisionCommand string
	RemoveCommands  bool
}

func (b *botImpl) invisionCommandString() string {
	if b.developmentMode {
		return "dev_" + b.invisionCommand
	}

	return b.invisionCommand
}

func (b *botImpl) invisionSettingsCommandString() string {
	if b.developmentMode {
		return "dev_" + b.invisionCommand + "_settings"
	}

	return b.invisionCommand + "_settings"
}

func New(cfg Config) (Bot, error) {
	if cfg.BotToken == "" {
		return nil, errors.New("missing bot token")
	}

	if cfg.GuildID == "" {
		return nil, errors.New("missing guild ID")
	}

	if cfg.InvisionQueue == nil {
		return nil, errors.New("missing invision queue")
	}

	if cfg.InvisionCommand == "" {
		return nil, errors.New("missing invision command")
	}

	botSession, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}

	botSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err = botSession.Open()
	if err != nil {
		return nil, err
	}

	bot := &botImpl{
		developmentMode:    cfg.DevelopmentMode,
		botSession:         botSession,
		invisionQueue:      cfg.InvisionQueue,
		registeredCommands: make([]*discordgo.ApplicationCommand, 0),
		invisionCommand:    cfg.InvisionCommand,
		removeCommands:     cfg.RemoveCommands,
	}

	err = bot.addInvisionCommand()
	if err != nil {
		return nil, err
	}

	err = bot.addInvisionSettingsCommand()
	if err != nil {
		return nil, err
	}

	botSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case bot.invisionCommandString():
				bot.processInvisionCommand(s, i)
			case bot.invisionSettingsCommandString():
				bot.processInvisionSettingsCommand(s, i)
			default:
				log.Printf("Unknown command '%v'", i.ApplicationCommandData().Name)
			}
		case discordgo.InteractionMessageComponent:
			switch customID := i.MessageComponentData().CustomID; {
			case customID == "invision_reroll":
				bot.processInvisionReroll(s, i)
			case strings.HasPrefix(customID, "invision_upscale_"):
				interactionIndex := strings.TrimPrefix(customID, "invision_upscale_")

				interactionIndexInt, intErr := strconv.Atoi(interactionIndex)
				if intErr != nil {
					log.Printf("Error parsing interaction index: %v", err)

					return
				}

				bot.processInvisionUpscale(s, i, interactionIndexInt)
			case strings.HasPrefix(customID, "invision_variation_"):
				interactionIndex := strings.TrimPrefix(customID, "invision_variation_")

				interactionIndexInt, intErr := strconv.Atoi(interactionIndex)
				if intErr != nil {
					log.Printf("Error parsing interaction index: %v", err)

					return
				}

				bot.processInvisionVariation(s, i, interactionIndexInt)
			case customID == "invision_dimension_setting_menu":
				if len(i.MessageComponentData().Values) == 0 {
					log.Printf("No values for invision dimension setting menu")

					return
				}

				sizes := strings.Split(i.MessageComponentData().Values[0], "_")

				width := sizes[0]
				height := sizes[1]

				widthInt, intErr := strconv.Atoi(width)
				if intErr != nil {
					log.Printf("Error parsing width: %v", err)

					return
				}

				heightInt, intErr := strconv.Atoi(height)
				if intErr != nil {
					log.Printf("Error parsing height: %v", err)

					return
				}

				bot.processInvisionDimensionSetting(s, i, widthInt, heightInt)

			// patch from upstream
			case customID == "invision_batch_count_setting_menu":
				if len(i.MessageComponentData().Values) == 0 {
					log.Printf("No values for invision batch count setting menu")

					return
				}

				batchCount := i.MessageComponentData().Values[0]

				batchCountInt, intErr := strconv.Atoi(batchCount)
				if intErr != nil {
					log.Printf("Error parsing batch count: %v", err)

					return
				}

				var batchSizeInt int

				// calculate the corresponding batch size
				switch batchCountInt {
				case 1:
					batchSizeInt = 4
				case 2:
					batchSizeInt = 2
				case 4:
					batchSizeInt = 1
				default:
					log.Printf("Unknown batch count: %v", batchCountInt)

					return
				}

				bot.processInvisionBatchSetting(s, i, batchCountInt, batchSizeInt)
			case customID == "invision_batch_size_setting_menu":
				if len(i.MessageComponentData().Values) == 0 {
					log.Printf("No values for invision batch count setting menu")

					return
				}

				batchSize := i.MessageComponentData().Values[0]

				batchSizeInt, intErr := strconv.Atoi(batchSize)
				if intErr != nil {
					log.Printf("Error parsing batch count: %v", err)

					return
				}

				var batchCountInt int

				// calculate the corresponding batch count
				switch batchSizeInt {
				case 1:
					batchCountInt = 4
				case 2:
					batchCountInt = 2
				case 4:
					batchCountInt = 1
				default:
					log.Printf("Unknown batch size: %v", batchSizeInt)

					return
				}

				bot.processInvisionBatchSetting(s, i, batchCountInt, batchSizeInt)

			default:
				log.Printf("Unknown message component '%v'", i.MessageComponentData().CustomID)
			}
		}
	})

	return bot, nil
}

func (b *botImpl) Start() {
	b.invisionQueue.StartPolling(b.botSession)

	err := b.teardown()
	if err != nil {
		log.Printf("Error tearing down bot: %v", err)
	}
}

func (b *botImpl) teardown() error {
	// Delete all commands added by the bot
	if b.removeCommands {
		log.Printf("Removing all commands added by bot...")

		for _, v := range b.registeredCommands {
			log.Printf("Removing command '%v'...", v.Name)

			err := b.botSession.ApplicationCommandDelete(b.botSession.State.User.ID, b.guildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	return b.botSession.Close()
}

func (b *botImpl) addInvisionCommand() error {
	log.Printf("Adding command '%s'...", b.invisionCommandString())

	cmd, err := b.botSession.ApplicationCommandCreate(b.botSession.State.User.ID, b.guildID, &discordgo.ApplicationCommand{
		Name:        b.invisionCommandString(),
		Description: "Ask the bot to invision something",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The text prompt to invision",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "negative_prompt",
				Description: "Negative prompt",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "sampler_name",
				Description: "sampler",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Euler a",
						Value: "Euler a",
					},
					{
						Name:  "DDIM",
						Value: "DDIM",
					},
					{
						Name:  "PLMS",
						Value: "PLMS",
					},
					{
						Name:  "UniPC",
						Value: "UniPC",
					},
					{
						Name:  "Heun",
						Value: "Heun",
					},
					{
						Name:  "Euler",
						Value: "Euler",
					},
					{
						Name:  "LMS",
						Value: "LMS",
					},
					{
						Name:  "LMS Karras",
						Value: "LMS Karras",
					},
					{
						Name:  "DPM2 a",
						Value: "DPM2 a",
					},
					{
						Name:  "DPM2 a Karras",
						Value: "DPM2 a Karras",
					},
					{
						Name:  "DPM2",
						Value: "DPM2",
					},
					{
						Name:  "DPM2 Karras",
						Value: "DPM2 Karras",
					},
					{
						Name:  "DPM fast",
						Value: "DPM fast",
					},
					{
						Name:  "DPM adaptive",
						Value: "DPM adaptive",
					},
					{
						Name:  "DPM++ 2S a",
						Value: "DPM++ 2S a",
					},
					{
						Name:  "DPM++ 2M",
						Value: "DPM++ 2M",
					},
					{
						Name:  "DPM++ SDE",
						Value: "DPM++ SDE",
					},
					{
						Name:  "DPM++ 2S a Karras",
						Value: "DPM++ 2S a Karras",
					},
					{
						Name:  "DPM++ 2M Karras",
						Value: "DPM++ 2M Karras",
					},
					{
						Name:  "DPM++ SDE Karras",
						Value: "DPM++ SDE Karras",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "use_hires_fix",
				Description: "use hires.fix or not. default=No for better performance",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Yes",
						Value: "true",
					},
					{
						Name:  "No",
						Value: "false",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error creating '%s' command: %v", b.invisionCommandString(), err)

		return err
	}

	b.registeredCommands = append(b.registeredCommands, cmd)

	return nil
}

func (b *botImpl) addInvisionSettingsCommand() error {
	log.Printf("Adding command '%s'...", b.invisionSettingsCommandString())

	cmd, err := b.botSession.ApplicationCommandCreate(b.botSession.State.User.ID, b.guildID, &discordgo.ApplicationCommand{
		Name:        b.invisionSettingsCommandString(),
		Description: "Change the default settings for the invision command",
	})
	if err != nil {
		log.Printf("Error creating '%s' command: %v", b.invisionSettingsCommandString(), err)

		return err
	}

	b.registeredCommands = append(b.registeredCommands, cmd)

	return nil
}

func (b *botImpl) processInvisionReroll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	position, queueError := b.invisionQueue.AddInvision(&invision_queue.QueueItem{
		Type:               invision_queue.ItemTypeReroll,
		DiscordInteraction: i.Interaction,
	})
	if queueError != nil {
		log.Printf("Error adding invision to queue: %v\n", queueError)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("I'm reimagining that for you... You are currently #%d in line.", position),
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func (b *botImpl) processInvisionUpscale(s *discordgo.Session, i *discordgo.InteractionCreate, upscaleIndex int) {
	position, queueError := b.invisionQueue.AddInvision(&invision_queue.QueueItem{
		Type:               invision_queue.ItemTypeUpscale,
		InteractionIndex:   upscaleIndex,
		DiscordInteraction: i.Interaction,
	})
	if queueError != nil {
		log.Printf("Error adding invision to queue: %v\n", queueError)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("I'm upscaling that for you... You are currently #%d in line.", position),
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func (b *botImpl) processInvisionVariation(s *discordgo.Session, i *discordgo.InteractionCreate, variationIndex int) {
	position, queueError := b.invisionQueue.AddInvision(&invision_queue.QueueItem{
		Type:               invision_queue.ItemTypeVariation,
		InteractionIndex:   variationIndex,
		DiscordInteraction: i.Interaction,
	})
	if queueError != nil {
		log.Printf("Error adding invision to queue: %v\n", queueError)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("I'm imagining more variations for you... You are currently #%d in line.", position),
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func (b *botImpl) processInvisionCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var position int
	var queueError error
	var prompt string
	negative := ""
	sampler := "DPM++ 2M"
	hiresfix := false

	if option, ok := optionMap["prompt"]; ok {
		prompt = option.StringValue()

		if nopt, ok := optionMap["negative_prompt"]; ok {
			negative = nopt.StringValue()
		}

		if smpl, ok := optionMap["sampler_name"]; ok {
			sampler = smpl.StringValue()
		}

		if hires, ok := optionMap["use_hires_fix"]; ok {
			hiresfix, _ = strconv.ParseBool(hires.StringValue())
		}

		position, queueError = b.invisionQueue.AddInvision(&invision_queue.QueueItem{
			Prompt:             prompt,
			NegativePrompt:     negative,
			SamplerName1:       sampler,
			Type:               invision_queue.ItemTypeInvision,
			UseHiresFix:        hiresfix,
			DiscordInteraction: i.Interaction,
		})
		if queueError != nil {
			log.Printf("Error adding invision to queue: %v\n", queueError)
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"I'm dreaming something up for you. You are currently #%d in line.\n<@%s> asked me to invision \"%s\", with sampler: %s",
				position,
				i.Member.User.ID,
				prompt,
				sampler),
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

// patch from upstream
func settingsMessageComponents(settings *entities.DefaultSettings) []discordgo.MessageComponent {
	minValues := 1

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:  "invision_dimension_setting_menu",
					MinValues: &minValues,
					MaxValues: 1,
					Options: []discordgo.SelectMenuOption{
						{
							Label:   "Size: 512x512",
							Value:   "512_512",
							Default: settings.Width == 512 && settings.Height == 512,
						},
						{
							Label:   "Size: 768x768",
							Value:   "768_768",
							Default: settings.Width == 768 && settings.Height == 768,
						},
					},
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:  "invision_batch_count_setting_menu",
					MinValues: &minValues,
					MaxValues: 1,
					Options: []discordgo.SelectMenuOption{
						{
							Label:   "Batch count: 1",
							Value:   "1",
							Default: settings.BatchCount == 1,
						},
						{
							Label:   "Batch count: 2",
							Value:   "2",
							Default: settings.BatchCount == 2,
						},
						{
							Label:   "Batch count: 4",
							Value:   "4",
							Default: settings.BatchCount == 4,
						},
					},
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:  "invision_batch_size_setting_menu",
					MinValues: &minValues,
					MaxValues: 1,
					Options: []discordgo.SelectMenuOption{
						{
							Label:   "Batch size: 1",
							Value:   "1",
							Default: settings.BatchSize == 1,
						},
						{
							Label:   "Batch size: 2",
							Value:   "2",
							Default: settings.BatchSize == 2,
						},
						{
							Label:   "Batch size: 4",
							Value:   "4",
							Default: settings.BatchSize == 4,
						},
					},
				},
			},
		},
	}
}

func (b *botImpl) processInvisionSettingsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	botSettings, err := b.invisionQueue.GetBotDefaultSettings()
	if err != nil {
		log.Printf("error getting default settings for settings command: %v", err)

		return
	}

	messageComponents := settingsMessageComponents(botSettings)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:      "Settings",
			Content:    "Choose defaults settings for the invision command:",
			Components: messageComponents,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func (b *botImpl) processInvisionDimensionSetting(s *discordgo.Session, i *discordgo.InteractionCreate, height, width int) {
	botSettings, err := b.invisionQueue.UpdateDefaultDimensions(width, height)
	if err != nil {
		log.Printf("error updating default dimensions: %v", err)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Error updating default dimensions...",
			},
		})
		if err != nil {
			log.Printf("Error responding to interaction: %v", err)
		}

		return
	}

	messageComponents := settingsMessageComponents(botSettings)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Choose defaults settings for the invision command:",
			Components: messageComponents,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func (b *botImpl) processInvisionBatchSetting(s *discordgo.Session, i *discordgo.InteractionCreate, batchCount, batchSize int) {
	botSettings, err := b.invisionQueue.UpdateDefaultBatch(batchCount, batchSize)
	if err != nil {
		log.Printf("error updating batch settings: %v", err)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Error updating batch settings...",
			},
		})
		if err != nil {
			log.Printf("Error responding to interaction: %v", err)
		}

		return
	}

	messageComponents := settingsMessageComponents(botSettings)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Choose defaults settings for the invision command:",
			Components: messageComponents,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}
