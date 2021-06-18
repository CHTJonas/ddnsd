# Dynamic DNS Daemon

This repo hosts the code for a small dynamic DNS daemon (ddnsd) that I use to setup `_acme-challenge` records when issuing [Let's Encrypt](https://letsencrypt.org/) records for some of my personal infrastructure. It's written in Go and designed to be as small and as simple as possible. Host machines poke ddnsd over HTTP by making a POST request to an API endpoint during a [dehydrated](https://github.com/dehydrated-io/dehydrated) hook step. HTTP Basic authentication is used with credentials stored in a .htpasswd file. The username dictates which record to update which makes sure that API credentials can't be shared.

## Usage

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

## Installation

TODO

## Copyright

ddnsd is licensed under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).

Copyright (c) 2019–2021 Charlie Jonas.
