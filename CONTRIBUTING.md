# Contributing

We're open to pull requests for additional tooling, but it must be relevant to static web site development.

## Dependencies

We use `dep` to manage package dependencies.  You should not manually change anything in the vendor directly, but use `make dep` to ensure dependencies are correct.  Running dep is a slow operation, so we do not do this on build.

## Pull Request Process

1. Increase the version number in VERSION, following [Semantic Versioning](http://semver.org/spec/v2.0.0.html)
2. Update the README.md and HISTORY.md files with details of changes
3. The Pull Request can be merged once Travis CI has successfully done a build
4. Travis CI will automatically release a new version, and the website will automatically be updated via GitHub Pages.