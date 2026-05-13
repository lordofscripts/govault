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
* Sample `tree_sample_docs.json` file with **Document** hierarchy.
* Sample `tree_sample_pics.json` file with **Pictures** hierarchy.
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

### Installation

To install the `govault` executables in your system:

`go get github.com/lordofscripts/govault@latest`

On **Windows** you have the `govault.exe` executable files.

##### Usage

First copy either of [Documents](../tree_sample_docs.json) or [Pictures](../tree_sample_pics.json)
to a customized file `tree.json`. Each of those has a single top JSON node, but you can
merge them into one do maintain both Pictures and Documents with a single specification file.
Ensure you [validate your JSON](https://jsonlint.com/) `tree.json` file before using it with
`govault`.

> `$govault {COMMAND} [OPTIONS]`

Where `COMMAND` can be any of `-create`, `-update`, `-tags`, `-query` or `-help`.

Options depend on the command but are any of:

* `-json FILE_PATH` your JSON file path, defaults to `tree.json` in the current directory.
* `-overwrite` overwrite existing folder metadata on creation.
* `-perm "PERM_VALUE"` folder permissions, defaults to `0755` or `rwxr-xr-x`
* `-root DIRECTORY` root directory for folder creation, defaults to `.` (current)

To Create the entire directory hierarchy (all options above apply):

> `$govault -create`

To Update the directory hierarchy and folder metadata after you update your JSON file
(all options from above apply, *except* `-overwrite` which is TRUE):

> `$govault -update`

To List all the *tags* you have mentioned in your JSON file and show you a tag cloud:

> `$govault -tags`

To Query your JSON file to suggest you where to store or find a file given one or more
tags:

> `$govault -query "TAG(S)"`

For the `-query` you can specify a single tag or multiple comma-separated tags.

Read the manual, there is plenty of info here and the source code is documented in the 
[Pkg Info](https://pkg.go.dev/github.com/lordofscripts/govault).


***

Copyright &copy;2025-2026 Lord of Scripts
