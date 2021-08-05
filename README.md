# Terraform Provider Clumio

The Terraform Clumio provider is a plugin for Terraform that allows for the full 
lifecycle management of Clumio resources. 

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.16

## Building The Provider

1. Clone the repository.
2. Enter the repository directory.
3. Open the GNUmakefile and update the OS_ARCH if required. The default value is set to darwin_amd64. Information about the recommended OS and Architectures can be found [here](https://www.terraform.io/docs/registry/providers/os-arch.html).
4. Build the provider using the `make install` command: 
```sh
$ make install
```
This will build the provider and put the provider binary in the Terraform plugins 
directory.
