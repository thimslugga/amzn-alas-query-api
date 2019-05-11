# amzn-alas-query-api

A scraper and API for the Amazon Linux (2) Security Advisories Bulletin

## Description

This repository contains an API that can be used to query for known vulnerabilities
against a list of `Amazon Linux (2)` packages.
It's pretty limited in functionality right now, but I thought what I had so far was kinda cool so figured I'd post to Github.

The inspiration behind this API was wanting an alternative to Inspector. Instead of running agents on my EC2 instances, I prefer to keep a package manifest for each AMI I build, usually in DynamoDB or something. That way I can query the AMI for vulnerabilities and then tag EC2's for critical patching by association.

Usually this just means running an `rpm -qa > package_list.txt` at the end of the AMI build, downloading the file, and then having a script parse the results into a table.

You could accomplish this more or less the same with Inspector, but at an extra cost and with increased complexity in your AMI build pipeline.

## Building

I use `docker-compose` and a scratch container for the whole thing, but you can also just build a normal executable.

With `docker-compose` you can just:

```
$> cd docker && docker-compose build
```

Otherwise:

```
$> go get -d ./...
$> go build .
```

## Running

The API uses Redis to keep a persistent cache of the ALAS feeds. This way it doesn't need to re-scrape the entire list every launch. If you use `docker-compose` you can:

```
$> cd docker && docker-compose up
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

```js
{
  "packages": [
    "amazon-linux-extras-1.6.7-1.amzn2.noarch"
    // ...
  ]
}
```

The package strings are expected to be in RPM format. This is the format you get when you do `rpm -qa`.

Here is an example of passing the entire output of `rpm -qa` to the API. I use an `amazonlinux` docker container with `jq` installed:

```shell
bash-4.2> rpm -qa | jq --raw-input --slurp 'split("\n") | map(select(. != ""))' - -M | curl -X GET http://localhost:8080/vulns --data @-
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
    },
    "gpgme-1.3.2-5.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "grep-2.20-3.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "info-5.1-5.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "jq-1.5-1.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "keyutils-libs-1.5.8-3.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "krb5-libs-1.15.1-20.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libacl-2.2.51-14.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libassuan-2.1.0-3.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libattr-2.4.46-12.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libblkid-2.30.2-2.amzn2.0.4.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libcap-2.22-9.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libcom_err-1.42.9-12.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libcrypt-2.26-32.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libcurl-7.61.1-9.amzn2.0.1.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libdb-5.3.21-24.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libdb-utils-5.3.21-24.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libffi-3.0.13-18.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libgcc-7.3.1-5.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libgcrypt-1.5.3-14.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libgpg-error-1.12-3.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libidn2-2.0.4-1.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libmetalink-0.1.2-7.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libmount-2.30.2-2.amzn2.0.4.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libnghttp2-1.31.1-1.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libselinux-2.5-12.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libsepol-2.5-8.1.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libssh2-1.4.3-10.amzn2.1.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libstdc++-7.3.1-5.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libtasn1-4.10-1.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libunistring-0.9.3-9.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libuuid-2.30.2-2.amzn2.0.4.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libverto-0.2.5-4.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "libxml2-2.9.1-6.amzn2.3.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "lua-5.1.4-15.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "ncurses-6.0-8.20170212.amzn2.1.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "ncurses-base-6.0-8.20170212.amzn2.1.2.noarch": {
      "vulns": [],
      "errors": []
    },
    "ncurses-libs-6.0-8.20170212.amzn2.1.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nspr-4.19.0-1.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nss-3.36.0-7.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nss-pem-1.0.3-5.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nss-softokn-3.36.0-5.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nss-softokn-freebl-3.36.0-5.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nss-sysinit-3.36.0-7.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nss-tools-3.36.0-7.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "nss-util-3.36.0-1.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "oniguruma-5.9.6-1.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "openldap-2.4.44-15.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "openssl-libs-1.0.2k-16.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "p11-kit-0.23.5-3.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "p11-kit-trust-0.23.5-3.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "pcre-8.32-17.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "pinentry-0.8.1-17.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "popt-1.13-16.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "pth-2.0.7-23.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "pygpgme-0.3-9.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "pyliblzma-0.5.3-11.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "python-2.7.14-58.amzn2.0.4.x86_64": {
      "vulns": [],
      "errors": []
    },
    "python-iniparse-0.4-9.amzn2.noarch": {
      "vulns": [],
      "errors": []
    },
    "python-libs-2.7.14-58.amzn2.0.4.x86_64": {
      "vulns": [],
      "errors": []
    },
    "python-pycurl-7.19.0-19.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "python-urlgrabber-3.10-8.amzn2.noarch": {
      "vulns": [],
      "errors": []
    },
    "pyxattr-0.5.1-5.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "readline-6.2-10.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "rpm-4.11.3-25.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "rpm-build-libs-4.11.3-25.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "rpm-libs-4.11.3-25.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "rpm-python-4.11.3-25.amzn2.0.3.x86_64": {
      "vulns": [],
      "errors": []
    },
    "sed-4.2.2-5.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "setup-2.8.71-10.amzn2.noarch": {
      "vulns": [],
      "errors": []
    },
    "shared-mime-info-1.8-4.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "sqlite-3.7.17-8.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "system-release-2-7.amzn2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "tzdata-2018i-1.amzn2.noarch": {
      "vulns": [],
      "errors": []
    },
    "vim-minimal-7.4.160-4.amzn2.0.16.x86_64": {
      "vulns": [],
      "errors": []
    },
    "xz-libs-5.2.2-1.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "yum-3.4.3-158.amzn2.0.2.noarch": {
      "vulns": [],
      "errors": []
    },
    "yum-metadata-parser-1.1.4-10.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    },
    "yum-plugin-ovl-1.1.31-46.amzn2.0.1.noarch": {
      "vulns": [],
      "errors": []
    },
    "yum-plugin-priorities-1.1.31-46.amzn2.0.1.noarch": {
      "vulns": [],
      "errors": []
    },
    "zlib-1.2.7-17.amzn2.0.2.x86_64": {
      "vulns": [],
      "errors": []
    }
  }
}
```

Obviously this is a bad example as there are no known vulnerabilities in the docker container. However, if we had a list of packages with known vulnerabilities, we'd get better output.

```shell
bash-4.2> curl -X GET localhost:8080/vulns --data '["gnupg2-2.0.21-5.amzn2.0.3.x86_64"]'
{
  "results": {
    "gnupg2-2.0.21-5.amzn2.0.3.x86_64": {
      "vulns": [
        {
          "alas": "ALAS-2018-1045",
          "cves": [
            "CVE-2018-12020"
          ],
          "packages": [
            "gnupg2"
          ],
          "priority": "important",
          "newPackages": [
            {
              "name": "gnupg2",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.2",
              "arch": "i686",
              "raw": "gnupg2-2.0.22-5.amzn2.0.2.i686"
            },
            {
              "name": "gnupg2-smime",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.2",
              "arch": "i686",
              "raw": "gnupg2-smime-2.0.22-5.amzn2.0.2.i686"
            },
            {
              "name": "gnupg2-debuginfo",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.2",
              "arch": "i686",
              "raw": "gnupg2-debuginfo-2.0.22-5.amzn2.0.2.i686"
            },
            {
              "name": "gnupg2",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.2",
              "arch": "src",
              "raw": "gnupg2-2.0.22-5.amzn2.0.2.src"
            },
            {
              "name": "gnupg2",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.2",
              "arch": "x86_64",
              "raw": "gnupg2-2.0.22-5.amzn2.0.2.x86_64"
            },
            {
              "name": "gnupg2-smime",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.2",
              "arch": "x86_64",
              "raw": "gnupg2-smime-2.0.22-5.amzn2.0.2.x86_64"
            },
            {
              "name": "gnupg2-debuginfo",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.2",
              "arch": "x86_64",
              "raw": "gnupg2-debuginfo-2.0.22-5.amzn2.0.2.x86_64"
            }
          ],
          "link": "https://alas.aws.amazon.com/AL2/ALAS-2018-1045.html",
          "pubDate": "Thu, 09 Aug 2018 23:13:00 GMT"
        },
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
              "arch": "aarch64",
              "raw": "gnupg2-2.0.22-5.amzn2.0.3.aarch64"
            },
            {
              "name": "gnupg2-smime",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "aarch64",
              "raw": "gnupg2-smime-2.0.22-5.amzn2.0.3.aarch64"
            },
            {
              "name": "gnupg2-debuginfo",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "aarch64",
              "raw": "gnupg2-debuginfo-2.0.22-5.amzn2.0.3.aarch64"
            },
            {
              "name": "gnupg2",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "i686",
              "raw": "gnupg2-2.0.22-5.amzn2.0.3.i686"
            },
            {
              "name": "gnupg2-smime",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "i686",
              "raw": "gnupg2-smime-2.0.22-5.amzn2.0.3.i686"
            },
            {
              "name": "gnupg2-debuginfo",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "i686",
              "raw": "gnupg2-debuginfo-2.0.22-5.amzn2.0.3.i686"
            },
            {
              "name": "gnupg2",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "src",
              "raw": "gnupg2-2.0.22-5.amzn2.0.3.src"
            },
            {
              "name": "gnupg2",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "x86_64",
              "raw": "gnupg2-2.0.22-5.amzn2.0.3.x86_64"
            },
            {
              "name": "gnupg2-smime",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "x86_64",
              "raw": "gnupg2-smime-2.0.22-5.amzn2.0.3.x86_64"
            },
            {
              "name": "gnupg2-debuginfo",
              "epoch": "0",
              "version": "2.0.22",
              "release": "5.amzn2.0.3",
              "arch": "x86_64",
              "raw": "gnupg2-debuginfo-2.0.22-5.amzn2.0.3.x86_64"
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
