# Staticli

A CLI tool for working with static site generators.

A [Cali](https://github.com/skybet/cali) App

This application will provide a number of static site generator applications over the course of time.

## Usage

`staticli rake` - runs `bundle && bundle exec rake preview` for the Jekyll site in the current directory, within a docker image.  You can now open a browser to http://127.0.0.1:4000 to view the site.

`staticli rake -p 2000` runs it on port 2000 instead of the default port of 4000

`staticli proselint README.md` runs proselint against the file README.md to check for best practises in writing.

[ ![Download](https://api.bintray.com/packages/wheresalice/staticli/staticli/images/download.svg) ](https://bintray.com/wheresalice/staticli/staticli/_latestVersion)
[![Build Status](https://travis-ci.org/WheresAlice/staticli.svg?branch=master)](https://travis-ci.org/WheresAlice/staticli)