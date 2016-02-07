" Copyright (c) 2016 OKUMURA Takahiro.
"
" gotests.vim: Vim command for integration of the gotests.
"
" Commands:
"
"   :Gotests
"
"       Filter the current Go buffer through gotests.
"
" Options:
"
"   g:go_tests_commands [default=1]
"
"       Flag to enable the gotests commands.

if exists("b:did_ftplugin_go_tests")
    finish
endif

if !exists("g:go_tests_commands")
    let g:go_tests_commands = 1
endif

if g:go_tests_commands
    command! -buffer Gotests call s:GoTests()
endif

function! s:GoTests()
    let view = winsaveview()
    silent execute "%!gotests -r " . expand('%:p')
    if v:shell_error
        let errors = []
        for line in getline(1, line('$'))
            let tokens = matchlist(line, '^\(.\{-}\):\s*\(.\{-}\):\s*\(.\{-}\):\(\d\+\):\(\d\+\)\s*\(.*\)')
            if !empty(tokens)
                call add(errors, {"filename": @%,
                                 \"lnum":     tokens[4],
                                 \"col":      tokens[5],
                                 \"text":     tokens[6]})
            endif
        endfor
        if empty(errors)
            % | " Couldn't detect gotests error format, output errors
        endif
        undo
        if !empty(errors)
            call setqflist(errors, 'r')
        endif
        echohl ErrorMsg | echomsg "Gotests failed to run." | echohl None
    endif
    call winrestview(view)
endfunction

let b:did_ftplugin_go_tests = 1

" vim:ts=4:sw=4:et
