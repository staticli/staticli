# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

[Unreleased]: https://github.com/skybet/cali/compare/v0.1.2...master
## [Unreleased]
### Changed
- Slight refactoring based on static analysis
- DockerClient.InitDocker does nothing if Docker client was previously initialised
- All DockerClient methods which rely on an initialised Docker client now run InitDocker. As a result, a few DockerClient methods now have the possibilities of returning errors: ContainerExists, ImageExists

[0.1.2]:      https://github.com/skybet/cali/compare/v0.1.1...v0.1.2
## [0.1.2] - 2018-04-07
### Fixed
- PullImage now always pulls as should be expected
- StartContainer now only calls PullImage when the image does not exist locally
- Other miscelaneous refactoring

[0.1.1]:      https://github.com/skybet/cali/compare/v0.1.0...v0.1.1
## [0.1.1] - 2018-02-11
### Added
- This CHANGELOG file

### Changed
- Git data containers now have the repo name, branch and directory in the container name


[0.1.0]:      https://github.com/skybet/cali/compare/init...v0.1.0
## [0.1.0] - 2018-02-03
### Added
- Tagged a stable version of Cali which devs can pin to
