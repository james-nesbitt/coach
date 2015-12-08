# .coach

This folder keeps settings and scripts for use with coach.  This folder
demarks the root of the project, and allows the coach CLI to be run at
any child path.  The folder is usually kept next to a source tree in order
to allow git separation, or encapsulation, where desired.

## configuring coach

### conf.yml

Base coach configuration, defining the project name (used to name containers)
, configure the docker client, add any custom paths that you may want to use
and to declare initial tokens

### secrets/secrets/yml

This YAML file contains additional tokens, loaded after the conf.yml tokens
are added, to allow sensitive tokens to be kept out of any git repository.
The file could be shared as an example.secrets.yml in repository.

Not demonstrated in this demo, is that if you have a file in your home directory
~/.coach/secrets/secrets.yml, then those tokens are loaded after these tokens
to allow user overrides.

### nodes.yml

A map of nodes which coach will run.

### tools.yml

A map of non-coach tools, usually scripts, that coach can try to run for you.

This allows for running of scripts from anywhere in the project, and allows
passing in coach tokens as ENV variables.

This accomodates the need to allow simple scripting integration, where writing
custom operations does not make sense.

### help.yml

A map of custom help topics, that could be accessed using $/> coach help {token}

### docker/{build}

An optional path that can be used to keep Docker builds that should be built
automatically as a part of a node.  To use one of these, there must be a node
that declares Build: docker/{build}
