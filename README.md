# Staticli

A CLI tool for working with static site generators.

A [Cali](https://github.com/skybet/cali) App

This application provides a number of tools related to static site generation, all running in docker containers from a single binary file.  Rather than go through the hassle of installing Ruby, rake, bundler, jekyll, python, proselint, gulp, etc. you only need the one single binary file and to have docker installed.

## Installation

1. Install Docker (`brew cask install docker` on MacOS)
2. Download the [correct binary](https://bintray.com/staticli/staticli/staticli/_latestVersion), move it into your $PATH as `staticli` and make it executable

## Usage

From this one single binary you can now run rake tasks to preview and validate Jekyll blogs, run proselint to check for best practises in writing, and run gulp to turn sass into css.  All of this happens inside docker containers, so you don't actually need to install any extra tooling.  You will need an internet connection the first time you run each subcommand in order to download the container though.

`staticli rake` - runs `bundle install --path=_vendor && bundle exec rake $@` in the current directory.  Assuming this is a Jekyll site with a default rake task of previewing the site, you can now open a browser and view the site at http://127.0.0.1:4000.  If you're already using port 4000 you can add `-p 2000` to change to port 2000 (or any port)

`staticli proselint README.md` runs proselint against the file README.md to check for best practises in writing.

`staticli gulp` runs the gulp watch task.  You can add `-t foo` to run the foo task instead.

`staticli heroku` runs the Heroku cli, and takes any parameters and subcommands you need (uses ~/.netrc for authentication)

[ ![Download](https://api.bintray.com/packages/staticli/staticli/staticli/images/download.svg) ](https://bintray.com/staticli/staticli/staticli/_latestVersion)