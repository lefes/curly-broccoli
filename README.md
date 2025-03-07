# Curly Broccoli Discord Bot

A fun and interactive Discord bot written in Go that responds to various triggers, provides weather information, and entertains server members with games and random responses.

## Features

### Basic Interactions
- Responds to greetings and farewells
- Reacts to morning/evening messages with appropriate emojis
- Adds reactions to messages about specific topics (Legion, Coven, games, etc.)
- Responds to various trigger words with custom messages

### Commands
- `!пиво <число>` - Beer drinking challenge with success/failure outcomes
- `!голосование` - Creates a poll to determine "who is a dick today"
- `!пенис` - Generates a random penis size visualization
- `!бубс` - Generates a random breast size visualization
- `!медведь` - Calculates your chance to defeat a bear
- `!ролл` or `!d20` - Rolls a D20 die
- `!писька` - Rates a user on their "писька" level
- `!письки` - Rates multiple users' "писька" levels
- `!гей` - Calculates and shows a user's gay percentage
- `!анекдот` - Fetches a random joke
- `!шар` - Magic 8-ball style responses to questions

### Weather System
- `!weather <city>` - Shows current weather for a location
- `!weather <city> <days>` - Shows weather forecast for specified number of days
- `!погода <город>` - Russian version of the weather command
- Supports city shortcuts for common Russian cities

## Installation

### Prerequisites
- Go 1.18 or higher
- Discord Bot Token
- Visual Crossing Weather API Key

### Setup
1. Clone the repository:
   ```
   git clone https://github.com/lefes/curly-broccoli.git
   cd curly-broccoli
   ```

2. Create a `.env` file in the root directory with your Discord token and weather API key:
   ```
   TOKEN=your_discord_bot_token
   WEATHER_API_KEY=your_weather_api_key
   ```

3. Build the application:
   ```
   go build
   ```

4. Run the bot:
   ```
   ./curly-broccoli
   ```

## Project Structure

```
├── consts.go           # Constants and message arrays
├── jokes
│   └── jokes.go        # Joke fetching functionality
├── main.go             # Main bot logic and command handlers
├── pkg
│   ├── logging         # Logging utilities
│   └── weather         # Weather API integration
```

## Weather API

The weather functionality uses the Visual Crossing Weather API. You can get a free API key by signing up at [Visual Crossing](https://www.visualcrossing.com/).

## Customization

- Add or modify responses in the constant arrays in `consts.go`
- Add new command handlers in `main.go`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [DiscordGo](https://github.com/bwmarrin/discordgo) for Discord API integration
- [Visual Crossing](https://www.visualcrossing.com/) for weather data
- [goquery](https://github.com/PuerkitoBio/goquery) for HTML parsing