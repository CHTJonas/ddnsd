# Dynamic DNS Daemon

This repo hosts the code for a small dynamic DNS daemon (ddnsd) that I use for some personal infrastructure. It's written in Go and designed to be as small and as simple as possible. Host machines poke ddnsd over HTTP by making a POST request to an API endpoint. HTTP Basic authentication is used with credentials stored in a .htpasswd file. The username dictates which record to update which makes sure that API credentials can't be shared.

## Usage

TODO

## Installation

TODO

## Copyright

ddnsd is licensed under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).

Copyright (c) 2019–2021 Charlie Jonas.
