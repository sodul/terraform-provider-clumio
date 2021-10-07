# Terraform Provider Clumio

The Terraform Clumio provider is a plugin for Terraform that allows for the full 
lifecycle management of Clumio resources. 

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.14.x
-	[Go](https://golang.org/doc/install) >= 1.17

## Installing the provider

To get started quickly and conveniently, install using the following command:
```sh
$ curl https://raw.githubusercontent.com/clumio-code/terraform-provider-clumio/main/installer.sh | bash -s -- VERSION
```
Replace VERSION with the appropriate release version.
### Alternative: Clone the repository and run installer.sh
```sh
$ git clone https://github.com/clumio-code/terraform-provider-clumio.git
$ ./installer.sh VERSION
```
Replace VERSION with the appropriate release version.
