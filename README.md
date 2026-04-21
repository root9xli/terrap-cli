# Terrap, by Sirrend
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  ![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/sirrend/terrap-cli?filename=go.mod)</br>
The tool offers clear and</br>
🔍 **Terrap is an alpha> **Personal note:** I'm using this tool to manage provider upgrades across several AWS-heavy workspaces. Works well paired with.

## Resources
* Documentation - <a href="https://www.sirrend.com/terrap-docs">sirrend.com/terrap-docs</a>

## Constraints 🧱
1. Supported Terraform Core versions: `>=0.13`.
2. Every provider which uses `Terraform Core 0.13` or higher.

## Good To Know 💡
Terrap decides which Terraform version to use in the following order:
1. The latest installed Terraform version found locally.
2. If the `TERRAP_TERRAFORM_VERSION` environment variable is set, Terrap will use the version specified as long as it matches the `>=0.13` constraint.</br>
    Set environment variable on mac/linux:</br>
    ```shell
   export <var>=<value>
   ```
   Set environment variable on windows:</br>
    ```shell
   $Env:TERRAP_TERRAFORM_VERSION = "0.13"
   ```
   
4. If none of the above is applicable, Terrap will download the latest available version.

## How to Download ⬇️
### Clone sirrend/terrap-cli
```shell
git clone https://github.com/sirrend/terrap-cli
cd terrap-cli

go build -o terrap .

chmod +x terrap
mv terrap /usr/local/bin/
```

### Brew
```shell
brew tap sirrend/products
brew install terrap
```

Validate terrap is working by executing `terrap`.

> **Tip (personal):** Add `alias tp="terrap"` to your `~/.zshrc` or `~/.bashrc` to save keystrokes.

## Quick Start ⏩

### Initialize my First Workspace
1. `CD` to the local Terraform folder you want to work with.</br>
   `cd < /terraform/folder/path >`</br></br>

2. Initialize a new Terrap workspace where you would run terraform apply with `terrap init`.</br></br>
    <strong>Important!</strong> </br>
    As Terrap runs <code>terraform init</code> under the hood, it would need every configuration component you normally use when executing the command.</br>
    It can be environment variables, the <code>.aws/credentials</code> file, etc.


3. Scan your workspace with: `terrap scan`

https://user-images.githubusercontent.com/47568615/236044850-3473952a-4169-4d63-beb7-cf1664afc35a.mov

## Features 🚀
### Scan for changes with `scan`
Scan your infrastructure for changes in the following provider version for a safe and easy upgrade!</br>
Looking for a specific resource type changes? Use the `--data-sources` `--resources` and `--provider` flags.

> **Personal note:** I primarily use `terrap scan --resources` scoped to the `aws` provider. Filtering by resource type significantly cuts down noise when working in large workspaces with multiple providers.

### Stay up-to-date with `whats-new`
Ready to explore what's new in the following version of your provider? Simply execute `terrap whats-new`.</br>
Looking to delve 
