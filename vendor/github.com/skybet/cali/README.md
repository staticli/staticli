# Cali - Cali Automation Layout Initialiser

![Cali logo](docs/cali.png)

Oh, you found this repo, huh?

Well, fair warning, we're still figuring out how the open source version of Cali should work.

Feel free to play with it, but while we're still on v0.x.x, be sure to pin to specific versions in your dependency manager of choice.

---

## Overview

**Ever wanted to be able to ensure that all your developers are working using the same versions of the same tools?**

**Want to be free from dowloading huge Vagrant boxes?**

**Want to be able to update and distribute new versions of your tooling quickly and easily?**

The answer to all these questions at our company was a resounding *YES!* 

We needed a way to reduce the overheads associated with onboarding new starters, so instead of a document giving a list of software that needed to be downloaded, e.g. ChefDK, Vagrant, Terraform _etc. etc._ they could just download a single _thing_ which gave them a ready made development environment whether they used Mac or Linux PCs. Vagrant was the obvious choice and for a while it worked, but as we scaled up as a tech team both in numbers and in diversity of tools, this quickly became unmanageable. Different teams started wanting different boxes with different thing and different versions of those things. We soon ended up with several huge box files to maintain, they were slow to update, slow to upload and slow to download and horrible to maintain. We spent a lot of time troubleshooting...

> "Oh yes I have seen that error before, are you using Vagrant version X with VirtualBox version Y?

> "Err... not sure this will work with Windows sorry!"

Then along came Docker for Mac and shortly after, Docker for Windows and this opened the world of containers up to non-Linux users and gave us the ability to distribute our tools using Docker across Mac, Linux and Windows for the first time. Being a Chef house, the first experiment was to run ChefDK out of a docker container and was distributed as some bash aliases (early days!).

```
alias buildtools="docker run --rm -it -v \$PWD:/root/build -v /var/run/docker.sock:/var/run/docker.sock -v $HOME/.aws:/root/.aws -e \"BUILD_ROOT=\$PWD\" -v ~/.gitconfig:/root/.gitconfig ourprivatedockerregistry.io/buildtools:stable"
alias kitchen="buildtools kitchen"
```

Typing `kitchen` would then execute Test Kitchen for Chef against your $PWD and no matter what laptop you were using, as long as you had docker, all our devs were now using the same version of ChefDK for developing Chef code against.

Clearly this was suboptimal and as soon as more tools came along such as Terraform which wanted to be able to obtain short leased AWS credentials from HashiCorp Vault, a new solution was needed to distribute and orchestrate these containers was needed. Also, would it not be cool if that same tool could be run directly, or as a container on CI? Then your developers and your CI is using exactly the same stuff?

So we broke out those Go ninja skills and started using the Docker API to programatically manage these containers, remembering at each step that at some point in the future, this tool would itself be containerised and and would need to be able to schedule its ephemeral job containers either directly on the host its running on, or tantalisingly, on a swarm. Each tool, whether it be ChefDK or Terraform or whatever would either need to be able to work with $PWD or be able to check out some git repo and do its thing within that clone, all within containers. The final pieces of the puzzle were including auth with Hashicorp Vault to allow us to connect our corporate AD to get on-demand AWS credentials and also an update command which would self update the code when a new version was available on our internal Artifact repository.

This has been a great success and we are really happy with how it works and how little day to day support we have to provide our developers to use it. Its been such a success that we are getting requests to put functionality into it which we feel does not really belong there. Rather than saying *NO!* to such feature requests, we are doing the only sane thing and open sourcing the guts of the system to allow anybody to go and create their own development tools to do whatever they want.

## Building a CLI tool

```
package main

import "github.com/skybet/cali"

func main() {
	cli := cali.NewCli("cali")
	cli.SetShort("Example CLI tool")
	cli.SetLong("A nice long description of what your tool actually does")

	cmdTerraform(cli)

	cli.Start()
}

func cmdTerraform(cli *cali.Cli) {

	terraform := cli.NewCommand("terraform [command]")
	terraform.SetShort("Run Terraform in an ephemeral container")
	terraform.SetLong(`Starts a container for Terraform and attempts to run it against your code. There are two choices for code source; a local mount, or directly from a git repo.

Examples:

  To build the contents of the current working directory using my_account as the AWS profile from the shared credentials file on this host.
  # cali terraform plan -p my_account

  Any addtional flags sent to the terraform command come after the --, e.g.
  # cali terraform plan -- -state=environments/test/terraform.tfstate -var-file=environments/test/terraform.tfvars
  # cali terraform -- plan -out plan.out
`)
	terraform.Flags().StringP("profile", "p", "default", "Profile to use from the AWS shared credentials file")
	terraform.BindFlags()

	terraformTask := terraform.Task("hashicorp/terraform:0.9.9")
	terraformTask.SetInitFunc(func(t *cali.Task, args []string) {
		t.AddEnv("AWS_PROFILE", cli.FlagValues().GetString("profile"))
	})
}
```

Now you can run containerised terraform :D

```
$ example terraform init
Terraform initialized in an empty directory!

The directory has no Terraform configuration files. You may begin working
with Terraform immediately by creating Terraform configuration files.
```

Or build from a git repo...

```
$ example terraform plan --git git@github.com:someone/terraform_code.git --git-branch master --git-path path/to/code
```

## Docs

The API [docs](https://godoc.org/github.com/skybet/cali) are available on GoDoc or take a look (and help edit) [the wiki](https://github.com/skybet/cali/wiki).

