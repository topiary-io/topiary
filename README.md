# Topiary

A content management system built on [Hugo](github.com/spf13/hugo), a static website generator.

### Status

Topiary is currently alpha software. It will build your Hugo source, start a web server, and make the input viewable at `/admin`. Topiary should not be used in anything remotely resembling a production environment.

### Installation

Right now the only way to install is by building from source.

0. Make sure to have [Git](https://git-scm.com/downloads), [Go](https://golang.org/dl/), and [Hugo](http://gohugo.io/overview/installing/) installed.

1. Clone the Topiary repo:
```
go get github.com/topiary-io/topiary
```

2. Build Topiary:
```
cd topiary
go install
```

2.5. If you don't have a Hugo site, clone [a bootstrapped site](https://github.com/enten/hyde-y):
```
cd $HOME/project-dir
git clone https://github.com/enten/hyde-y
```

3. Start Topiary:
```
cd path/to/site/root
topiary
```

4. Point a web browser to localhost:3000 to view your site, and localhost:3000/admin to view the site input.

### License
Apache 2.0

