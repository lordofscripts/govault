# GO Vault

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/govault)
[![Go Report Card](https://goreportcard.com/badge/github.com/lordofscripts/govault?style=flat-square)](https://goreportcard.com/report/github.com/lordofscripts/govault)
![Build](https://github.com/lordofscripts/govault/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/govault.svg)](https://pkg.go.dev/github.com/lordofscripts/govault)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/govault)](https://github.com/lordofscripts/govault/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This small project was born out of a necessity. When you work on a series of laptops, PCs and
smartphones it will become difficult to manage your documents and pictures. Recently, I started
a project of reorganizing and indexing a Terabyte disk of decades of photos and documents. Some
are shadowed in different devices. I decided to use a new Document & Picture directory 
hiearchy in the backup and ALL devices I use, then it becomes easier to find stuff.

The first version used a YAML file that defined the structure. Recently I decided to add
tag support so that if I had a new document or picture I could ask this utility where it
would suggest me to store it based on a series of tags.

Once you decide on a structure, you copy the sample JSON file and modify it to suit your
needs. After that simply use the `-create` switch and in a second everything is created!
No need to manually create a hiearchy and then do the same elsewhere (or copy a ZIP).

### Features:

* Works for `Documents` and `Pictures` folders, but you decide how to name them.
* Sample `tree_sample_docs.json` file with **Document** hierarchy. *available in your config folder*
* Sample `tree_sample_pics.json` file with **Pictures** hierarchy. *available in your config folder*
* You can merge both JSON trees (Docs & Pics) into a single file!
* Create entire structure in click
* List your Tag Cloud
* Update the metadata folder files after you update the source JSON file
* Query your structure with a list of tags, it will suggest you which folders to use.
* You can flag a folder as *encrypted* if you plan to store there sensitive or encrypted data.

|                                                                                       | Show your support                                                                                                 |
| ------------------------------------------------------------------------------------- |:-----------------------------------------------------------------------------------------------------------------:|
| [ ![AllMyLinks](./assets/allmylinks.png)](https://allmylinks.com/lordofscripts)       | visit <br> Lord of Scripts&trade; <br> on [AllMyLinks.com](https://allmylinks.com/lordofscripts)                  |
| [ ![Buy me a coffee](./assets/buymecoffee.jpg)](https://allmylinks.com/lordofscripts) | buy Lord of Scripts&trade; <br> a Capuccino on <br>[BuyMeACoffee.com](https://www.buymeacoffee.com/lostinwriting) |

#### License

Read the [LICENSE](../LICENSE) for more details.

#### Requirements

* GO Language v1.23 or higher

---

## Motivation

While I lack an up-to-date computer, I work with several computing devices:
my day-to-day laptop, an old backup laptop, my former PC hard disks, a
Raspberry Pi and my Android phone.

Organizing a terabyte of data with only directories and filenames isn't enough.
At a point I realized it was better to use the exact same directory tree structure
on all devices, no more guessing where things are.

But, at the same time I needed directory metadata such as **Tags** and 
**Filename Templates**. Some directories needed to have documents organized
with a pre-defined filename template, just as modern digital camera images
get named.

When in doubt I also wanted the ability to query my file system with questions
like:

* Which Document/Picture directories would you suggest me for `legal, certificate` documents?
* Which directories have been marked as having encrypted/sensitive files?

You get the idea. I loved `SQL` since the first time I started using databases. I
wanted to query my file "vault" in a similar way, but with a simplified syntax
that was intuitive.

## Features

* A *Reproduceable Directory Structure* for your Documens, Pictures and others
  that you can deploy on all your devices.
* The logical directory hieararchy is maintained in a JSON file, where it is
  annotated with *metadata* such as **Tags, Templates** and whether it is 
  supposed to be **Encrypted*.
* Tag & (filename) Template metadata could be inherited.  

## Installation

To install the `govault` and `qvault` executables in your system:

`go get github.com/lordofscripts/govault@latest`

On **Windows** you have the `govault.exe` and `qvault.exe` executable files.

### Usage

These are some command-line flags you should familiarize yourself with
in order to use `govault` and `qvault` safely. I repeat, I am NOT liable
for any data loss! Also familiarize yourself with the configuration files
that you are expected to customize to suit your Document/Image/Other Repository.

* `-json FILE_PATH` your JSON file path. For using built-in configuration (which you can customize) you can use `FILE_PATH` set as `app:pics` or `app:docs` to point to the config files you edited.
* `-overwrite` overwrite existing folder metadata on `create` subcommand. *Does not apply to* `update` subcommand.
* `-perm "PERM_VALUE"` folder permissions, defaults to `0755` or `rwxr-xr-x`
* `-root DIRECTORY` root directory for folder creation, defaults to `.` (current) which would
  typically be your HOME directory. The root nodes of your JSON files are relative to THIS path.

Read the manual, there is plenty of info here and the source code is documented in the 
[Pkg Info](https://pkg.go.dev/github.com/lordofscripts/govault).


***
Copyright &copy;2025-2026 Lord of Scripts™ 
