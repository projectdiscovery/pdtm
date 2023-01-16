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
  <a href="#running-tlsx">Running tlsx</a> •
  <a href="https://discord.gg/projectdiscovery">Join Discord</a>
</p>


A simple project that enables users to list, install, update, and remove the projectdiscovery tools.

## Features:

* List all the available pd tools (default)
* Install given pd tool
* Install all the available pd tools
* Remove given pd tool
* Remove all the pd tool
* Update given pd tool
* Update all the pd tools

## Usage: 


```console
Usage:
  ./pdtm [flags]

Flags:
CONFIG:
   -config string            cli flag configuration file (default "$HOME/.config/pdtm/config.yaml")
   -bp, -binary-path string  custom location to download project binary (default "$HOME/.pdtm/go/bin")
   -nc, -no-color            disable output content coloring (ANSI escape codes)

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
   -version  show version of the project
   -v        show verbose output
```

## Running pdtm

```console

$ pdtm -i httpx,nuclei -u naabu,dnsx
                __ __           
    ____   ____/ // /_ ____ ___ 
   / __ \ / __  // __// __  __ \
  / /_/ // /_/ // /_ / / / / / /
 / .___/ \__,_/ \__//_/ /_/ /_/ 
/_/                      v0.0.1

		projectdiscovery.io

[INF] Installed httpx v1.1.1
[INF] Installed nuclei v2.6.3
[INF] Updated to naabu v2.6.3
[INF] Updated to dnsx v2.6.3

```