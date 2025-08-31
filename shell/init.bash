#!/bin/bash

__babbr_expand() {
    local output
    output="$(babbr expand --lbuffer="$READLINE_LINE" --rbuffer="")"
    if [[ $? -eq 0 && -n "$output" ]]; then
        eval "$output"
    fi
}

__babbr_expand_and_execute() {
    __babbr_expand
    builtin echo
}

__babbr_expand_and_space() {
    __babbr_expand
    builtin echo -n " "
}

bind -x '"\C-M": __babbr_expand_and_execute'
bind -x '" ": __babbr_expand_and_space'
bind '"\C-x ": " "'
bind '"\C-x\C-m": accept-line'
