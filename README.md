# gconf

An HTTP service for storing versioned configuration data in a Git repository

## Overview

gconf treats all files in the given Git repository as plaintext configuration
files. Each top-level directory in the repository represents a separate
project. Files are read from master/HEAD of the repository. The server updates
the files periodically according to a configurable delay.

TODO(#1): Each project may specify an optional HTTP endpoint for validating
configuration changes. If a change is invalid, gconf does not update its current
snapshot.

## API

### List the set of projects

```sh
GET /projects
```

#### Example Response

```sh
project-1
project-2
```

### Fetch file contents

```sh
GET /file/:project/:filename
```

#### Example Response

```sh
These are the contents of project-1/data.txt
```

## Command Line Usage Usage

```sh
# Serve configs from github.com/owner/repo. Refresh every 5 min.
gconf -repo-owner=owner -repo-name=repo -t=300
```

## Command Line Options

```sh
  -port int
        The port to listen on (default 8080)
  -repo-name string
        The name of the Git repository
  -repo-owner string
        The name of the  Git repository owner
  -t int
        The number of seconds to wait between config updates
```
