package main

import (
	"context"
	"flag"
	"kinshi_vision_bot/databases/sqlite"
	"kinshi_vision_bot/discord_bot"
	"kinshi_vision_bot/invision_queue"
	"kinshi_vision_bot/repositories/default_settings"
	"kinshi_vision_bot/repositories/image_generations"
	"kinshi_vision_bot/stable_diffusion_api"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// getEnvVar retrieves the environment variable with the provided key and a default value.
func getEnvVar(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

var (
	invisionCommand    = flag.String("invision", "invision", "Invision command name. Default is \"invision\"")
	removeCommandsFlag = flag.Bool("remove", false, "Delete all commands when bot exits")
	devModeFlag        = flag.Bool("dev", false, "Start in development mode, using \"dev_\" prefixed commands instead")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.Parse()

	guildID := getEnvVar("GUILD_ID", "")
	botToken := getEnvVar("BOT_TOKEN", "")
	apiHost := getEnvVar("API_HOST", "")

	if guildID == "" {
		log.Fatal("Guild ID is required")
	}

	if botToken == "" {
		log.Fatal("Bot token is required")
	}

	if apiHost == "" {
		log.Fatal("API host is required")
	}

	if invisionCommand == nil || *invisionCommand == "" {
		log.Fatalf("Invision command flag is required")
	}

	devMode := false

	if devModeFlag != nil && *devModeFlag {
		devMode = *devModeFlag

		log.Printf("Starting in development mode.. all commands prefixed with \"dev_\"")
	}

	removeCommands := false

	if removeCommandsFlag != nil && *removeCommandsFlag {
		removeCommands = *removeCommandsFlag
	}

	stableDiffusionAPI, err := stable_diffusion_api.New(stable_diffusion_api.Config{
		Host: apiHost,
	})
	if err != nil {
		log.Fatalf("Failed to create Stable Diffusion API: %v", err)
	}

	ctx := context.Background()

	sqliteDB, err := sqlite.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create sqlite database: %v", err)
	}

	generationRepo, err := image_generations.NewRepository(&image_generations.Config{DB: sqliteDB})
	if err != nil {
		log.Fatalf("Failed to create image generation repository: %v", err)
	}

	defaultSettingsRepo, err := default_settings.NewRepository(&default_settings.Config{DB: sqliteDB})
	if err != nil {
		log.Fatalf("Failed to create default settings repository: %v", err)
	}

	invisionQueue, err := invision_queue.New(invision_queue.Config{
		StableDiffusionAPI:  stableDiffusionAPI,
		ImageGenerationRepo: generationRepo,
		DefaultSettingsRepo: defaultSettingsRepo,
	})
	if err != nil {
		log.Fatalf("Failed to create invision queue: %v", err)
	}

	bot, err := discord_bot.New(discord_bot.Config{
		DevelopmentMode: devMode,
		BotToken:        botToken,
		GuildID:         guildID,
		InvisionQueue:   invisionQueue,
		InvisionCommand: *invisionCommand,
		RemoveCommands:  removeCommands,
	})
	if err != nil {
		log.Fatalf("Error creating Discord bot: %v", err)
	}

	bot.Start()

	log.Println("Gracefully shutting down.")
}
