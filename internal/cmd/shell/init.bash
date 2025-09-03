#!/bin/bash

expand_abbr() {
    local key="$1"
    
    local line_before_cursor="${READLINE_LINE:0:READLINE_POINT}"
    local suffix="${READLINE_LINE:READLINE_POINT}"
    
    if [[ -z "$line_before_cursor" ]]; then
        if [[ -n "$key" ]]; then
            READLINE_LINE="${READLINE_LINE:0:READLINE_POINT}${key}${READLINE_LINE:READLINE_POINT}"
            ((READLINE_POINT += ${#key}))
        fi
        return 0
    fi

    local output
    output="$(babbr expand --lbuffer="$line_before_cursor" --rbuffer="$suffix" 2>/dev/null)"
    local exit_code=$?

    if [[ $exit_code -eq 0 && -n "$output" ]]; then
        eval "$output"
        # Add trigger character only if cursor wasn't manually positioned
        if [[ "$output" != *"SET_CURSOR=1"* && -n "$key" ]]; then
            READLINE_LINE="${READLINE_LINE:0:READLINE_POINT}${key}${READLINE_LINE:READLINE_POINT}"
            ((READLINE_POINT += ${#key}))
        fi
    else
        if [[ -n "$key" ]]; then
            READLINE_LINE="${READLINE_LINE:0:READLINE_POINT}${key}${READLINE_LINE:READLINE_POINT}"
            ((READLINE_POINT += ${#key}))
        fi
    fi
}

# Optimized space handler with minimal pre-screening
__babbr_handle_space() {
    local line_before_cursor="${READLINE_LINE:0:READLINE_POINT}"
    local word_before_cursor="${line_before_cursor##*[( ]}"
    
    # Skip expansion for obviously non-expandable patterns
    if 
        # Early return for empty word
        [[ -z "$word_before_cursor" ]] || \
        # Skip overly long words
        [[ ${#word_before_cursor} -gt 30 ]] || \
        # Skip obvious path patterns
        [[ "$word_before_cursor" =~ ^\.\.?/.*$ ]] || \
        [[ "$word_before_cursor" =~ ^(/|~)[^[:space:]]+/.*$ ]] || \
        # Skip pure numbers and version patterns
        [[ "$word_before_cursor" =~ ^[0-9]+(\.[0-9]+){0,2}$ ]] || \
        # Skip common command-line patterns
        [[ "$word_before_cursor" =~ ^((--.*)|(-[a-zA-Z]+)|(\$.*))$ ]];
    then
        READLINE_LINE="${READLINE_LINE:0:READLINE_POINT} ${READLINE_LINE:READLINE_POINT}"
        ((READLINE_POINT++))
    else
        expand_abbr " "
    fi
}

# Setup key bindings for space, semicolon, and pipe
bind -x '" ": __babbr_handle_space'
bind -x '";": expand_abbr ";"'
bind -x '"|": expand_abbr "|"'
bind -m vi-insert -x '" ": __babbr_handle_space'
bind -m vi-insert -x '";": expand_abbr ";"'
bind -m vi-insert -x '"|": expand_abbr "|"'

# Setup enter key bindings
key_seq_expand_abbr_enter='\C-x\C-['
key_seq_accept_line='\C-j'

bind -m emacs "\"\C-m\": \"$key_seq_expand_abbr_enter$key_seq_accept_line\""
bind -m emacs -x "\"$key_seq_expand_abbr_enter\": expand_abbr"

bind -m vi-insert "\"\C-m\": \"$key_seq_expand_abbr_enter$key_seq_accept_line\""
bind -m vi-insert -x "\"$key_seq_expand_abbr_enter\": expand_abbr"
