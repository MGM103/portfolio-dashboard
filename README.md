# Portfolio Dashboard TUI 📊

### Description

Portfolio dashboard is a Terminal UI application for seeing the value of one's digital asset portfolio.

Currently each ecosystem has tools for ascertaining portfolio information and other data points. Therefore in the current landscape one must use a combination of tools to get a holistic valuation of a portfolio of digital assets. These tools are mostly web based applications that use various on-chain indexers to gather information and can feel quite laggy and slow to use. This makes the task of ascertaining portfolio details laborious, tedious and unpleasant.

Enter portfolio dashboard, this TUI runs locally and stores your data on your own device. It is quite simple in design, digital asset positions are added manually, the valuation of the portfolio is then able to be calculated in a moments notice. There is also a watchlist feature to monitor assets that are not yet owned.

Currently the data for digital asset pricing is sourced from [Coin Market Cap](https://coinmarketcap.com/).

## Table of contents

- [Portfolio Dashboard 📊](#portfolio-dashboard-tui)
  - [Description](#description)
  - [Table of contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Requirements](#requirements)
    - [Quickstart](#quickstart)
  - [Usage](#usage)
    - [Local Dev](#local-dev)
    - [Installing on System](#installing-on-system)
  - [Testing](#testing)
    - [Unit tests](#unit-tests)
    - [Test Coverage](#test-coverage)
  - [Acknowledgements](#acknowledgements)

## Getting Started

### Requirements:

The following must be installed on your machine:

- [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git/)
- [go](https://go.dev/doc/install)
- [CMC api key](https://coinmarketcap.com/api/)

### Quickstart:

```
git clone git@github.com:MGM103/portfolio-dashboard.git
cd portfolio-dashboard
make
```

## Usage

To run the application you must create a `.env` file with a `CMC_API_KEY` as seen below:

```sh
CMC_API_KEY=<your-api-key>
```

### Local Dev

To run the application in dev you can use the following commands:

```sh
make
```

or

```sh
make run
```

### Installing on System

To use the commands `portfolio-dashboard` & `portfolio-app` which are the application executable and shell script respectively you can run the following command:

```sh
make install
```

This command will create a config and binary directory if it does not exist and copy the binary and shell script into the bin directory and the db and .env dependencies in the config directory.

## Testing

### Unit tests

To run the unit tests for this project you can run the following commands:

```
cd ./data/
go run test
```

## Acknowledgements

This project was built using [bubble tea](https://github.com/charmbracelet/bubbletea) TUI application framework, leveraging [bubbles](https://github.com/charmbracelet/bubbles) the component library.

Additionally, this [tutorial](https://www.youtube.com/watch?v=_gzypL-Qv-g&t=1s) from [package main](https://www.youtube.com/@packagemain) was a great aid in learning how to structure and build the application.
