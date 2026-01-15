# Collingo CLI Developer Tool

This is the CLI tool to work with Collingo using your console.

## Install

When you are using Go:

```
go install github.com/MatthiasSchild/collingo@latest
```

## Authenticate

Before you can use it, you need to login your device to
your account. Visit [Collingo](https://collingo.app),
login to your account and go to settings.

There you can find your access tokens. Create one and copy it.
Afterwards use the following command:

```
collingo login
```

Now you will be asked to enter your access token.
After an successful login, you can use the `collingo` command.

This will create a `.collingo` file in your home directory.
This contains the meta data for using collingo on your device.

## First steps

First, you need to create a **project** to work on.
Your free plan contains until 5 projects you can create.
To create a project, you can either create it using the
web dashboard or the following command:

```
collingo projects create
```

Afterwards, you need a **workspace**. A workspace is your code,
which is connected to a Collingo project. E.g. when you create a
web app, your web app code is your workspace, and you manage your labels
and translations in a Collingo project you've just created.

To initialize a workspace, navigate to the directory of your code and
run `collingo init`, you will be asked what project should be used for this workspace.

```
cd ../my-website
collingo init
```

It will create a `.collingo.json` file, which contains the metadata for your workspace.
This contains workspace-related data, and no user related data, and can therefore
be also commited, if you are using git.

Now you can use all commands. usually you can either provide flags (e.g. `--technical-name demo`),
or when you don't provide a flag, you will be interactively asked. Use `--help` to see what
flags are offered.

```
# Create a group
collingo groups create --technical-name demo
# Create an entry
collingo entries create --group demo
```

## Documentation

I try to do my best to provide you with a documentation,
also containing how to use the CLI tool.
You can find it in the [Collingo Docs](https://docs.collingo.app).

But you can also always use `--help` to display more information about
the command, e.g:

```
collingo groups create --help
```
