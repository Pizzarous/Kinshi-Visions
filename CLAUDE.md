# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Kinshi-Visions is a Discord bot written in Go that integrates with the [Automatic1111 Stable Diffusion WebUI](https://github.com/AUTOMATIC1111/stable-diffusion-webui) API to generate AI images on demand from Discord slash commands.

## Build & Run

```bash
# Build
go build

# Run (requires .env to be configured)
./kinshi_vision_bot

# Run with flags
./kinshi_vision_bot -dev                        # Development mode (commands prefixed with "dev_")
./kinshi_vision_bot -invision <command_name>    # Override the slash command name
./kinshi_vision_bot -remove                     # Delete all registered commands on exit
```

## Environment Configuration

Copy `.env_example` to `.env` and fill in:
- `BOT_TOKEN` — Discord bot token
- `GUILD_ID` — Discord server ID
- `API_HOST` — Automatic1111 WebUI URL (e.g., `http://127.0.0.1:7860`, no trailing slash)

The Automatic1111 WebUI must be started with the `--api` flag. Use `--listen` too if it's on a different machine.

## Architecture

The application follows a layered architecture with interface-based dependency injection:

```
Discord Events (discordgo)
    → discord_bot/          — Slash command handling, button interactions
    → invision_queue/       — Async request queue, processes generations
    → stable_diffusion_api/ — HTTP client for the Automatic1111 API
    → repositories/         — Data access (image_generations, default_settings)
    → databases/sqlite/     — SQLite with sequential auto-migrations
```

**Initialization order in `main.go`:** SD API client → SQLite DB (auto-migrates) → repositories → queue → Discord bot → `bot.Start()` (blocks until Ctrl+C).

**Key packages:**
- `invision_queue/queue.go` (~1282 lines) — core generation logic: queues requests, polls SD API progress, builds responses, handles variations/upscale/reroll
- `discord_bot/discord_bot.go` (~710 lines) — registers slash commands, dispatches interaction events to the queue
- `stable_diffusion_api/` — thin HTTP wrapper for `/sdapi/v1/txt2img`, `/sdapi/v1/upscaler`, and progress endpoints
- `composite_renderer/` — assembles individual generated images into a 2×2 grid for Discord display
- `entities/` — `ImageGeneration` and `DefaultSettings` structs (SD parameter containers, not ORM models)

**Database:** SQLite file `sd_discord_bot.sqlite` created on first run. Migrations in `databases/sqlite/sqlite.go` run sequentially at startup.

## Prompt Parameters

Parameters parsed from the prompt string in `invision_queue/queue.go`:
- `--ar <w>:<h>` — aspect ratio
- `--step <n>` — sampling steps
- `--cfgscale <f>` — CFG scale
- `--seed <n>` — fixed seed
- `--zoom <f>` — hires.fix zoom rate (enables hires.fix when set)
- `--px <x>,<y>` — explicit pixel dimensions

## CI/CD

- `.github/workflows/golangci-lint.yaml` — runs golangci-lint on PRs
- `.github/workflows/release.yml` — builds cross-platform binaries (Linux/Windows/Darwin, amd64/arm64) on tag push
