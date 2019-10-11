# remod

remod is a tool to work with local copies of libraries and go modules.
It provides a cli to manage replacement directives in a file called `go.mod.dev`
that can be ignored from VSC.

remod also provide `remod go args...` to wrap a standard go command and transparently combining
the original `go.mod` with the `go.mod.dev`, executing the command, and restoring the original file.

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

## Developement mode

remod allows to switch on and off a development mode, where it will
add replace directives in a `go.mod.dev` file.

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

The `go.mod.dev` file will now look like:

```mod
replace (
  github.com/spf13/viper => ../viper
)
```

To to delete the `go.mod.dev`, run:

```shell
remod off
```

The modules you pass are actually used as prefix, so you can replace all the packages from `github.com/spf13`
by doing:

```shell
remod on -m github.com/spf13
```

Which will modify the `go.mod.dev` file like so:

```mod
replace (
  github.com/spf13/cobra => ../cobra
  github.com/spf13/viper => ../viper
)
```

To set a different base path, you can use the option
`--prefix`.

For instance:

```shell
remod on -m github.com/spf13/viper --prefix github.com/me/ --replace-version dev
```

Which will turn the `go.mod.dev` file to:

```mod
replace (
  github.com/spf13/viper => github.com/me/viper dev
)
```

## Wraping go command

If you run a classic go command, the `go.mod.dev` will of course be ignored.
You can wrap the go command using `remod go` so `go.mod.dev` will be used.

For example:

```shell
remod go build
```

```shell
remod go test -race ./...
```

If there is no `go.mod.dev`, `remod go` will simply run the go command, so it will always work.

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
