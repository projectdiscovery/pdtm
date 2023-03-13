<h1 align="center">
<img src="https://user-images.githubusercontent.com/8293321/211602034-411e38e9-e5df-429e-89ee-a97e3e09ebf0.png" width="200px">
<br>
</h1>

<h4 align="center">ProjectDiscovery's Open Source Tool Manager</h4>

<p align="center">
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-_red.svg"></a>
<a href="https://goreportcard.com/badge/github.com/projectdiscovery/pdtm"><img src="https://goreportcard.com/badge/github.com/projectdiscovery/pdtm"></a>
<a href="https://github.com/projectdiscovery/pdtm/releases"><img src="https://img.shields.io/github/release/projectdiscovery/pdtm"></a>
<a href="https://twitter.com/pdiscoveryio"><img src="https://img.shields.io/twitter/follow/pdiscoveryio.svg?logo=twitter"></a>
<a href="https://discord.gg/projectdiscovery"><img src="https://img.shields.io/discord/695645237418131507.svg?logo=discord"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#running-pdtm">Running pdtm</a> •
  <a href="https://discord.gg/projectdiscovery">Join Discord</a>

**pdtm** is a simple and easy-to-use golang based tool for managing open source projects from ProjectDiscovery.

</p>

<h1 align="center">
<img src="https://user-images.githubusercontent.com/8293321/212781914-bae85495-5a7b-40d7-9e05-964a8edf3b61.png" width="700px">
</h1>

## Installation


**`pdtm`** requires **go1.19** to install successfully. Run the following command to install the latest version:

1. Install using go install -

```sh
go install -v github.com/projectdiscovery/pdtm/cmd/pdtm@latest
```

2. Install by downloading binary from https://github.com/projectdiscovery/pdtm/releases

<table>
<tr>
<td>  

> **Notes**:

> - *Currently, projects are installed by downloading the released project binary. This means that projects can only be installed on the platforms for which binaries have been published.*
> - *The path $HOME/.pdtm/go/bin is added to the $PATH variable by default*

</table>
</tr>
</td> 

## Usage: 


```console
Usage:
  ./pdtm [flags]

Flags:
CONFIG:
   -config string            cli flag configuration file (default "$HOME/.config/pdtm/config.yaml")
   -bp, -binary-path string  custom location to download project binary (default "$HOME/.pdtm/go/bin")
   -nsp, -no-set-path        disable adding path to environment variables

INSTALL:
   -i, -install string[]  install single or multiple project by name (comma separated)
   -ia, -install-all      install all the projects

UPDATE:
   -u, -update string[]  update single or multiple project by name (comma separated)
   -ua, -update-all      update all the projects

REMOVE:
   -r, -remove string[]  remove single or multiple project by name (comma separated)
   -ra, -remove-all      remove all the projects

DEBUG:
   -sp, -show-path  show the current binary path then exit
   -version         show version of the project
   -v, -verbose     show verbose output
   -nc, -no-color   disable output content coloring (ANSI escape codes)
```

## Running pdtm

```console
$ pdtm -install-all
                ____          
     ____  ____/ / /_____ ___ 
    / __ \/ __  / __/ __ __  \
   / /_/ / /_/ / /_/ / / / / /
  / .___/\__,_/\__/_/ /_/ /_/ 
 /_/                          v0.0.1

      projectdiscovery.io

[INF] Installed httpx v1.1.1
[INF] Installed nuclei v2.6.3
[INF] Installed naabu v2.6.3
[INF] Installed dnsx v2.6.3
``` 

### Todo

- support for go setup + project install from source
- support for installing from source as fallback option

--------

<div align="center">

**pdtm** is made with ❤️ by the [projectdiscovery](https://projectdiscovery.io) team and distributed under [MIT License](LICENSE).


<a href="https://discord.gg/projectdiscovery"><img src="https://raw.githubusercontent.com/projectdiscovery/nuclei-burp-plugin/main/static/join-discord.png" width="300" alt="Join Discord"></a>

</div>