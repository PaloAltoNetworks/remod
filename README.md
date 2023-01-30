# remod

> DEPRECATED: Go devs, after dismissing the use case remod was solving, finally
> implemented a native way to do what remod was doing with
> [go workspaces](https://go.dev/blog/get-familiar-with-workspaces).
> So this repository is now deprecated and won't get any new updates.

remod is a tool to work with local copies of libraries and go modules.
It provides a cli to manage replacement directives in a file called `remod.dev`
that can be ignored from VSC. Remod also uses git attributes to make the changes
applied to the `go.mod` file completely transparent.

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

## Working with remod

### Initialization

If you never used remod on that repo, you need to
first run:

```shell
remod init
```

This will add a the necessary `.gitattributes`, and configure your
clone to ignore the relevant changes in the `go.mod` and `go.sum` files.

You can then edit the `remod.dev` file with the replacements you want.
You can also decide to preconfigure the replacements using the `--include` and `--exclude` flags.

The modules you pass are actually used as prefix, so you can replace all the packages from `github.com/spf13`
by doing:

```shell
remod init -m github.com/spf13
```

Which will modify the `remod.dev` file like so:

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
remod init -m github.com/spf13/viper --prefix github.com/me/ --replace-version dev
```

Which will turn the `remod.dev` file to:

```mod
replace (
  github.com/spf13/viper => github.com/me/viper dev
)
```

It is safe to edit the `remod.dev` file manually as you wish.

### Activating remod for the branch

> Note: If you've cloned a fresh repo, or switched to a new branch,
> you need to run `remod on` again before being able to use it.
> Running `remod on` is idempotent and will align what needs to be aligned.

remod allows to switch on and off the development mode at will, where it will
add replace directives from `remod.dev` file in your `go.mod` file.

For instance, if we have the following `go.mod`:

```mod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)
```

And `remod.dev` contains:

```mod
replace github.com/spf13/viper => ../viper
```

You can enable the development mode by running:

```shell
remod on
```

The `go.mod` file will now look like:

```mod
module go.aporeto.io/remod

go 1.12

require (
  github.com/spf13/cobra v0.0.5
  github.com/spf13/viper v1.4.0
)

// remod:start

replace github.com/spf13/viper => ../viper

// remod:end
```

You should see no change from the git point of view.

### Updating or installing new modules

If you need to update or add a new module that must end up in the final `go.mod` you need to
run the following command:

```shell
remod get github.com/user/repo
```

This will transparently restore the original `go.mod`, run the add the new dependency, and reactivate
the development mode. If development mode was not active, it will work as usual.

`remod get` will blindly pass all arguments to the underlying `go get` command, so anything supported by
go get command can be done through remod.

> Note: if you simply run go get while remod is on, you will loose the change.

### Deactivating remod for the branch

You cam turn off the development mode for the current branch by running:

```shell
remod off
```

This will restore the original `go.mod` and `go.sum` file.
Note that this only affects the current branch. You can run `remod on` again
at any time to start development mode again.

### Updating modules in bulk

The command `remod update` allows to update all modules with a prefix at once.

For instance, to update all modules from `go.aporeto.io` in the current project, you can run:

```shell
remod up go.aporeto.io
```

It will udpate the the `latest` tag by default.

You can change this behavior by setting the `--version` flag, and exclude some modules with
the `--exclude` flag.

For instance:

```shell
remod up go.aporeto.io --version master --exclude go.aporeto.io/trireme-lib
```
