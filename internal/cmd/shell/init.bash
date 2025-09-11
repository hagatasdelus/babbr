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

__babbr_handle_space() {
    local line_before_cursor="${READLINE_LINE:0:READLINE_POINT}"
    local word_before_cursor="${line_before_cursor##*[( ]}"
    
    # Skip expansion for obviously non-expandable patterns
    if 
        [[ -z "$word_before_cursor" ]] || \
        [[ ${#word_before_cursor} -gt 30 ]] || \
        [[ "$word_before_cursor" =~ ^(\$.*)$ ]];
    then
        READLINE_LINE="${READLINE_LINE:0:READLINE_POINT} ${READLINE_LINE:READLINE_POINT}"
        ((READLINE_POINT++))
    else
        expand_abbr " "
    fi
}

bind -x '" ": __babbr_handle_space'
bind -x '";": expand_abbr ";"'
bind -x '"|": expand_abbr "|"'
bind -m vi-insert -x '" ": __babbr_handle_space'
bind -m vi-insert -x '";": expand_abbr ";"'
bind -m vi-insert -x '"|": expand_abbr "|"'

key_seq_expand_abbr_enter='\C-x\C-['
key_seq_accept_line='\C-j'

bind -m emacs "\"\C-m\": \"$key_seq_expand_abbr_enter$key_seq_accept_line\""
bind -m emacs -x "\"$key_seq_expand_abbr_enter\": expand_abbr"

bind -m vi-insert "\"\C-m\": \"$key_seq_expand_abbr_enter$key_seq_accept_line\""
bind -m vi-insert -x "\"$key_seq_expand_abbr_enter\": expand_abbr"
