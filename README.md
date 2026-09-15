> [!NOTE]
> **This is a test fixture, not the real project.**
> `dispatrick/clidle-test` is a clone of [`ajeetdsouza/clidle`](https://github.com/ajeetdsouza/clidle)
> used purely as a sandbox for testing [Dispatch](https://upsun.com). It exists so we
> have a small, real-feeling Go codebase to open throwaway pull requests against.
> It is not maintained, not deployed, and not affiliated with the upstream author.
> All credit for the actual game goes to the original project (MIT licensed — see `LICENSE`).

<div align="center">

# clidle

**Wordle, now over SSH.**

<img align="right" alt="Preview" width="50%" src="preview.png" />

**Try it:**

```sh
ssh clidle.duckdns.org -p 3000
```

**Or, run it locally:**

```sh
go install github.com/ajeetdsouza/clidle@latest
```

</div>

## Play over SSH

clidle can serve itself over SSH, so anyone with an SSH client can play without
installing anything. Start the server:

```sh
clidle serve          # or: clidle --ssh
```

Then connect from another terminal:

```sh
ssh localhost -p 23234
```

Each connection gets its own game, rendered with the colors and window size of
the connecting terminal. Sessions without a terminal (for example `ssh -T`) are
rejected.

### Configuration

| Flag         | Environment variable | Default                              | Description                                     |
| ------------ | -------------------- | ------------------------------------ | ----------------------------------------------- |
| `--host`     | `SSH_HOST`           | `127.0.0.1`                          | Address to bind to                              |
| `--port`     | `SSH_PORT`           | `23234`                              | Port to listen on                               |
| `--host-key` | `SSH_HOST_KEY_PATH`  | `$CLIDLE_DATA_DIR/hostkey`           | SSH host key, generated on first run if missing  |

Flags take precedence over environment variables, which take precedence over the
defaults. `--serve 0.0.0.0:1337` remains supported and sets host and port at
once.

To expose the server beyond the local machine, bind it to all interfaces:

```sh
SSH_HOST=0.0.0.0 clidle serve
```

> [!WARNING]
> The server accepts **any** public key and does not authenticate players, so
> exposing it publicly makes the game playable by anyone who can reach the port.

## How to play

You have 6 attempts to guess the correct word. Each guess must be a valid 5 letter
word.

After submitting a guess, the letters will turn green, yellow, or gray.

- **Green:** The letter is correct, and is in the correct position.
- **Yellow:** The letter is present in the solution, but is in the wrong position.
- **Gray:** The letter is not present in the solution.

## Scoring

Your final score is based on how many guesses it took to arrive at the solution:

| Guesses | Score |
| ------- | ----- |
| 1       | 100   |
| 2       | 90    |
| 3       | 80    |
| 4       | 70    |
| 5       | 60    |
| 6       | 50    |
