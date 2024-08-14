# Amazon Linux ALAS Query API

A scaper and local API soltuion for the Amazon Linux Security Advisories Bulletin (ALAS).

## Description

This repository contains an API that can be used to query for known vulnerabilities against a list of `Amazon Linux` packages.

The inspiration behind this API was wanting an alternative to Inspector. Instead of running agents on my EC2 instances, I prefer to keep a package manifest for each AMI I build, usually in DynamoDB or something. That way I can query the AMI for vulnerabilities and then tag EC2's for critical patching by association.

Usually this just means running an `rpm -qa > package_list.txt` at the end of the AMI build, downloading the file, and then having a script parse the results into a table.

You could accomplish this more or less the same with Inspector, but at an extra cost and with increased complexity in your AMI build pipeline.

## Building

To build with Docker Compose, you can just run the following:

```sh
cd docker && docker compose build
```

To build without docker:

```sh
go init mod github.com/thimslugga/amzn-alas-query-api
go get -d ./...
go build .
```

## Running

The API utilizes Redis to keep a persistent cache of the ALAS feeds in memory. This way it doesn't need to re-scrape the entire list every launch. 

You can deploy Redis and the API using Docker Compose:

```sh
cd docker && docker compose up -d
```

Which will start both the API and a local redis instance. If you have a different Redis instance you'd like to use you can configure it in the environment. Below is a table of all the configuration options.

|Environment Variable|Description|Default|
|:------------------:|:----------|:-----:|
|`LISTEN_ADDR`|The address/port to listen on|`:8080`|
|`CACHE_TTL`|How often, in seconds, to check for changes in the ALAS feeds|`300`|
|`REDIS_HOST`|The host and port to connect to Redis on|`redis:6379` (docker-compose instance)|
|`REDIS_PASSWORD`|If needed, a password for connecting to Redis|`None`|
|`REDIS_DATABASE`|The Redis database to use|`0`|

The first time the API runs against a given Redis database, it will need to parse the entire ALAS feed before it can serve requests. This can take upwards of 60-ish seconds.

## Usage

There is currently just a single endpoint available: `GET /vulns`.
This endpoint expects a `JSON` payload in the following format:

```json
// Just a list of RPM package strings
[
  "amazon-linux-extras-1.6.7-1.amzn2.noarch"
  // ...
]
```

The package strings are expected to be in RPM format. This is the format you get when you do `rpm -qa`.

Here is an example of passing the entire output of `rpm -qa` to the API. You can use an `amazonlinux` based container with `curl` and `jq` installed:

```sh
rpm -qa \
  | jq --raw-input --slurp 'split("\n") | map(select(. != ""))' - -M \
  | curl -X GET http://localhost:8080/vulns --data @-
```

Output:

```json
{
  "results": {
    "amazon-linux-extras-1.6.7-1.amzn2.noarch": {
      "vulns": [],
      "errors": []
    },
    "basesystem-10.0-7.amzn2.0.1.noarch": {
      "vulns": [],
      "errors": []
    },
    "bash-4.2.46-30.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "bzip2-libs-1.0.6-13.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "ca-certificates-2018.2.22-70.0.amzn2.noarch": {
      "vulns": [],
      "errors": []
    },
    "chkconfig-1.7.4-1.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "coreutils-8.22-21.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "cpio-2.11-27.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "curl-7.61.1-9.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "cyrus-sasl-lib-2.1.26-23.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "diffutils-3.3-4.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "elfutils-libelf-0.170-4.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "expat-2.1.0-10.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "file-libs-5.11-33.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "filesystem-3.2-25.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "findutils-4.5.11-5.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "gawk-4.0.2-4.amzn2.1.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "gdbm-1.13-6.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "glib2-2.54.2-2.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "glibc-2.26-32.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "glibc-common-2.26-32.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "glibc-langpack-en-2.26-32.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "glibc-minimal-langpack-2.26-32.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "gmp-6.0.0-15.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "gnupg2-2.0.22-5.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "gpg-pubkey-c87f5b1a-593863f8": {
      "vulns": [],
      "errors": [
        "Could not parse arch from package string: Splitting gpg-pubkey-c87f5b1a-593863f8 at . produced one result"
      ]
...
```

Obviously this is a bad example as there are no known vulnerabilities in the docker container. However, if we had a list of packages with known vulnerabilities, we'd get better output.

```shell
curl -X GET localhost:8080/vulns --data '["gnupg2-2.0.22-5.amzn2.0.2.x86_64"]'
```

Output:

```
{
  "results": {
    "gnupg2-2.0.22-5.amzn2.0.2.x86_64": {
      "vulns": [
        {
          "alas": "ALAS-2019-1203",
          "cves": [
            "CVE-2014-4617"
          ],
          "packages": [
            "gnupg2"
          ],
          "priority": "medium",
          "newPackages": [
            {
              "name": "gnupg2",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "x86_64",
              "raw": "gnupg2-2.0.22-5.amzn2.0.3.x86_64"
            }
          ],
          "link": "https://alas.aws.amazon.com/AL2/ALAS-2019-1203.html",
          "pubDate": "Fri, 03 May 2019 00:13:00 GMT"
        }
      ],
      "errors": []
    }
  }
}
```
