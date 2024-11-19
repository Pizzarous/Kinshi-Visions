# Stable Diffusion Discord Bot

This Discord bot interacts with the Automatic1111 API, part of the [Stable Diffusion WebUI project](https://github.com/AUTOMATIC1111/stable-diffusion-webui).

Watch a demo of the current features [here](https://www.youtube.com/watch?v=of5MBh3ueMk).

---

## Building

1. Clone this repository.
2. Install Go using the [official installer](https://golang.org/dl/).
3. Build the bot using:

   ```bash
   go build
   ```

---

## Environment Configuration

Before running the bot, configure an `.env` file to specify your bot’s details and the Automatic1111 instance.

### Steps

1. Duplicate the `.env_example` file in the project directory and rename it to `.env`.
2. Open the `.env` file in a text editor and fill in the required fields:
   - `BOT_TOKEN` — Your Discord bot token.
   - `GUILD_ID` — The Guild ID (server ID) for your Discord server.
   - `API_HOST` — The URL of the Automatic1111 API instance.

   **Important Notes for `API_HOST`:**
   - If the Automatic1111 WebUI is running on the same computer as the bot, use `http://127.0.0.1:7860`.
   - If running on a different machine, replace `127.0.0.1` with the host's IP address (e.g., `http://192.168.1.100:7860`).
   - Do not include a trailing slash in the URL (e.g., use `http://192.168.1.100:7860`, not `http://192.168.1.100:7860/`).

---

## Usage

1. Create a Discord bot via the [Discord Developer Portal](https://discord.com/developers/applications) and obtain its token.
2. Add the bot to your Discord server with all necessary text permissions.
3. Start the Automatic1111 WebUI with the `--api` flag (use `--listen` as well if it's on a different machine).
4. Run the bot executable created during the build process:

   ```bash
   ./kinshi_vision_bot
   ```

   Ensure the `.env` file is correctly configured (see the [Environment Configuration](#environment-configuration) section).

### Notes

- On the first run, a SQLite database file will be generated in the same directory as the bot.
- Use the `-invision <new command name>` option to avoid conflicts with other bots.

---

## Commands

### `/invision_settings`

Displays buttons in Discord to update default settings for the `/invision` command.  

![Invision Settings](https://user-images.githubusercontent.com/7525989/211077599-482536ef-1a70-4f58-abf0-314c773c64c6.png)

### `/invision`

Generates an image based on a text prompt. Example:

```bash
/invision cute kitten riding a skateboard
```

#### Options

- Specify Aspect Ratio:

  ```bash
  /invision cute kitten --ar 16:9
  ```

---

## How it Works

The bot operates as a queue-based system, processing user requests and sending them to the Automatic1111 WebUI API. It supports:

- Image generation.
- Interaction updates (e.g., re-rolling, variations, up-scaling).

All image data is logged locally in a SQLite database.  

![Bot Workflow](https://user-images.githubusercontent.com/7525989/209247280-4318a73a-71f4-48aa-8310-7fdfbbbf6820.png)

---

## Fork Information

This bot includes updates from [pitapan5376](https://github.com/pitapan5376/stable-diffusion-discord-bot).

### Recent Changes

#### 01. Aspect Ratio Support

Added the `--ar` parameter to set aspect ratios without upscaling:

- `1girl --ar 4:3`  
  ![Sample for AR 4:3](https://github.com/pizzarous/kinshi-visions/blob/master/document/003_aspect_ratio_4_3.png?raw=true)
- `1girl --ar 1:2`  
  ![Sample for AR 1:2](https://github.com/pizzarous/kinshi-visions/blob/master/document/003_aspect_ratio_1_2.png?raw=true)

#### 02. Sampling Steps

Introduced `--step` to control generation steps:

- `--step 7` (512x512)  
  ![Sample for Step 7](https://github.com/pizzarous/kinshi-visions/blob/master/document/004_steps_param_7.png?raw=true)

- `--step 50 --ar 2:1` (1024x768)  
  ![Sample for Step 50](https://github.com/pizzarous/kinshi-visions/blob/master/document/004_steps_param_50.png?raw=true)

#### 03. CFG Scale

Added `--cfgscale` to adjust the CFG scale:

- `--cfgscale 1.2`  
  ![Sample for CFG Scale Low](https://github.com/pizzarous/kinshi-visions/blob/master/document/005_cfg_scale_1.png?raw=true)

- `--cfgscale 15.3`  
  ![Sample for CFG Scale High](https://github.com/pizzarous/kinshi-visions/blob/master/document/005_cfg_scale_15.png?raw=true)

#### 04. Seed Parameter

Use `--seed` to specify a seed value for deterministic results:

- `--seed 111`  
  ![Sample for Seed](https://github.com/pizzarous/kinshi-visions/blob/master/document/006_seed.png?raw=true)

#### 05. Negative Prompts

Support for negative prompts via `negative_prompt` parameter.  
![Negative Prompt Param](https://github.com/pizzarous/kinshi-visions/blob/master/document/007_negative_prompt.png?raw=true)
  
![Seed on Bigint](https://github.com/pizzarous/kinshi-visions/blob/master/document/008_seed_bigint.png?raw=true)

#### 06. Sampler Selection

Added a pop-up to choose the sampler during image generation.  
![Sampler Choice](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_selection.png?raw=true)

- DPM++ S2 a Karras  
  ![Sampler: DPM++ S2 a Karras](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_DPMppS2aKarras.png?raw=true)

- DPM Adaptive  
  ![Sampler: DPM Adaptive](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_DPMAdaptive.png?raw=true)

- UniPC  
  ![Sampler: UniPC](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_UniPC.png?raw=true)

#### 07. Hires.fix

Partial support for `hires.fix`. Added `hr_scale` and `hr_upscaler` to the table.  
![Hiresfix1](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix1.png?raw=true)

#### 08. Hires.fix with Zoom Rate

Added the `--zoom` parameter to switch `hires.fix` on/off with a specified zoom rate.
  
![Hiresfix2](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix2.png?raw=true)

- Hires.fix ON with Zoom 1.2  
  ![Hiresfix3](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix3.png?raw=true)

- Hires.fix OFF  
  ![Hiresfix4](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix4.png?raw=true)

---
