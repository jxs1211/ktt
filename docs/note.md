## Build Guidelines

### Prerequisites

- Go (latest version)
- Node.js >= 16
- NPM >= 9

### Install Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Pull the Code

```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### Build Frontend

```bash
npm install --prefix ./frontend
```

or

```bash
cd frontend
npm install
```

### Compile and Run

```bash
wails dev
```

## About

### Wechat Official Account

<img src="docs/images/wechat_official.png" alt="wechat" width="360" />

### Sponsor

If this project helpful for you, feel free to buy me a cup of coffee ☕️.

- Wechat Sponsor

<img src="docs/images/wechat_sponsor.jpg" alt="wechat" width="200" />

Cli Tab
// Ok, let's change it a bit, I login in the terminal when I switch to the Console tab, like what it did in @ContentCli.vue for now
// I want to seperate the frontend of gotty from backend, and move it into the @ContentCli.vue, user can operate cmd in the terminal.
// let's do it step by step
// - 1 init connnection, @terminal_service.go is responsible for manage all connections, if connecion is build ok, save connection info to sqilte3,include port, address, cmd, start_time
// - manage connection, every connnection will be alive for 2 hours, there is a job for maintain the all the connections, how to support reconnection
// - From frontend perspective, every time need to start a new session with info includes: address,port,writable,cmd
// - Backend start a gotty server with info sent by frontend
// - Cmd:
// - if create session from Console tab, cdebug to the default worker pod
// - if coming from item of diagnose tab
// - pod
// - running container: connect to it directly
// - failed container: return tlsConfig, nil

// ai client
// In preference ai options, check user ai condfiguration using like `func (c *OllamaClient) Configure(config IAIConfig) error {}`
// the available models in the chat is original from configured models in Preference ai options

- is there is not configured model, pop up preferences dialog to enable one
// save the configuration to preference.yaml, reload it to map of ai/client_service.go when KT reboot, map goes like `map[string]ai.IAI`, key goes like `provider-model`
//
todo:
- load and configure saved ai models when bootup app for ai client(async loadPreference)
- validate and cache all supported ai model providers when saving prefernece at PreferenceDialog
- sync validated ai model providers to chat input's model options


Let's combine the two components @ContentErrorPane.vue and @CliBar.vue together, and implement an amazing functionality step by step here is the thing:
- when I click btn debugWithAI, the @CliBar.vue will popup from bottom to the half height of the tab pane, the data-table on it's top will shrink to half height
- the session name is composed of `resource kind: resource name:session.address : session.port : session.cmds`
- The clicked item's info will sent as init prompt to the chat box
- After the response msg is received, there are 2 options(Copy and Apply) can be shown up at the top right of the msg box when the mouse is hoved on any place on the response msg box, I can copy it and paste to the right terminal to be executed. And I can also click Apply btn to execute the command in the response box in the right terminal
- The popuped CliBar is draggable, and can be closed by clicking the close button on the top right of the CliBar, and when it's displayed, the data-table in @ContentErrorPane.vue is also visible(include the pagination on the bottom right), I can opreate it normally

Let's continue to implement @CliBar.vue 's functionality:
- I want the content of the chat box and the terminal can be keeped, when I swith back from others tab to current tab, so that I can continue to operate in the chat and terminal with previous context


### CliBar:

#### Issue:
- the current context of chat and terminal disappeared when switch back from others tab to current tab
- Run button is not working(feat/cli-bar-run-btn)
- the input can't be fixed in the bottom of chat container 


!-- 
cdebug exec --namespace=test -it --privileged --image=nixery.dev/shell/vim/ps/tshark/kubectl/zsh pod/my-distroless

-->
