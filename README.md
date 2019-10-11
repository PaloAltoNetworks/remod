# remod

Remod is a tool to work with local copies of libraries and go modules.
It provides a cli to manage replacement directives and uses git attributes
to make the changes invisible from git's point of view.

When you work on projects with various internal libraries that are working
closely together, you usually need to have a bit of velocity. For example, the Aporeto ci pipelines
are able combine multiple pull requests accross multiple github repositories to test them together.

While go modules help in a lot of scenarios, being able to do such things is not one them. remod
made that workflow easier by selectively restoring the GOPATH behavior on only a subset
of the dependencies, while benefiting from go modules for the others.

## Installation

To install remod, run:

```shell
go get go.aporeto.io/remod
```

## Init

If you want to use remod on a repo that does not handle it already, you can run:

```shell
remod init
```

This will create a `.gitattributes` file and will add the correct git configuration in `./.git/config`

> NOTE: if the git attributes file already exists, remod will not do anything but will print
> a command to append what it needed.

## Developement mode

remod allows to switch on and off a development mode, where it will
add replace directives to your `go.mod` file. These changes will be invisible to git so you cannot commit them.

Your replacements will survive accross checkouts (branching, etc.) and will not interfere with other configurations from other people.

You can replace one or multiple libraries, pointing to a local fork
or to a remote one.

You can write manual replacements, as long as they are between the special comments:

```mod
// remod:replacements:start

<manual replacements here>

// remod:replacements:end
```

For instance, if we have the following `go.mod`:

```mod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)
```

To use a local copy of viper, run:

```shell
remod on -m github.com/spf13/viper
```

The `go.mod` file will now look like:

```mod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)

// remod:replacements:start

replace github.com/spf13/viper => ../viper

// remod:replacements:end
```

You should see no change when doing `git status`.

To turn it off and reset the `go.mod` file, run:

```shell
remod off
```

The modules you pass are actually used as prefix, so you can replace all the packages from `github.com/spf13`
by doing:

```shell
remod on -m github.com/spf13
```

Which will modify the `go.mod` file like so:

```mod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)

// remod:replacements:start

replace (
  github.com/spf13/cobra => ../cobra
  github.com/spf13/viper => ../viper
)

// remod:replacements:end
```

To set a different base path, you can use the option
`--prefix`.

For instance:

```shell
remod on -m github.com/spf13/viper --prefix github.com/me/ --replace-version dev
```

Which will turn the `go.mod` file to:

```mod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)

// remod:replacements:start

replace github.com/spf13/viper => github.com/me/viper dev

// remod:replacements:end
```

## Updating modules

remod allows to simply perform batch modules updgrade.

For instance, to update viper:

```shell
remod up -m github.com/spf13/viper --version master
```

The given modules is also matched on a prefix, so to update all modules from spf13:

```shell
remod up -m github.com/spf13 --version master
```

By default, `remod up` targes the working directory, but you can target one or more folders:

```shell
remod up -m github.com/spf13 --version master /path/to/my/module1 /path/to/my/module2
```

Or simply do it recursively from one folder:

```shell
remod up -m github.com/spf13 --version master /path/to/my/modules -r
```
