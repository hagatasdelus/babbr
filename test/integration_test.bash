#!/usr/bin/env bash

failed=0

export USERNAME="BABBR"
export ABBREV_MODE="DEBUG"
export LONG_PATH="$HOME/some/really/long/path/to/directory"

expand_test() {
    local lbuffer="$1"
    local rbuffer="$2"
    local expected_lbuffer="$3"
    local expected_rbuffer="$4"
    
    echo -e "\$ ${lbuffer}\033[1;31m|\033[0;39m${rbuffer}"

    local out exit_code
    out="$(./babbr expand --lbuffer="$lbuffer" --rbuffer="$rbuffer" 2>&1)"
    exit_code=$?
    if [ "$exit_code" -ne 0 ]; then
        echo "  babbr expand failed with status $exit_code" >&2
        ((failed++))
        return
    fi

    local LBUFFER RBUFFER READLINE_LINE READLINE_POINT SET_CURSOR
    
    eval "$out"
    
    LBUFFER="${READLINE_LINE:0:$READLINE_POINT}"
    RBUFFER="${READLINE_LINE:$READLINE_POINT}"

    if [ "$LBUFFER" != "$expected_lbuffer" ]; then
        echo -n "  -- FAILED: LBUFFER not matched"
        ((failed++))
    else
        echo -n "  -- PASSED: LBUFFER matched"
    fi
    echo "  (expected: '$expected_lbuffer', actual: '$LBUFFER')"

    if [ "$RBUFFER" != "$expected_rbuffer" ]; then
        echo -n "  -- FAILED: RBUFFER not matched"
        ((failed++))
    else
        echo -n "  -- PASSED: RBUFFER matched"
    fi
    echo "  (expected: '$expected_rbuffer', actual: '$RBUFFER')"
}

expand_test "g"                     ""          "git"                                   ""
expand_test "  g"                   ""          "  git"                                 ""
expand_test "g"                     "add"       "git"                                   "add"
expand_test "g"                     " add"      "git"                                   " add"
expand_test "echo g"                ""          "echo g"                                ""
expand_test "echo test; g"          ""          "echo test; git"                        ""
expand_test "echo TEST && git s"    ""          "echo TEST && git status"               ""
expand_test "cat test.txt L"        ""          "cat test.txt | less"                   ""
expand_test "welcome"               ""          "echo 'Hello, BABBR! Welcome back.'"    ""
expand_test "echo calc"             ""          "echo 3 * 4 * 5 = 60"                   ""
expand_test "calc"                  ""          "3 * 4 * 5 = 60"                        ""
expand_test "git s"                 ""          "git status"                            ""
expand_test "hg s"                  ""          "hg s"                                  ""
expand_test "echo docker RMI"       ""          "echo docker RMI"                       ""
expand_test "git cm"                ""          "git commit -m '"                       "'"
expand_test "git cm"                " -v"       "git commit -m '"                       "' -v"
expand_test "docker compose bnc"    ""          "docker compose build --no-cache"       ""
expand_test "vim /etc/shells"       ""          "sudo vim /etc/shells"                  ""
expand_test "../"                   ""          "cd ../"                                ""
expand_test "../../"                ""          "cd ../../"                             ""
expand_test "../../../"             ""          "cd ../../../"                          ""
expand_test "test.ts"               ""          "deno run -A test.ts"                   ""
expand_test "echo Test | ./test.ts" ""          "echo Test | deno run -A ./test.ts"     ""
expand_test "deno run test.ts"      ""          "deno run test.ts"                      ""
expand_test "test.tsx"              ""          "test.tsx"                              ""
expand_test "log"                   ""          "echo '[DEBUG] Log message here'"       ""
expand_test "log-user"              ""          "log-user"                              ""
expand_test "del"                   ""          "rm -i"    ""
expand_test "ls -l | p2"            ""          "ls -l | awk '{ print \$2 }'"           ""
expand_test "long"                  ""          "cd $HOME/some/really/long/path/to/directory"  ""

if [ "$failed" -ne 0 ]; then
    echo "❎ ${failed} Test Failed!" >&2
    exit 1
fi

echo "✅ All Test Passed!"
