# babbr

Abbreviation extension plugin for bash

## Installation

```
go install github.com/hagatasdelus/babbr
```

## Usage
```
Usage of babbr
...
```

### Configuration

Configuration can be done with a `config.yaml` file.

### Example for config.yaml

Location of config.yaml is:

* UNIX: `$XDG_CONFIG_HOME/babbr/config.yaml` or `$HOME/.config/babbr/config.yaml`

```yaml
abbreviations:
  - name: git checkout
    abbr: gco
    options:
      position: command # "command"(Default) or "anywhere"
    - name: '% | less'
      abbr: L
      options:
        position: anywhere
        set_cursor: true # If true, ‘%’ in expansion becomes the cursor position
    - name: vim %
      abbr: vtxt
      options:
        position: command
        regex: '.+\.txt' # PCRE2 compatible regular expressions
    - name: checkout
      abbr: co
      options:
        # Expanded only when used as an argument to the ‘git’ command
        command: git
        position: anywhere
```

## License

MIT

## author

Yoshiki Horinaka
