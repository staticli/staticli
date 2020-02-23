# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [1.8.8] - 2020-02-23
### Fixed
- Do not need to bind a port if we're only running Jekyll Build

## [1.8.2] - 2018-04-10
### Fixed
- Incorrect name in version and update commands

## [1.8.1] - 2018-04-09
### Changed
- Rebuilt against a new version of Cali, no functional changes
### Added
- `version --check-update=false` because Travis is failing to connect to the GitHub releases API

## [1.8.0] - 2018-04-07
### Added
- Update can now update local Docker images if `-i` is passed in

## [1.7.1] - 2018-03-19
### Changed
- Update now prints status bar while updating

## [1.7.0] - 2018-03-09
### Added
- mkdocs

## [1.6.0] - 2018-03-07
### Added
- ponysay

## [1.5.2] - 2018-03-04
### Changed
- Correct environment variable is set for Jekyll UID & GID
- Jekyll now uses cache directory on host (needs internet access for first run)

## [1.5.1] - 2018-02-28
### Changed
- Added option to run specific tagged version of Rake-based commands, allowing you to use `--tag ruby2.4`

## [1.5.0] - 2018-02-22
### Added
- Silver Searcher (`ag`)

## [1.4.1] - 2018-02-22
### Changed
- Jekyll now exposes a port for `jekyll serve`

## [1.4.0] - 2018-02-21
### Added
- jekyll

## [1.3.0] - 2018-02-18
### Added
- npm and npx commands from node

## [1.2.0] - 2018-02-11
### Added
- update command to update to the latest version of staticli

## [1.1.1] - 2018-02-10
### Changed
- version cmd now checks for latest version from github api

## [1.1.0] - 2018-02-08
### Added
- github-release for releasing new versions of staticli

## [1.0.0] - 2018-02-03
### Changed
- Default port is 2000 for everything that listens
- Log port being used
### Added
- Bundle command which reuses the rake image
- Hugo command

## [0.7.0] - 2018-02-01
### Added
- Heroku

## [0.6.0] - 2018-01-28
### Added
- Python SimpleHTTPServer

## [0.5.0] - 2018-01-26
### Added
- surge.sh

## [0.4.0] - 2018-01-25
### Added
- gulp

## [0.3.0] - 2018-01-25
### Changed
- Switch to staticli/rake docker image which refactors the rake docker image to provide arbitrary rake commands

## [0.2.0] - 2018-01-25
### Added
- proselint

## [0.1.0] - 2018-01-25
### Changed
- Run rake preview on port 4000 by default

## [0.0.1] - 2018-01-24
### Added
- rake
