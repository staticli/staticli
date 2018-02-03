# Staticli

A CLI tool for working with static site generators.  A [Cali](https://github.com/skybet/cali) App

This application provides a number of tools related to static site generation, all running in docker containers from a single binary file.  Rather than go through the hassle of installing Ruby, rake, bundler, jekyll, python, proselint, gulp, etc. you only need the one single binary file and to have docker installed.

## Tools

Staticli provides the tools that its developers regularly use to work with static websites.  We're open to pull requests for additional tooling, but it must be relevant to static web site development (feel free to create your own Cali app for other tools).  The tools we currently provide include:

* bundle and rake for working with Jekyll and github-pages websites
* gulp for compiling scss to css
* heroku for deploying sites to heroku
* hugo for providing the hugo static site generator
* proselint for checking files for best practises in writing
* python simplehttp for serving the current directory over http
* surge for deploying sites to surge.sh

## Installation

1. Install Docker (`brew cask install docker` on MacOS)
2. Download the [correct binary](https://bintray.com/staticli/staticli/staticli/_latestVersion), move it into your $PATH as `staticli` and make it executable

[![asciicast](https://asciinema.org/a/159883.png)](https://asciinema.org/a/159883)

## Usage

From this one single binary you can now run rake tasks to preview and validate Jekyll blogs, run proselint to check for best practises in writing, and run gulp to turn sass into css.  All of this happens inside docker containers, so you don't actually need to install any extra tooling.  You will need an internet connection the first time you run each subcommand in order to download the container though.

For any command which exposes a port (typically an http server to render a site) we default to exposing this on port 2000.  You can override this by setting `--port 4000` or `-p 4000` to listen on (for example) port 4000 instead.

`staticli rake` runs `bundle install --path=_vendor && bundle exec rake $@` in the current directory.  This assumes the default rake task runs preview on a Jekyll site, and therefore exposes port 4000 on the container as port 2000 on the host.  This means you can browse to http://127.0.0.1:2000 to view the site.  You can change the port exposed on the host by setting `--port 4000` to use (for example) port 4000.

`staticli bundle` runs `bundle` in the current directory.  Since the `rake` command installs required gems anyway, this command is mostly useful as a way of upgrading gems.

`staticli proselint README.md` runs proselint against the file README.md to check for best practises in writing.

`staticli gulp` runs the gulp watch task.  You can add `-t foo` to run the foo task instead.

`staticli heroku` runs the Heroku cli, and takes any parameters and subcommands you need (uses ~/.netrc for authentication)

`staticli hugo` runs the hugo static site generator

`staticli simplehttp` runs Python SimpleHTTP in the current directory.

`staticli surge` runs surge, allowing you to deploy the current directory to a surge.sh site.

[![asciicast](https://asciinema.org/a/159884.png)](https://asciinema.org/a/159884)

[ ![Download](https://api.bintray.com/packages/staticli/staticli/staticli/images/download.svg) ](https://bintray.com/staticli/staticli/staticli/_latestVersion)