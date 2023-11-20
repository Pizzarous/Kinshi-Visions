# Stable Diffusion Discord Bot

This Discord bot interacts with the Automatic1111 API, part of the [Stable Diffusion WebUI project](https://github.com/AUTOMATIC1111/stable-diffusion-webui).

Watch a demo of the current features [here](https://www.youtube.com/watch?v=of5MBh3ueMk).

## Installation

1. Download the suitable version from the [releases page](https://github.com/AndBobsYourUncle/stable-diffusion-discord-bot/releases):
   - Windows: `windows-amd64`
   - Intel Macs: `darwin-amd64`
   - M1 Macs: `darwin-arm64`
   - Raspberry Pi: `linux-arm64`
   - Most other Linux devices: `linux-amd64`
2. Extract the archive to your preferred location.

## Building (Optional)

1. Clone this repository.
2. Install Go using the [official installer](https://golang.org/dl/).
3. Build the bot with `go build`.

## Usage

1. Create a Discord bot and obtain the token.
2. Add the bot to your Discord server with all text permissions.
3. Ensure the Automatic1111 webui is running with `--api` (and `--listen` if on a different machine).
4. Run the bot with `./stable_diffusion_bot -token <token> -guild <guild ID> -host <webui host, e.g., http://127.0.0.1:7860>`.
   - Ensure `-host` matches the A1111's IP address. Use `127.0.0.1` if on the same computer.
   - No trailing slash after the port number (e.g., `http://127.0.0.1:7860`, not `http://127.0.0.1:7860/`).
5. The first run generates a new SQLite DB file in the current directory.

Use `-invision <new command name>` to avoid conflicts with other bots.

## Commands

### `/invision_settings`

Responds with a message containing buttons to update default settings for `/invision`.

![Invision Settings](https://user-images.githubusercontent.com/7525989/211077599-482536ef-1a70-4f58-abf0-314c773c64c6.png)

### `/invision`

Creates an image from a text prompt (e.g., `/invision cute kitten riding a skateboard`).

Options:
- Aspect Ratio: `/invision cute kitten --ar 16:9`

## How it Works

The bot uses a FIFO queue to process user commands. It sends interactions to the Automatic1111 WebUI API, updates messages, and handles buttons for re-rolling, variations, and up-scaling.

All image data is stored in a local SQLite database.

![Bot Workflow](https://user-images.githubusercontent.com/7525989/209247280-4318a73a-71f4-48aa-8310-7fdfbbbf6820.png)

## Contributing

Pull requests are welcome. Major changes? Open an issue first.

Potential Features:
- [x] Move defaults to the database
- [ ] Per-user defaults/settings, usage limits
- [x] Re-roll image
- [x] Generate multiple images
- [x] Upscale images
- [x] Generate variations on a grid
- [ ] More settings for `/invision` command
- [ ] Image to image processing

## Fork Information

Implemented changes from [pitapan5376 fork](https://github.com/pitapan5376/stable-diffusion-discord-bot) on April 8, 2023, at 07:10:00 (JST).

### Changes:

#### 001. Button Order and Icon Captions

Reordered buttons for better alignment on iPhone Discord.

#### 002. Changed Fonts

![Change Prompt Font to Monospace](https://github.com/pizzarous/kinshi-visions/blob/master/document/002_change_prompt_font.png?raw=true)

#### 003. Enable Aspect Ratio (Without Upscaler)

Parsed the `--ar` parameter and computed new values, allowing for different aspect ratios without upscaling. Examples:

- `1girl --ar 4:3`
  ![Sample for AR 4:3](https://github.com/pizzarous/kinshi-visions/blob/master/document/003_aspect_ratio_4_3.png?raw=true)

- `1girl --ar 1:2`
  ![Sample for AR 1:2](https://github.com/pizzarous/kinshi-visions/blob/master/document/003_aspect_ratio_1_2.png?raw=true)

#### 004. Add Sampling Steps

Introduced a `--step` parameter to control the number of steps processed during image generation.

- `--step 7` (512x512)
  ![Sample for Step 7](https://github.com/pizzarous/kinshi-visions/blob/master/document/004_steps_param_7.png?raw=true)

- `--step 50 --ar 2:1` (1024x768)
  ![Sample for Step 50](https://github.com/pizzarous/kinshi-visions/blob/master/document/004_steps_param_50.png?raw=true)

#### 005. Add CFG Scale Parameter

Introduced the `--cfgscale` parameter to control CFG scale values.

- `--cfgscale 1.2`
  ![Sample for CFG Scale Low](https://github.com/pizzarous/kinshi-visions/blob/master/document/005_cfg_scale_1.png?raw=true)

- `--cfgscale 15.3`
  ![Sample for CFG Scale High](https://github.com/pizzarous/kinshi-visions/blob/master/document/005_cfg_scale_15.png?raw=true)

#### 006. Seed Parameter

Added a `--seed` parameter to specify the seed value for image generation.

- `--seed 111`
  ![Sample for Seed](https://github.com/pizzarous/kinshi-visions/blob/master/document/006_seed.png?raw=true)

#### 007. Negative Prompt Parameter

Introduced a `negative_prompt` parameter to provide a negative prompt for image generation.

![Negative Prompt Param](https://github.com/pizzarous/kinshi-visions/blob/master/document/007_negative_prompt.png?raw=true)

#### 008. Bugfix: Seed Value for Big Int

Fixed an issue where the bot crashed when receiving large seed values.

![Seed on Bigint](https://github.com/pizzarous/kinshi-visions/blob/master/document/008_seed_bigint.png?raw=true)

#### 009. Selection Pop-up for Sampler

Added a pop-up for selecting a sampler for image generation.

![Sampler Choice](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_selection.png?raw=true)

Sample Samplers:
- DPM++ S2 a Karras
  ![Sampler: DPM++ S2 a Karras](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_DPMppS2aKarras.png?raw=true)

- DPM Adaptive
  ![Sampler: DPM Adaptive](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_DPMAdaptive.png?raw=true)

- UniPC
  ![Sampler: UniPC](https://github.com/pizzarous/kinshi-visions/blob/master/document/009_sampler_UniPC.png?raw=true)

#### 010. Hires.fix

Partial support for `hires.fix`. Added `hr_scale` and `hr_upscaler` to the table.

![Hiresfix1](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix1.png?raw=true)

#### 011. Bugfix: NegativePrompt

Fixed an issue where the negative prompt wasn't applied.

#### 012. Hires.fix with Zoom Rate

Added the `--zoom` parameter to switch `hires.fix` on/off with a specified zoom rate.

![Hiresfix2](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix2.png?raw=true)

- Hires.fix ON with Zoom 1.2
  ![Hiresfix3](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix3.png?raw=true)

- Hires.fix OFF
  ![Hiresfix4](https://github.com/pizzarous/kinshi-visions/blob/master/document/012_hiresfix4.png?raw=true)

#### 013. Apply Upstream Update

Included updates from the upstream repository.

Please note that the content has been reviewed, and non-English elements have been removed. If you have any specific questions or need further clarification on any point, feel free to ask.
