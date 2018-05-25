# Chalepoxenus Kutteri

[![Build Status](https://travis-ci.org/containous/kutteri.svg?branch=master)](https://travis-ci.org/containous/kutteri)
[![Docker Build Status](https://img.shields.io/docker/build/containous/kutteri.svg)](https://hub.docker.com/r/containous/kutteri/builds/)

## Usage

```
Chalepoxenus Kutteri: Track a GitHub repository and publish on Slack.

Usage: kuterri [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: kuterri [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Available Commands:
        version                                            Display the version.
Use "kuterri [command] --help" for more information about a command.

Flags:
    --bot-icon      Icon of the bot in Slack.          (default ":bento:")
    --bot-name      Name of the bot in Slack.          (default "Bender")
-c, --channel       Slack channel. [required]          
    --dry-run       Dry run mode.                      (default "true")
    --ghtoken       GitHub Token.                      
-o, --owner         Repository owner. [required]       
    --port          Server port.                       (default "80")
-r, --repo-name     Repository name. [required]        
-f, --search-filter Search filter. [required]          
    --server        Server mode.                       (default "false")
    --sltoken       Slack Token. [required]            
-h, --help          Print Help (this message) and exit 
```

## Information

| issue   | filter                                                           |
|---------|------------------------------------------------------------------|
| new     | `type:issue state:open created:>[DATE] in:title traefik`         |
| updated | `type:issue state:open updated:>[DATE] in:title traefik` - [NEW] |
| ~closed~  | `type:issue closed:>[DATE] in:title traefik`                     |


| PR      | filter                                                        |
|---------|---------------------------------------------------------------|
| new     | `type:pr state:open created:>[DATE] in:title traefik`         |
| updated | `type:pr state:open updated:>[DATE] in:title traefik` - [NEW] |
| merged  | `type:pr merged:>[DATE] in:title traefik`                     |
| ~closed~  | `type:pr closed:>[DATE] in:title traefik`                     |
