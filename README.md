# Staticli

A CLI tool for working with static site generators.

A [Cali](https://github.com/skybet/cali) App

This application will provide a number of static site generator applications over the course of time.  For now all it does is provide a way of running rake preview against a Jekyll site without requiring any dependencies.

## Usage

`staticli rake` - runs `bundle && bundle exec rake preview` for the Jekyll site in the current directory, within a docker image.

It serves the port on a random number because we haven't figured how to make that static yet.  You'll need to run `docker ps` in order to find out what port is being used.


[ ![Download](https://api.bintray.com/packages/wheresalice/staticli/staticli/images/download.svg) ](https://bintray.com/wheresalice/staticli/staticli/_latestVersion)
[![Build Status](https://travis-ci.org/WheresAlice/staticli.svg?branch=master)](https://travis-ci.org/WheresAlice/staticli)