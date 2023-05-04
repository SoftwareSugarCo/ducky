# Ducky
![ducky2.png](assets%2Fducky2.png)

## General info and links


Ducky is a simple CLI tool for querying and chatting with an OpenAI Model; the purpose of this project is to limit
distractions and time wasted going through a browser or other UI to query a model. Since the most likely user of ducky 
is a software engineer, it is designed to be used in a terminal and to extract out, and format, any code snippets it detects.

[Why Ducky?](#why-ducky)

[Install GoLang](#installing-golang)

[Installing Ducky](#installing-ducky)

[Configuring Ducky](#configuring-ducky)

[Using Ducky](#using-ducky)


## Why Ducky?
One of the first - and very important - things a software engineer learns to do is to debug their code.
Sometimes this is an excruciating process, and sometimes it is a simple process. We've all had that one 
project or class that we just couldn't figure out and spend hours or days trying to debug.

It truly does help to talk it out! This is where the idea of Rubber Ducky Debugging comes in; you can think of
Ducky as your virtual or AI rubber ducky that can not only listen to you (via your textual input) but also respond
with helpful insights and suggestions.

Links:

[Rubber Ducky Debugging Wiki](https://en.wikipedia.org/wiki/Rubber_duck_debugging)

[Rubber Ducky Debugging](https://rubberduckdebugging.com/)


## Installing Golang
Ducky is written in GoLang, so you will need to install GoLang on your machine. If you already have GoLang installed,
you can completely skip this section.

### Windows

1. Download the Go installer for Windows from the official website: https://golang.org/dl/.
2. Run the installer and follow the instructions.
3. After installation, open the command prompt and type `go version` to verify that Go has been installed successfully.

### macOS

1. Install Homebrew by opening the terminal and running the following command:
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```
2. Once Homebrew is installed, open the terminal and run the following command to install Go:
```
brew install go
```
3. After installation, type `go version` in the terminal to verify that Go has been installed successfully.


### Linux

1. Open the terminal and run the following commands to install Go:
```
sudo apt-get update
sudo apt-get install golang-go
```
2. After installation, type `go version` in the terminal to verify that Go has been installed successfully.

## Configuring GOPATH

Go requires a GOPATH environment variable to be set. This is the path where Go packages are installed on your machine. Here's how you can set it:

### Windows

1. Right-click on My Computer and select Properties.
2. Click on Advanced System Settings.
3. Click on Environment Variables.
4. Under System Variables, click on New.
5. Enter "GOPATH" for the variable name and the path where you want to install your Go packages for the variable value.
6. Click OK.

### macOS and Linux

1. Open the terminal.
2. Run the following command to open the `.bashrc` file:
```
nano ~/.bashrc
```
3. Add the following line to the file, replacing `/path/to/gopath` with the path where you want to install your Go packages:
```
export GOPATH=/path/to/gopath
```
4. Press `Ctrl+X`, `Y`, and then `Enter` to save the changes.

## Installing Ducky

Ducky (and other Go packages) can be installed globally using the `go install` command. Here's how:

1. Open the terminal.
2. Run the following command replacing `<VERSION>` with the version you want to install:
```
go install github.com/SoftwareSugarCo/ducky@<VERSION>
```
3. The package will be installed globally and can be accessed from anywhere within your terminal using the `ducky` command.

> ***RECOMMENDATION:*** *Ducky has two modes, chat and query, and each can be configured with a specific model via a flag. If you have 
a mode or model you always want to use then you should create an alias in your shells config file.*
> 
> For example, if you are using ZSH, you can add the following line to your `.zshrc` file:
> ```
> alias ducky="ducky chat --chat_model=gpt4"
> ```

## Configuring Ducky
### Getting your OpenAI API Key
Ducky requires an OpenAI API key to be set in order to query models. You can get an API key by signing up for an OpenAI account and creating an API key. Here's how you can set it:

1. Sign up for an OpenAI account at https://beta.openai.com/.
2. Once you have an account, go to https://beta.openai.com/account/api-keys.
3. Click on the Create New API Key button and create a new API key.

### Setting your OpenAI API Key for Ducky to use
Ducky is looking for a global environment variable called `DUCKY_API_KEY` to be set. Here's how you can set it:

1. Open the terminal.
2. Edit your terminal's config file. For example, if you are using ZSH, run the following command:
```
nano ~/.zshrc
```
3. Add this line to the bottom of your config file:
```
export DUCKY_API_KEY=<YOUR_OPENAPI_KEY>
```
4. Save the config file however your editor requires.
5. Restart your terminal or `source` the config file to apply the changes.

## Using Ducky
### Ducky Chat
Ducky Chat is a simple CLI tool for chatting with an OpenAI Model. It is designed to be used in a terminal and to extract out, and format, any code snippets it detects.
This mode keeps the history of your chat saved while chatting so the model can "remember" what you've said.

#### Flags:
* `--chat_model` - The model to use for chatting. Defaults to `gpt4`.

#### Chat Commands:
* `/done`, `/exit`, `/quit` - Exits the chat.
* `/m`, `/multi`, `/multiline` - Enters multiline mode. This mode allows you to enter multiple lines of text before sending it to the model.
    * End multiline mode by entering `/end`

#### Examples:
```shell
$ ducky chat --chat_model=gpt4
Using model: gpt4
Ducky: Yes, How may I help you?
You: <question>
Ducky: <response>
You: <question>
Ducky: <response>
You: /m
<enter multiple lines of text>

/end
Ducky: <response>
You: /done
```

### Ducky Query
Ducky Query is a simple CLI tool for querying an OpenAI Model. It is designed to be used in a terminal and to extract out, and format, any code snippets it detects.
This mode has no memory of previous queries, so each query is treated as a new query and the model has no context.

#### Flags:
* `--query_model` - The model to use for querying. Defaults to `gpt4`.

#### Examples:
```shell
ducky query --query_model=gpt3.5Turbo "<question>"
```