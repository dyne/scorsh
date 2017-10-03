# scorsh

**Signed-Commit Remote Shell**


**scorsh** lets you trigger commands on a remote **git** server through commits, optionally signed with **gnupg**.

**scorsh** is written in Go. 


## Why scorsh

...if you have ever felt that git hooks fall too short to your standards...

...because you would like each specific push event to trigger _something
different_ on the git repo...

..and you want only authorised users to be able to trigger that
_something_...

...then **scorsh** might be what you have been looking for. 

**scorsh** is a simple system to execute commands on a remote host by
using git commits containing customisable commands
(scorsh-tags) that can be authenticated using a gnupg signature . **scorsh** consists of three components:

* the `scorsh-commit` executable (client-side)

* a `post-receive` git hook

* the `scorshd` binary itself (server-side)

The `scorsh-commit` executable is used to inject scorsh-commands in a
regular gpg-signed git commit. 

For each new push event, the `post-receive` hook creates a file in a
configurable spool directory, containing information about the repo,
branch, and commits of the push.

The `scorshd` binary processes inotify events from the spool, parses
each new file there, walks through the new commits looking for signed
ones, checks if the message of a signed commit contains a recognised
scorsh-command, verifies that the user who signed the message is
allowed to use that scorsh-command, and executes the actions
associated to the scorsh-command. 

The set of scorsh-commands accepted on a repo/branch is configurable,
and each scorsh-command can be associated to a list of
actions. Actions are just URLs, at the moment restricted to two
possible types:

* `file://path/to/file` - in this case `scorsh` tries to execute the
  corresponding file (useful to execute scripts)
  
* `http://myserver.com/where/you/like` - in this case `scorsh` makes an
  HTTP request to the specified URL (useful to trigger other actions,
  e.g., Jenkins or Travis builds -- **currently not working**)
  


## Build notes

**scorsh** depends on the availability of a native build of `libgit2`
version `0.26` or greater on the native system where ***scorsh** is
built. This dependencies is easily satisfied on various operating
systems by using their respective package manager. For instance in
Devuan ASCII one can simply do:

```
sudo apt install libgit2-dev
```

In most distributions unfortunately `libgit2` is older than `0.26` so
one should first build this exact release version from source,
available
here:
[https://github.com/libgit2/libgit2/releases/tag/v0.26.0](libgit2 release 0.26)

Then proceed installing dependencies for **scorsh**:
```
make deps
```

And finally build its binary:
```
make
```

## Configuration walkthrough (DRAFT)

`scorshd` reads its configuration from a yaml file, normally passed on
the command line through the option `-c CFG_FILE`. An example is the
following:

```
---
s_spooldir: "./spool"
s_logfile: "./scorsh.log"
s_logprefix: "[scorsh]"

s_workers:
  [
     {
       w_name: worker1,
       w_repos: [".*:.*"], # All branches in all repos
       w_folder: ./worker1,
       w_logfile: ./worker1/worker1.log,
       w_cfgfile: "./worker1/worker1.cfg",
     },
     {
       w_name: worker2,
       w_repos: [".*:master"], # Branch master in all repos
       w_folder: ./worker2,
       w_logfile: ./worker2/worker2.log,
       w_cfgfile: "./worker2/worker2.cfg",
     }
]
...

```

This files defines two workers. Each worker is associated to a pair of
`repo:branch` regexps. A worker will be activated only on pushes made
on a matching `repo:branch`. Each worker has a configuration file
`w_cfgfile`, where the list of accepted scorsh-commands is
defined. For instance, for `worker1` we could have:

```
---
w_commands:
    [
     {
       c_name: "LOG",
       c_keyrings: ["allowed_users.asc"],
       c_actions: [
                    {
                     a_url: "file:///home/katolaz/bin/scorsh_script_log.sh"
                    }
                   ]
      },
     {
       c_name: "build",
       c_keyrings: ["allowed_users.asc"],
       c_actions: [
                    {
                     a_url: "file:///home/katolaz/bin/scorsh_script.sh",
                     a_hash: "c129d4a12998c44dfb9a9fd61ec3159bf29606e0f7280f28bbd98fc6f972fa27"
                    }
                   ]
      },
     {
      c_name: "preview",
      c_keyrings: ["allowed_users.asc"],
      c_actions: [
                  {
                  a_url: "file:///home/katolaz/bin/scorsh_preview.sh"
                  }
                 ]
     }
      
    ]
...
```

In this example, `worker1` has three configured scorsh-commands,
namely `LOG`, `build`, and `preview`.  Commands are
*case-sensitive*. Each command is associated to a list of keyblocks
(containg the public keys of the users allowed to run that command),
and to a list of actions. 

**TBC**

## License

**scorsh** is Copyright (2017) by Vincenzo "KatolaZ" Nicosia.

**scorsh** is free software. You can use, modify, and redistribute it
  under the terms of the GNU Affero General Public Licence, version 3
  of the Licence or, at your option, any later version. Please see
  LICENSE.md for details.
