# Staticli

A CLI tool for working with static site generators.

A [Cali](https://github.com/skybet/cali) App

This application will provide a number of static site generator applications over the course of time.

## Usage

`staticli rake` - runs `bundle install --path=_vendor && bundle exec rake $@` in the current directory, within a docker image.  Assuming this is a Jekyll site with a default rake task of previewing the site, you can now open a browser and view the site at http://127.0.0.1:4000

`staticli rake -p 2000` runs it on port 2000 instead of the default port of 4000

`staticli proselint README.md` runs proselint against the file README.md to check for best practises in writing.

[ ![Download](https://api.bintray.com/packages/wheresalice/staticli/staticli/images/download.svg) ](https://bintray.com/wheresalice/staticli/staticli/_latestVersion)