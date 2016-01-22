# Topiary

A content management system built on [Hugo](https://github.com/spf13/hugo), a static website generator.

### Status

Topiary is currently alpha software. It will build your Hugo source, start a web server, and make the input viewable at `/admin`. Topiary should not be used in anything remotely resembling a production environment.

### Use cases

1. Run Topiary on your computer, for individual use. Manually copy or automatically sync the output to your production environment, either as changes are made or at regular intervals.
2. Run Topiary on a local server or intranet, for team use. Manually copy or automatically sync the output to your production environment, either as changes are made or at regular intervals.
3. Run Topiary in your production environment, for individual or team use. Have your site rebuild as you make changes.

Possible production environments are GitHub Pages, a VPS, a shared host, [IPFS](ipfs.io), or your own server. Since the output is static, it's lightweight, secure, cheap to host, and simple to set up.

### Installation

Right now the only way to install is by building from source.

1. Make sure to have [Git](https://git-scm.com/downloads), [Go](https://golang.org/dl/), and [Hugo](http://gohugo.io/overview/installing/) installed.

2. Clone the Topiary repo:
  ```
  $ go get github.com/topiary-io/topiary
  ```

3. Build Topiary:
  ```
  $ cd topiary
  $ go install
  ```

4. If you don't have a Hugo site, clone [a bootstrapped site](https://github.com/enten/hugo-boilerplate):
  ```
  $ cd path/to/project/dir
  $ git clone https://github.com/enten/hugo-boilerplate
  ```

5. Start Topiary:
  ```
  $ cd path/to/hugo/site/root
  $ topiary
  ```

6. Point a web browser to `localhost:3000` to view your site, and `localhost:3000/admin` to view the site input.

### License

Apache 2.0
