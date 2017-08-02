# scorsh

Signed-Commit Remote Shell


**scorsh** lets you trigger commands on a remote git server through
signed git commits.

**scorsh** is written in Go. 


## WTF

...if you have ever felt that git hooks fall too short to your standards...

...because you would like each specific push event to trigger _something
different_ on the git repo...

...and you want only authorised users to be able to trigger that
_something_...

...then **scorsh** might be what you have been looking for. 

**scorsh** is a simple system to execute commands on a remote host by
using GPG-signed commits containing customisable commands
(scorsh-tags). **scorsh** consists of two components:

* a `post-receive` git hook

* the `scorsh` binary itself

For each new push event, the `post-receive` hook creates a file in a
configurable spool directory, containing information about the repo,
branch, and commits of the push.

The `scorsh` binary processes inotify events from the spool, parses
each new file there, walks through the new commits looking for signed
ones, checks if the message of a signed commit contains a recognised
scorsh-tag, verifies that the user who signed the message is allowed
to use that scorsh-tag, and executes the commands associated to the
scorsh-tag. Or, well, this is what `scorsh` should be able to do when
it's finished ;-)

The set of scorsh-tags accepted on a repo/branch is configurable, and
each scorsh-tag can be associated to a list of commands. Commands are
just URLs, at the moment restricted to two possible types:

* `file://path/to/file` - in this case `scorsh` tries to execute the
  corresponding file (useful to execute scripts)
  
* `http://myserver.com/where/you/like` - in this case `scorsh` makes an
  HTTP request to the specified URL (useful to trigger other actions,
  e.g., Jenkins or Travis builds...)
  



## Build notes

**scorsh** depends from the availability of a native build of
`libgit2` version `0.25` or greater on the native system where
***scorsh** is built. This dependencies is easily satisfied on various
operating systems by using their respective package manager. For
instance in Devuan ASCII one can simply do:

```
sudo apt install libgit2-dev
```

In Devuan Jessie unfortunately `libgit2` is older than `0.25` so one
should first build `git2go` from its repository, in which `libgit2` is a
submodule to be built from scratch.

```
git clone https://github.com/libgit2/git2go
cd git2go
git submodule init
git submodule update
cd libgit2
cmake .
make
sudo make install
```

Then proceed installing dependencies for **scorsh**:
```
make deps
```

And finally build its binary:
```
make
```


## License

**scorsh** is Copyright (2017) by Vincenzo "KatolaZ" Nicosia.

**scorsh** is free software. You can use, modify, and redistribute it
  under the terms of the GNU Affero General Public Licence, version 3
  of the Licence or, at your option, any later version. Please see
  LICENSE.md for details.
