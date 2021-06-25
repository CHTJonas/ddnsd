# Dynamic DNS Daemon

This repo hosts the code for a small dynamic DNS daemon (ddnsd) that I use to setup `_acme-challenge` records when issuing [Let's Encrypt](https://letsencrypt.org/) records for some of my personal infrastructure. It's written in Go and designed to be as small and as simple as possible. It is **not** a full DNS server, rather it parses DNS zonefiles and updates them. To serve your DNS zone you will need a full authoritative server; I use [knot](https://www.knot-dns.cz).

Host machines poke ddnsd over HTTP by making a POST request to an API endpoint during a [dehydrated](https://github.com/dehydrated-io/dehydrated) hook step. HTTP Basic authentication is used with credentials stored in a `.htpasswd` file. The username dictates the record to update which makes sure that API credentials can't be shared.

## Example

Imagine your `.htpasswd` file contains the line `test1.domain.tld.:$2y$05$DKfxz32xyO3giwz0eV1Qi.et.DL2AokHFxjYk1j78Vb7kPJCcmsyi` and that your zonefile has the following contents:

```
domain.tld. 21600 IN SOA server.domain.tld. admin.domain.tld. 2021061800 900 300 4320000 1200
domain.tld. 21600 IN NS server.domain.tld.
test1.domain.tld. 60 IN TXT test1
test2.domain.tld. 60 IN TXT test2
test3.domain.tld. 60 IN TXT test3
```

After running the command `curl http://test1.domain.tld.:testing@server.domain.tld:8080/update -X POST -d 'contents=__Tm_C9ZjVv8RZGi2AXLDWBey1HmTHVsS8c_A7a4hqk'`, your zonefile will now look like this:

```
domain.tld. 21600 IN SOA server.domain.tld. admin.domain.tld. 2021061800 900 300 4320000 1200
domain.tld. 21600 IN NS server.domain.tld.
test1.domain.tld. 60 IN TXT __Tm_C9ZjVv8RZGi2AXLDWBey1HmTHVsS8c_A7a4hqk
test2.domain.tld. 60 IN TXT test2
test3.domain.tld. 60 IN TXT test3
```

Obviously you should substitute `server.domain.tld:8080` for the actual hostname and port that you are using. Make sure that you put ddnsd behind a reverse proxy which handles TLS!

## Usage

```
Usage:
  ddnsd [flags]

Flags:
  -b, --bind string     address and port to bind to (default "localhost:8080")
  -h, --help            help for ddnsd
  -H, --hook string     full path to command/script to run after updating zonefile
  -p, --passwd string   path to .htpasswd file (default ".htpasswd")
  -z, --zone string     path to DNS zonefile (default "ddns.zone")
```

## Installation

Pre-built binaries for a variety of operating systems and architectures are available to download from [GitHub Releases](https://github.com/CHTJonas/ddnsd/releases). If you wish to compile from source then you will need a suitable [Go toolchain installed](https://golang.org/doc/install). After that just clone the project using Git and run Make! Cross-compilation is easy in Go so by default we build for all targets and place the resulting executables in `./bin`:

```bash
git clone https://github.com/CHTJonas/ddnsd.git
cd ddnsd
make clean && make all
```

## Copyright

ddnsd is licensed under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).

Copyright (c) 2019–2021 Charlie Jonas.
