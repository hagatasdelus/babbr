# babbr

Abbreviation extension plugin for bash that completely reproduces fish shell's abbr functionality.

## Features

- **Basic abbreviations**: Simple alias-like expansions
- **Global abbreviations**: Expand anywhere on the command line with `position: anywhere`  
- **Context-aware abbreviations**: Expand only when used with specific commands (e.g., git)
- **Regex-based abbreviations**: Pattern matching for file extensions and more
- **Conditional abbreviations**: Expand only when certain conditions are met
- **Dynamic evaluation**: Execute commands and insert their output
- **Cursor positioning**: Place cursor at specific positions after expansion
- **High performance**: Caching and optimized processing
- **TDD development**: Comprehensive test coverage for reliability

## Installation

```bash
go install github.com/hagatasdelus/babbr@latest
```

Add the following line to your `~/.bashrc`:

```bash
eval "$(babbr init)"
```

## Usage

```
Usage: babbr [OPTIONS] [COMMAND]

Commands:
  init      Initialize the plugin for the current shell session
  list      List all configured abbreviations
  expand    (Internal) Expand an abbreviation based on buffer content
  help      Print this message or the help of the given subcommand(s)

Options:
  -v, --version  Print version information
  -h, --help     Print help

Configuration:
  Config file: $XDG_CONFIG_HOME/babbr/config.yaml or $HOME/.config/babbr/config.yaml

Examples:
  eval "$(babbr init)"     # Add to your ~/.bashrc
  babbr list              # Show configured abbreviations
```

### Configuration

Configuration can be done with a `config.yaml` file.

### Example for config.yaml

Location of config.yaml is:

* UNIX: `$XDG_CONFIG_HOME/babbr/config.yaml` or `$HOME/.config/babbr/config.yaml`

```yaml
abbreviations:
  # --------------------------------------------------------------------
  # I. Basic abbreviations (alias usage)
  # --------------------------------------------------------------------
  - name: list files with details
    abbr: l
    snippet: ls -l
    # Condition: Command part
    # Input: `l<Space>` or `l<CR>`
    # Expansion: `ls -l`

  - name: list all files with details
    abbr: la
    snippet: ls -la
    # Condition: Command part
    # Input: `la<Space>` or `la<CR>`
    # Expansion: `ls -la`

  - name: clear screen
    abbr: cls
    snippet: clear
    # Condition: Command part
    # Input: `cls<Space>` or `cls<CR>`
    # Expansion: `clear`

  # --------------------------------------------------------------------
  # II. Global abbreviation (expanded anywhere on command line)
  # --------------------------------------------------------------------
  - name: Redirect stdout and stderr to /dev/null
    abbr: 'null'
    snippet: '>/dev/null 2>&1'
    options:
      # `position: anywhere` allows expansion anywhere in the row.
      position: anywhere
    # Condition: Any position on the command line
    # Input: `some_command null<Space>` or `some_command null<CR>`
    # Expansion: `some_command >/dev/null 2>&1`

  - name: Pipe to less
    abbr: L
    snippet: '| less'
    options:
      position: anywhere
    # Condition: Any position on the command line
    # Input: `cat long_file.txt L<Space>` or `cat long_file.txt L<CR>`
    # Expansion: `cat long_file.txt | less`

  # --------------------------------------------------------------------
  # III. Context-aware abbreviations (to speed up Git workflow)
  # --------------------------------------------------------------------
  - name: git status
    abbr: s
    snippet: status
    options:
      position: anywhere
      # By `command: git`, it is expanded only as an argument to the `git` command.
      command: git
    # Condition: Any position after `git` is typed.
    # Input: `git s<Space>` or `git s<CR>`
    # Expansion: `git status`

  - name: git commit
    abbr: c
    snippet: commit
    options:
      position: anywhere
      command: git
    # Condition: Any position after `git` is typed.
    # Input: `git c<Space>` or `git c<CR>`
    # Expansion: `git commit`

  - name: git commit with message
    abbr: cm
    snippet: "commit -m '%'"
    options:
      position: anywhere
      command: git
      # 'set_cursor: true' will move the cursor to the `%` position after expansion
      set_cursor: true
    # Condition: Any position after `git` is typed.
    # Input: `git cm<Space>` or `git cm<CR>`
    # Expansion: `git commit -m ''`  <-- cursor moves between ''

  - name: git add all
    abbr: a
    snippet: 'add .'
    options:
      position: anywhere
      command: git
    # Condition: Any position after `git` is typed.
    # Input: `git a<Space>` or `git a<CR>`
    # Expansion: `git add .`

  - name: git push current branch to origin (dynamic)
    abbr: po
    snippet: 'push origin $(git symbolic-ref --short HEAD)'
    options:
      position: anywhere
      command: git
      # `evaluate: true` causes the command in the snippet to be executed and its standard output to be expanded.
      # The interpretation of `evaluate` here is “insert a string that the shell will evaluate later”.
      evaluate: true
    # Condition: Any position after `git` is typed.
    # Input: `git po<Space>` or `git po<CR>`
    # Expansion: `git push origin main`  <-- if current branch is 'main

  # --------------------------------------------------------------------
  # IV. Conditional and environmentally-aware abbreviations
  # --------------------------------------------------------------------
  - name: Use trash-cli if available
    abbr: del
    snippet: trash
    options:
      # This abbreviation is valid only if `condition` is true (the trash command is present).
      # The string is executed in a shell and the condition is determined by its success or failure (exit status).
      condition: 'command -v trash &> /dev/null'
    # Condition: The `trash` command must be installed on the system.
    # Input: `del some_file.txt<Space>` or `del some_file.txt<CR>`
    # Expansion: `trash some_file.txt`

  - name: Fallback to interactive rm if trash-cli is not installed
    abbr: del
    snippet: 'rm -i'
    # Condition: If the `del` condition above is not satisfied (candidates with the same abbr are used here).
    # Input: `del some_file.txt<Space>` or `del some_file.txt<CR>`
    # Expansion: `rm -i some_file.txt`

  - name: Open file with VSCode if code command exists
    abbr: code
    snippet: 'code .'
    options:
      condition: 'type code &> /dev/null'
    # Condition: The `code` command must be installed on the system.
    # Input: `code<Space>` or `code<CR>`
    # Expansion: code .

  # --------------------------------------------------------------------
  # V. Suffix alias (execution by regular expression pattern)
  # --------------------------------------------------------------------
  - name: Execute Python scripts with filename
    # The `$file` variable in the snippet is replaced by the value captured in options.regex.
    snippet: 'python3 $file'
    options:
      # Matches words ending with `.py` and captures the entire word as `file`.
      regex: '^(?<file>\S+\.py)$'
      # 'evaluate: true' enables the process of assigning the named group (`file`) captured by regex to the variable of the same name (`$file`) in snippet.
      evaluate: true
    # Condition: When a word ending with `.py` is entered at the beginning of a line (command part).
    # Input: `myscript.py<Space>` or `myscript.py<CR>`
    # Expansion: `python3 myscript.py`

  - name: Execute shell scripts
    snippet: 'bash %'
    options:
      # Capture words ending in `.sh` with the name `script`.
      regex: '^(?<script>\S+\.sh)$'
      evaluate: true
    # Condition: When a word ending in `.sh` is entered at the beginning of a line.
    # Input: `./run.sh<Space>` or `./run.sh<CR>`
    # Expansion: `bash ./run.sh`

  - name: Unpack tar.gz files
    snippet: 'tar -xzvf $archive'
    options:
      # Capture words ending in `.tar.gz` with the name `archive`.
      regex: '^(?<archive>\S+\.tar\.gz)$'
      evaluate: true
    # Condition: When a word ending with `.tar.gz` is entered at the beginning of a line.
    # Input: `archive.tar.gz<Space>` or `archive.tar.gz<Enter>`
    # Expansion: `tar -xzvf archive.tar.gz`
```

## License

MIT

## Author

Yoshiki Horinaka
