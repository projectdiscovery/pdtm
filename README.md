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
```bash
$ pdtm
                __ __           
    ____   ____/ // /_ ____ ___ 
   / __ \ / __  // __// __  __ \
  / /_/ // /_/ // /_ / / / / / /
 / .___/ \__,_/ \__//_/ /_/ /_/ 
/_/                      v0.0.1

Available ProjectDiscovery FOSS Tools

1. subfinder v2.5.0 (latest)
2. cloudlist - v1.0.0 (latest)
3. dnsx - v1.0.9 (outdated)
4. uncover - v1.0.0 (latest)
5. naabu - (not installed)
6. httpx - (not installed)
7. nuclei - (not installed)
8. notify - (not installed)
9. proxify - (not installed)
10. interactsh-client - (not installed)
11. interactsh-server - (not installed)
12. chaos-client - (not installed)
13. mapcidr - (not installed)
14. simplehttpserver - (not installed)
```

```bash
pdtm -h
projectdiscovery foss tool manager

Usage:
  pdtm [flags]

Flags:
INSTALL:
   -i, -install string[]  	install given pd-tool (comma separated)
   -ia, -install-all		install all pd-tools (false)

UPDATE:
   -u, -update string[]		update given pd-tool (comma separated)
   -ua, -update-all    		update all pd-tools (false)

REMOVE:
   -r, -remove string[]   	remove given pd-tool (comma separated)
   -ra, -remove-all   		remove all pd-tools (false)

DEBUG:
   -silent   	disable banner in output
   -version  	display version of the project
   -v        	display verbose output
```
```bash

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