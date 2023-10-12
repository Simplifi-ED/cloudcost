# azureprice

![build workflow](https://github.com/muandane/azureprice/actions/workflows/release.yml/badge.svg)

azureprice is a Go CLI that retrieves pricing information for Azure services using the Azure pricing API.

![Azureprice Demo](./demo.gif)

## Features

- Retrieves pricing info for Azure services
- Outputs pricing in easy-to-read tables
- Built using [lipgloss](https://github.com/charmbracelet/lipgloss) for beautiful terminal output

## Installation

```sh
brew tap muandane/gitmoji
brew install azureprice
```

Or build from source:

```sh
git clone https://github.com/muandane/azureprice
cd azureprice
go build -o azureprice ./src/cmd  
```

## Usage

```sh
azureprice [flags]
  Usage of azureprice:
    -h, --help            help for azureprice
    -r, --region string   Azure region (default "westus")
    -s, --service string  Azure service
    -t, --type string     Azure VM type (default "Standard_B4ms")
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0)
