# Developing the UI

## Installation

From this directory (hereby known as `$ui`) run:  
`npm install`

go-bindata is also required, install with one of:
`sudo apt-get install go-bindata`
or
`go get -u github.com/jteeuwen/go-bindata/...`

## Start Compilin'

`npm run dev`

## Compiling Topiary w/ Mithril UI

From `$GOPATH/src/github.com/topiary-io/topiary`, run `./dev.sh`

This shell script compiles `bindata.go` and then installs `topiary`

Now start your hugo server with `topiary` and  navigate to `localhost:3000/mithril-ui`

Your index template and javascript files will update after refreshes.
