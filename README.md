# scorsh
Signed-Commit Remote Shell


**scorsh** lets you trigger commands on a remote git server through
signed git commits.

**scorsh** is written in Go. 

**This is still work-in-progress, not ready to be used yet**

# WTF

...if you have ever felt that git hooks fall too short to your standards...

...because you would like each specific push event to trigger _something
different_ on the git repo...

...and you want only authorised users to be able to trigger that
_something__....

..then **scorsh** might be what you have been looking for. 

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
scorsh-tag. Or, well, this is what `scorsh` will do when it's ready.

The set of scorsh-tags accepted on a repo/branch is configurable, and
each scorsh-tag can be associated to a list of commands. Commands are
just URLs, at the moment restricted to two possible types:

* file://path/to/file - in this case `scorsh` tries to execute the
  corresponding file (useful to execute scripts)
  
* http://myserver.com/where/you/like - in this case `scorsh` makes an
  HTTP request to the specified URL (useful to trigger other actions,
  e.g., Jenkins or Travis builds...)
  






