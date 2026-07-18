# Scrabble Scorer (Fastly Compute WASM)

A stateless Scrabble scoring application built with Go and running on Fastly Compute. This project allows 2-4 players to track their Scrabble scores, calculate word values with multipliers, and manage game history directly at the edge.

## Features

- **Stateless Architecture**: Game state is managed on the client side and passed to the WASM backend for processing, ensuring reliability in ephemeral environments.
- **Dynamic Scoring**: Automatically calculates Scrabble word scores based on standard letter values.
- **Multipliers Support**: Supports Double Letter (DL), Triple Letter (TL), Double Word (DW), and Triple Word (TW) multipliers.
- **Player Management**: Supports 2 to 4 players with individual score tracking and word history.
- **Tailwind CSS Frontend**: A clean, responsive UI built with Tailwind CSS.

## Getting Started

### Prerequisites

- [Go 1.23](https://go.dev/dl/) or later.
- [Fastly CLI](https://github.com/fastly/cli) for local development and deployment.

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/mealies/wasmScrabbleScorer.git
   cd wasmScrabbleScorer
   ```

2. Serve the application locally:
   ```bash
   fastly compute serve
   ```

3. Open your browser and navigate to `http://127.0.0.1:7676`.

## How to Use

1. **Start a Game**: Enter the names of 2 to 4 players on the initial screen and click "Start Game".
2. **Add Scores**:
   - Select a player.
   - Enter the word played.
   - (Optional) Assign multipliers to specific letters by selecting them from the dropdowns below the word input.
   - Alternatively, you can enter a raw numeric score if you've already calculated it.
   - Click "Add Word". You can add multiple words per round.
   - Once all words for the round are added, click "End Round" to finalize the score for that turn.
3. **Finish Game**:
   - Once the game is over, click "Finish Game".
   - Select the player who went out first (the winner).
   - Enter the remaining point values for the other players.
   - Click "Submit Final Scores" to see the final standings.

## Technical Details

This project is compiled to WebAssembly (WASM) and runs on Fastly's Compute platform. It uses the `fsthttp` package for handling HTTP requests and responses at the edge. The frontend is embedded directly into the WASM binary using Go's `embed` directive and served as a single-page application.

## Security issues

Please see our [SECURITY.md](SECURITY.md) for guidance on reporting security-related issues.
