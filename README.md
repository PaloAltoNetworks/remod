# remod

Remod is a tool making it easier to work with local copies of some library when using go modules.
It provides a cli to manage replacement and can also be used alongside a git attribute file
to make the changes invisible from git's point of view.

## Installation

To install remod, run:

```shell
go get go.aporeto.io/remod
```

>  use `go get github.com/aporeto-inc/remod` until vanity url is available

## Init

If you want to use remod on a repo that does not handle it already, you can run

```shell
remod init
```

This will create a `.gitattributes` file and will add the correct git configuration.

> NOTE: if the git attributes already exists, remod will not do anything but will print
> a command to append what it needed.

## Developement mode

remod allows to switch on and off a development mode, where it will
add replace directives to your go.mod that you cannot commit upstream.
You can replace one or multiple libraries, either pointing to a local fork
or using remote one.

For instance, if we have the following `go.mod`:

```mod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)
```

If you want to use a local copy of viper, you can run:

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

You should also see no change when doing `git status`.

To turn it off, and reset the `go.mod` file, run:

```shell
remod off
```

The modules you pass are actually used as prefix, so you can replace all package from `github.com/spf13`
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

If you want to set a different path instead of using `../` you can use the option
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

remod also allows to simply perform bunch modules updgrade.
For instance to update viper:

```shell
remod up -m github.com/spf13/viper --version master
```

The given modules is also matched on a prefix, so to update all modules from spf13:

```shell
remod up -m github.com/spf13 --version master
```

By default `remod up` will work locally, but you can pass one or more folders:

```shell
remod up -m github.com/spf13 --version master /path/to/my/module1 /path/to/my/module2
```

Or simply do it recursively from one folder:

```shell
remod up -m github.com/spf13 --version master /path/to/my/modules -r
```
