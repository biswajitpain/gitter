# gitter

`gitter` is a command-line wrapper for Git that enhances your commit workflow with an intelligent "commit review" (`cr`) command, including optional LLM-powered commit message generation. It acts as a drop-in replacement for your daily `git` commands while providing powerful automation for creating well-structured commit messages.

## Features

-   **Git Command Passthrough**: Use `gitter` just like `git` for all standard commands (e.g., `gitter status`, `gitter pull`). All arguments are passed directly to the underlying `git` command.
-   **`cr` (Commit Review) Command**: A custom command that streamlines the commit process:
    -   Stages all current changes (`git add .`).
    -   Generates a comprehensive diff of staged changes.
    -   Prompts you for a brief, high-level description of your changes.
    -   Generates a detailed commit message based on the diff and your input (either via LLM or a structured template).
    -   Asks for your confirmation before committing.
    -   Offers to unstage changes if the commit is cancelled.
-   **LLM-Powered Commit Message Generation**: Integrate with Large Language Models (LLMs) like OpenAI to generate high-quality, conventional commit messages.
-   **Configurable**: Easily set up your preferred LLM provider and API key.

## Installation

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.23 or higher)
-   [Git](https://git-scm.com/downloads)

### Building from Source

1.  **Clone the repository (if applicable):**
    ```bash
    git clone <repository_url>
    cd gitter
    ```
    *(Assuming you are in the `gitter` project directory)*

2.  **Tidy modules:**
    ```bash
    go mod tidy
    ```

3.  **Build the executable:**
    To build a simple executable, run:
    ```bash
    go build -o gitter .
    ```

    To build the executable with embedded version information, use the following command:
    ```bash
    go build -o gitter -ldflags="-X 'gitter/cmd.version=$(git describe --tags --always --dirty)'" .
    ```

### Making it Executable and Accessible

To use `gitter` from any directory, you should move the executable to a directory included in your system's `PATH`.

```bash
chmod +x gitter
sudo mv gitter /usr/local/bin/ # Adjust path based on your system's PATH
```

## Usage

### Versioning

To check the version of your `gitter` build, use the `version` command.

```bash
gitter version
```

If you built the application with the version information embedded, this will display the git tag, commit hash, and a `-dirty` suffix if you have uncommitted changes. Otherwise, it will show `dev`.

### Basic Git Commands

You can use `gitter` as a direct replacement for `git` for most commands. Any command not recognized as a `gitter` specific command will be passed directly to `git`.

```bash
gitter status
gitter log --oneline
gitter push origin main
gitter checkout -b feature/new-feature
```

### The `cr` Command (Commit Review)

The `cr` command automates the process of creating well-formed commit messages.

```bash
gitter cr
```

**Workflow:**

1.  `gitter` will automatically stage all your current changes (`git add .`).
2.  It will then generate a diff of these staged changes.
3.  You will be prompted to enter a brief, high-level description of your changes. This serves as a hint for the commit message generation.
4.  `gitter` will then generate a more detailed commit message based on the diff and your input. If an LLM is configured, it will attempt to use it; otherwise, it will use a structured template.
5.  The generated message will be displayed, and you'll be asked to confirm it.
6.  If you confirm (`y`), the changes will be committed with the generated message.
7.  If you cancel (`n`), the commit will be aborted, and you'll be given the option to unstage your changes.

### Configuring LLM Integration

To enable AI-powered commit message generation, you need to configure your LLM provider and API key.

**Command:**

```bash
gitter config --provider <provider_name> --api-key <your_api_key>
```

**Example for OpenAI:**

```bash
gitter config --provider openai --api-key "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

-   Replace `"sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxx"` with your actual OpenAI API key.
-   This command creates a configuration file at `~/.config/gitter/config.json` (or creates the directory if it doesn't exist) and stores your provider and API key. The file permissions are set to `0600` for security.

**LLM Fallback:**

-   If you have not configured an LLM provider, the `gitter cr` command will automatically fall back to using its simple, template-based message generator.
-   If an LLM is configured but fails to generate a message (e.g., due to network issues, invalid API key, or API errors), `gitter` will print a warning and gracefully fall back to the simple generator, ensuring your commit workflow is not interrupted.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
