# azureprice

![build workflow](https://github.com/muandane/azureprice/actions/workflows/release.yml/badge.svg)

cloudcost is a Go CLI that retrieves pricing information for Azure services using the Azure pricing API.

![Azureprice Demo](./demo.gif)

## Features

- Retrieves pricing info for Azure services
- Outputs pricing in easy-to-read tables
- Built using [lipgloss](https://github.com/charmbracelet/lipgloss) for beautiful terminal output
- Ability to calculate the monthly cost of running a service or it's pricing depending on Bandwidth and Region.

## Installation

```sh
brew tap muandane/gitmoji
brew install cloudcost
```

Or build from source:

```sh
git clone https://github.com/muandane/cloudcost
cd cloudcost
go build -o cloudcost ./src/cmd  
```

## Usage

```sh
cloudcost is a Go CLI that retrieves pricing information for Azure services using the Azure pricing API.

Usage:
  cloudcost [command]

Available Commands:
  azure       Main command for Azure resource management.
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Display application version information.

Flags:
  -h, --help   help for cloudcost

Use "cloudcost [command] --help" for more information about a command.
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0)
