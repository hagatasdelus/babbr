# babbr

[![Actions Status](https://github.com/hagatasdelus/babbr/workflows/CI/gadge.svg)](https://github.com/hagatasdelus/babbr/actions)

Fish shell-style abbreviations for bash.
A tool that provides fish shell abbr functionality for bash.

## Installation

Build from source:

```bash
go install github.com/hagatasdelus/babbr@latest
```

Execute the following line or add it to `$HOME/.bashrc`:

```bash
eval "$(babbr init)"
```

## Usage

### Configuration

Configuration can be done with a `config.yaml` file.

### Example for config.yaml

Location of config.yaml is:

- UNIX: `$XDG_CONFIG_HOME/babbr/config.yaml` or `$HOME/.config/babbr/config.yaml`

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

> **Note**
> A minor visual flicker may occur during abbreviation expansion. This is a known limitation of the bind -x macro within Bash's readline library.  This behavior does not impact the functionality of the tool.

## Features

This tool reproduces fish shell's abbr functionality completely:

- **Inline expansion**: Abbreviations expand in the command line before execution
- **Visibility**: Users can see and edit the full command before execution
- **History integrity**: Full commands are saved in bash history, not abbreviations
- **Position control**: Abbreviations can be restricted to command position or allowed anywhere
- **Command-specific**: Abbreviations can be limited to specific parent commands
- **Conditional expansion**: Abbreviations can be enabled based on shell conditions (e.g., command existence)
- **Regex support**: Pattern-based abbreviations using regular expressions
- **Variable substitution**: Substitute named capture groups from regex patterns into snippets
- **Cursor positioning**: Set cursor position after expansion using the `%` marker

## License

MIT

## Author

Hagata
