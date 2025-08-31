#!/bin/bash

expand_abbr() {
    local key="$1"

    # Pseudo-tokenize (prefix / word_before_cursor / cursor / suffix)
    local line_before_cursor="${READLINE_LINE:0:READLINE_POINT}"
    local word_before_cursor="${line_before_cursor##*[( ]}"
    local word_before_start_offset=$((READLINE_POINT - ${#word_before_cursor}))
    local prefix="${READLINE_LINE:0:word_before_start_offset}"
    local suffix="${READLINE_LINE:READLINE_POINT}"

    # Determine characters to add
    local add_char="$key"
    local add_char_if_no_expansion="$key"
    if [[ "$key" = "enter" ]]; then
        add_char=""
        add_char_if_no_expansion=""
    fi

    # Call babbr expand
    local output
    output="$(babbr expand --lbuffer="$line_before_cursor" --rbuffer="$suffix" 2>/dev/null)"
    local exit_code=$?

    if [[ $exit_code -eq 0 && -n "$output" ]]; then
        # Expansion occurred - eval the output to set READLINE_LINE and READLINE_POINT
        eval "$output"
        # Check if SET_CURSOR is set (cursor position manually set)
        if [[ "$output" == *"SET_CURSOR=1"* ]]; then
            # Don't add trigger character when cursor is manually positioned
            :
        elif [[ -n "$add_char" ]]; then
            # Add the trigger character after expansion
            READLINE_LINE="${READLINE_LINE:0:READLINE_POINT}${add_char}${READLINE_LINE:READLINE_POINT}"
            ((READLINE_POINT += ${#add_char}))
        fi
    else
        # No expansion - just add the character
        READLINE_LINE="${prefix}${word_before_cursor}${add_char_if_no_expansion}${suffix}"
        READLINE_POINT=$((word_before_start_offset + ${#word_before_cursor} + ${#add_char_if_no_expansion}))
    fi
}

# Setup intercepts for space
bind -x '" ": expand_abbr " "'
bind -x '";": expand_abbr ";"'
bind -x '"|": expand_abbr "|"'
bind -m vi-insert -x '" ": expand_abbr " "'
bind -m vi-insert -x '";": expand_abbr ";"'
bind -m vi-insert -x '"|": expand_abbr "|"'

# Setup intercept for enter
key_seq_expand_abbr_enter='\C-x\C-['
key_seq_accept_line='\C-j'

# Emacs keymap
bind -m emacs "\"\C-m\": \"$key_seq_expand_abbr_enter$key_seq_accept_line\""
bind -m emacs -x "\"$key_seq_expand_abbr_enter\": expand_abbr enter"

# VI insert keymap
bind -m vi-insert "\"\C-m\": \"$key_seq_expand_abbr_enter$key_seq_accept_line\""
bind -m vi-insert -x "\"$key_seq_expand_abbr_enter\": expand_abbr enter"
