# Hash

Hash is a terminal based password manager for everyone. If you spend most of your time in your terminal, Hash may be a good fit in your workflow.

Watch this video to understand the functionality of Hash.

<video width="auto" height="auto" controls>
  <source src="hash-v1-trailer.mov" type="video/mp4">
</video>

## Motivation for this project

I always prefer email password based signups in websites. But managing password is a hassle for me. I have used [Bitwarden password manager](https://bitwarden.com/) and it is great. But I want to access my passwords from my terminal. I don't want to click around with the mouse to get a single password. Therefore I created Hash.

## Getting started

1. Go to the [GitHub release section](https://github.com/hrishiksh/hash/releases) and pick the binary for your operating system. All the binaries are statically linked. Hence, binaries are portable and don't depend on OS libraries.

2. Rename the binary file to `hash`.

3. Make sure that the downloaded binary is executable. If not, run `chmod +x hash`.

4. After downloading, place the binary in your path. e.g, For Linux users, put the binary inside `~/.local/bin` or `~/bin` directory. For macOS, put the binary inside `/usr/local/bin` directory.

5. Start using Hash by running `hash` in your terminal.

## Tech Stack

- This project is written in [Golang](https://go.dev/).
- The beautiful TUI (Terminal User Interface) is made using [BubbleTea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles) and [Lipgloss](https://github.com/charmbracelet/lipgloss).
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) is used for generating the Salt
- [PBKDF2](https://pkg.go.dev/golang.org/x/crypto/pbkdf2) is used to generate the Secret Key
- [Nacl Secretbox](https://pkg.go.dev/golang.org/x/crypto/nacl/secretbox) is used to symmetrically encrypt the password.
- All the password are [Hex encoded](https://pkg.go.dev/encoding/hex) before storing on the Database
- [Sqlite](https://github.com/mattn/go-sqlite3) is used to store the encrypted passwords.

## Build this project in your machine

1. Clone the repository

```bash
git clone https://github.com/hrishiksh/golang-oauth2-starter.git
```

2. Download all the dependency

```bash
go mod download
go mod tidy
```

3. Build binary

```bash
make build
```

**Note: ** As some part of this project uses CGO, native go cross compilation is not easy. To statically link the binary, I have used the Zig compiler. This is [a tutorial](https://hrishikeshpathak.com/tips/build-static-binary-cross-compile-cgo-project-zig-compiler/) on how to do so in your project.
