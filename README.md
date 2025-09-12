# babbr

[![Actions Status](https://github.com/hagatasdelus/babbr/workflows/CI/gadge.svg)](https://github.com/hagatasdelus/babbr/actions)

Fish shell-style abbreviations for bash, providing an experience inspired by fish's abbr functionality.

## Installation

```bash
go install github.com/hagatasdelus/babbr@latest
```

Execute the following line or add it to `$HOME/.bashrc`:

```bash
eval "$(babbr init)"
```

## Usage

> **Note**
> A minor visual flicker may occur during abbreviation expansion. This is a known limitation of the `bind -x` macro within Bash's readline library.  This behavior does not impact the functionality of the tool.

### Configuration

Configuration can be done with a `config.yaml` file.
You can see all currently defined abbreviations by running.

```bash
babbr list
```

### Example for config.yaml

Location of config.yaml is:

- UNIX: `$XDG_CONFIG_HOME/babbr/config.yaml` or `$HOME/.config/babbr/config.yaml`
- Windows: `%APPDATA%\babbr\config.yaml`

```yaml
abbreviations:
  # I. Basic abbreviations (alias usage)
  - name: list files with details
    abbr: l
    snippet: ls -l

  - name: list all files with details
    abbr: la
    snippet: ls -la

  - name: clear screen
    abbr: cls
    snippet: clear

  # II. Global abbreviation (expanded anywhere on command line)
  - name: Redirect stdout and stderr to /dev/null
    abbr: "null"
    snippet: ">/dev/null 2>&1"
    options:
      # `position: anywhere` allows expansion anywhere in the row
      position: anywhere

  - name: Pipe to less
    abbr: L
    snippet: "| less"
    options:
      position: anywhere

  # III. Context-aware abbreviations
  - name: git status
    abbr: s
    snippet: status
    options:
      position: anywhere
      # By `command: git`, it is expanded only as an argument to the `git` command
      command: git

  - name: git commit
    abbr: c
    snippet: commit
    options:
      position: anywhere
      command: git

  - name: git commit with message
    abbr: cm
    snippet: "commit -m '%'"
    options:
      position: anywhere
      command: git
      # 'set_cursor: true' will move the cursor to the `%` position after expansion
      set_cursor: true

  - name: git add all
    abbr: a
    snippet: "add ."
    options:
      position: anywhere
      command: git

  - name: git push current branch to origin (dynamic)
    abbr: po
    snippet: "push origin $(git symbolic-ref --short HEAD)"
    options:
      position: anywhere
      command: git
      # `evaluate: true` causes the command in the snippet to be executed and its standard output to be expanded
      evaluate: true

  # IV. Conditional and environmentally-aware abbreviations
  - name: Use trash-cli if available
    abbr: del
    snippet: trash
    options:
      # This abbreviation is valid only if `condition` is true (the trash command is present)
      # The string is executed in a shell and the condition is determined by its success or failure (exit status)
      condition: "command -v trash &> /dev/null"

  - name: Fallback to interactive rm if trash-cli is not installed
    abbr: del
    snippet: "rm -i"

  - name: Open file with VSCode if code command exists
    abbr: code
    snippet: "code ."
    options:
      condition: "type code &> /dev/null"

  # V. Suffix alias (execution by regular expression pattern)
  - name: Execute Python scripts with filename
    # The `$file` variable in the snippet is replaced by the value captured in options.regex
    snippet: "python $file"
    options:
      # Matches words ending with `.py` and captures the entire word as `file`
      regex: '^(?<file>\S+\.py)$'
      evaluate: true

  - name: Execute shell scripts
    snippet: "bash $script"
    options:
      regex: '^(?<script>\S+\.sh)$'
      evaluate: true

  - name: Unpack tar.gz files
    snippet: "tar -xzvf $archive"
    options:
      regex: '^(?<archive>\S+\.tar\.gz)$'
      evaluate: true
```

## Configuration Reference

The following keys can be used in your `config.yaml` to define an abbreviation.

  * `name`: A brief, human-readable description of what the abbreviation does.
  * `abbr`: The short word to be typed that will trigger the expansion.
  * `snippet`: The text that will replace the `abbr` when it is expanded.
  * `options`: A block of advanced configuration keys for more specific behavior.
      * `position`: Controls where the abbreviation can be expanded.
          * `command` (default): Only expands when the `abbr` is the first word on the command line.
          * `anywhere`: Allows expansion at any point on the command line.
      * `command`: Restricts the abbreviation to only expand when it appears as an argument to a specific command.
      * `condition`: A shell command that is executed to determine if the abbreviation should be active. The abbreviation is only enabled if the command exits with a status of 0 (success).
      * `regex`: A regular expression that serves as the trigger instead of a fixed `abbr` string. This is useful for creating "suffix aliases" that act on patterns, like file extensions.
      * `evaluate`: Enables dynamic features in the `snippet`.
          * When used with `regex`, it substitutes named capture groups (e.g., `(?<name>...)`) from the regex into variables in the snippet (e.g., `$name`).
          * When used without `regex`, it allows for shell command substitution (e.g., `$(...))`) within the snippet to be evaluated by the shell.
      * `set_cursor`: A boolean (`true`/`false`). If true, the cursor will be moved to the position of the `%` character in the snippet after expansion. The `%` character is then removed.

## Features

This tool brings a fish shell-inspired `abbr` experience to bash, including key features like:

| Feature                  | Description                                                     |
| ------------------------ | --------------------------------------------------------------- |
| **Inline expansion** | Abbreviations expand in the command line before you run the command. |
| **Visibility** | You can see and edit the full command before execution. |
| **History integrity** | Your shell history saves the expanded commands, not the abbreviations. |
| **Position control** | You can restrict abbreviations to the command position or allow them anywhere on the line. |
| **Command-specific** | Abbreviations can be set to expand only as arguments for a specific command (e.g., `git`). |
| **Conditional expansion** | You can enable abbreviations based on shell conditions, such as checking if a command exists. |
| **Cursor positioning** | You can set the cursor's position after an expansion using a `%` marker in the snippet. |
| **Regex-powered snippets** | You can use regular expressions to trigger abbreviations, primarily to substitute named capture groups. |

## License

MIT

## Author

Hagata
