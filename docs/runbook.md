# Staticli Runbook

## General
* The Staticli website runs on GitHub Pages with Cloudflare in front of it.
* The Staticli build infrastructure uses Travis CI, GitHub, and GitHub releases.
* The staticli.tech is registered through NameCheap.

## Dashboards
* [Cloudflare Analytics](https://www.cloudflare.com/a/analytics/staticli.tech) (requires a login)
* [GitHub Insights](https://github.com/staticli/staticli/pulse) provides some insights into the project

The GitHub API also provides download counts for releases, which is not currently tracked in any dashboard.

## Alerts
There is currently no alerting for these services
 
## Contact Info
Most contacts should go via GitHub issues, but Alice can also be contacted via Twitter (@WheresAlice).  Any escalations will go via Alice.

## Latest Deployments
All releases are tagged [GitHub Releases](https://github.com/staticli/staticli/releases).  There has been a history of releases not being correctly tagged, and so you can get detailed version information by running `staticli version -d` and compare the commit hash to the git history.

The website is automatically published via GitHub Pages using the latest README.md file in the root of this repository.

## Deployment
Deployment is via GitHub and Travis CI.  Travis CI will create a new GitHub release when the VERSION file is incremented.  The website gets automatically rebuilt from the latest README.md in the master branch.

 1. Push a branch to your own fork of the repository, including updates to VERSION, HISTORY.md, and README.md
 2. Create a Pull Request on GitHub from your own fork of the repository
 3. Wait for Travis CI to mark the PR as successfully building
 4. Approve and merge the PR and wait for Travis CI to perform a release
  
### Rollback Deploy
We prefer to roll forwards rather than backwards, and therefore create a new release which fixes or reverts the issue.  We will generally only remove releases in the case of security vulnerabilties.