nmap ,r :call VimuxRunCommand("go run *.go")<cr>
nmap ,t :call VimuxRunCommand("ruby test.rb")<cr>
nmap ,k :!killall a.out; killall go<cr><cr>

nmap ,ct :!gotags *.go > tags<cr>
