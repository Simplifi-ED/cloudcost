# cloudcost

**Query Azure retail pricing from the command line and turn infrastructure cost assumptions into something explicit and reproducible.**

`cloudcost` is a Go CLI for exploring Azure service prices without jumping between pricing pages or copying values into spreadsheets. It queries the Azure Retail Prices API, filters results by service, region, SKU and pricing model, and can calculate usage-oriented cost estimates for common billing units.

| Dark theme | Light theme |
| --- | --- |
| ![cloudcost demo in a dark terminal](./demo-dark.gif) | ![cloudcost demo in a light terminal](./demo-light.gif) |

## Why this exists

Cloud architecture decisions often contain hidden cost assumptions: a region is selected because it is "usually cheaper", a SKU is copied from an old estimate, or hourly pricing is compared without normalising it to expected usage.

`cloudcost` makes those assumptions easier to inspect from the terminal. It is useful during architecture reviews, FinOps investigations, infrastructure design and quick comparisons while working on Terraform or platform changes.

## What it does

- Queries the public Azure Retail Prices API.
- Filters Azure pricing by region, service, SKU/type and pricing model.
- Supports currencies exposed by the Azure pricing API.
- Displays retail price, unit of measure, region, meter and product information in a terminal table.
- Calculates estimates for hourly, monthly, data-volume and event-based billing units.
- Provides shell completion through Cobra.

## Example workflow

Use the Azure calculator command to narrow a pricing query:

```sh
cloudcost azure calculator --region <region> --service <service> --type <sku> --currency EUR
```

Add usage parameters when you want to estimate a cost rather than only inspect the retail unit price:

```sh
cloudcost azure calculator \
  --region <region> \
  --service <service> \
  --currency EUR \
  --days 30 \
  --bandwidth 10
```

Run the built-in help for the exact options supported by your version:

```sh
cloudcost --help
cloudcost azure --help
cloudcost azure calculator --help
```

## Installation

### Build from source

```sh
git clone https://github.com/Simplifi-ED/cloudcost.git
cd cloudcost
go build -o cloudcost .
```

Then move the binary somewhere on your `PATH`.

> The repository previously referenced an external Homebrew tap. Until a maintained package is published from this organization, building from source is the most explicit installation path.

## How it works

`cloudcost` builds a filter from the command-line parameters and queries:

```text
https://prices.azure.com/api/retail/prices
```

Azure can return multiple pages of matching meters. The CLI follows the API pagination links and renders the resulting SKUs and meters as a terminal table.

For hourly resources, the displayed monthly reference price uses 730 hours. Usage estimates are derived from the billing unit returned by Azure, so they should be treated as architecture and comparison aids rather than invoices.

## Scope and limitations

- Azure is currently the supported cloud provider.
- Prices come from the Azure Retail Prices API and may differ from negotiated enterprise pricing, private offers, credits, taxes or contractual discounts.
- Cost calculations are estimates. Validate material financial decisions against the pricing terms that apply to your subscription and contract.
- A returned meter does not prove that a given architecture is valid or available for your specific subscription.

## Development

Requirements:

- Go

Build:

```sh
go build -o cloudcost .
```

Run tests:

```sh
go test ./...
```

## Contributing

Issues and pull requests are welcome. For substantial behavioural changes, open an issue first so the intended pricing model and expected output can be discussed explicitly.

## Maintained by Omnivya

This project is part of the public engineering work maintained by [Omnivya](https://www.omnivya.fr/), focused on cloud, platform engineering, reliability, security and infrastructure automation.

## License

Licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).
