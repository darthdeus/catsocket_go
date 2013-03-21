nmap ,r :call VimuxRunCommand("go run *.go")<cr>
nmap ,t :call VimuxRunCommand("go test")<cr>
nmap ,k :!killall a.out; killall go<cr><cr>

nmap ,ct :!gotags *.go > tags<cr>
