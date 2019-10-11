# remod

Remod is a tool making it easier to work with local copy of some library when using go modules.
It provides a cli to easily manage replacement and can also be used alongside a a git attribute file
to make the change invisible to git.

## Installation

To install remod, run:

```shell
go get go.aporeto.io/remod
```

## Usage

### Init

If you want to use remod on a repo that does not handle it already, you can run

```shell
remod init
```

This will create a .gitattributes file and will add the correct git configuration.

### Developement mode

For instance, if we have the follwoging `go.mod`:

```gomod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)
```

You can use remod to use a local copy of viper, you can run

```shell
remod on -m github.com/spf13/viper
```

The go mod will look like:

```gomod
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

You should also see not changes when doing `git status`.

To turn it off, run

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

```go.mod
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
