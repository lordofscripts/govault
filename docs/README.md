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

First copy either of [Documents](../cmd/json/tree_sample_docs.json) or [Pictures](../cmd/json/tree_sample_pics.json)
to a customized file `tree.json`. Each of those has a single top JSON node, but you can
merge them into one do maintain both Pictures and Documents with a single specification file.
Ensure you [validate your JSON](https://jsonlint.com/) `tree.json` file before using it with
`govault`.

> `$govault {COMMAND} [OPTIONS]`

Where `COMMAND` can be any of the following:

* `-help` to get usage help
* `-create` create a directory hiearchy of Pictures or Documents
* `-update` update directory hiearchy and folder metadata after you edited your reference JSON folder file.
* `-query` to query your tag cloud and suggest you where to deposit the file.
* `-list tokens|rules` (Filename template feature) list the *tokens* or *rules* defined in your JSON filenames templates files.

Options depend on the command but are any of:

* `-json FILE_PATH` your JSON file path. For using built-in configuration (which you can customize) you can use `FILE_PATH` set as `app:pics` or `app:docs` to point to the config files you edited.
* `-overwrite` overwrite existing folder metadata on creation. *Does not apply to* `-update`
* `-perm "PERM_VALUE"` folder permissions, defaults to `0755` or `rwxr-xr-x`
* `-root DIRECTORY` root directory for folder creation, defaults to `.` (current)

Read the manual, there is plenty of info here and the source code is documented in the 
[Pkg Info](https://pkg.go.dev/github.com/lordofscripts/govault).

#### Folder Hierarchy Commands

These commands or options are related to the upkeeping of the Documents, Pictures or
whatever directory hiearchy.

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

#### Filename Template Commands

Usually the pictures from your devices have a specific filename that follows a
pattern like `IMG_20261025_230015.jpg`. Or perhaps you want (recommended!) to
organize your *Digital Legacy* with well-formed names that identify owners,
institutions, categories or whatever.

For that, as of `goVault v1.1.0` you can specify filename templates in your
*directory hiearchy* JSON files. The `govault_pics.json` and `govault_docs.json`
files in your configuration directory are a good start, feel free to customize
them to your own hierarchy.

The templates you *optionally* specify in the two config files just mentioned
could look like `IMG_%D_%T*.%X` (example for images) so that you can validate 
compliance of your filenames with the suggested template. This highly flexible 
system uses **Tokens** that you define in the `govault_templates.json` 
configuration file. There you define Predefined tokens as a list of  *key-value* 
pairs (valid name and description), a list of values (`IMG|PIC|DCIM`) or
even *Regular Expression* Rules for more complex variations. You can examine
a [sample template](../cmd/json/empty_template.json). I have found this
very valuable for both my Document and Picture vaults where my family (or survivors)
can quickly identify the relevance of files. Obviously, very helpful for me
as well.

Get a list of Tokens in your `govault_templates.json` file.

> `$govault -list tokens`

Get a list of Rules in your `govault_templates.json` file. These are for more
complex filename patterns.

> `$govault -list rules`

Use the templates from the local `.foldermeta.json` referencing the
tokens and rules from your `govault_templates.json` file to verify compliance
of the filenames in the *current directory*.

> `$govault -check`

## QVault

This is an auxillary executable that allows the user to do free-form
SQL-like queries on the command line. It uses your logical folder structure
which is contained in one or more configuration files. This logical structure
is reflected on the physical file system.

![QVault](./assets/qvault_query_sample.png)

At this moment it has full `SELECT` support and partial `UPDATE` support. The
update is not completed because it needs to be persisted.

Here are some examples of FQL (Folder Query Language) statements:

* `SELECT Name,Path FROM Pictures WHERE Encrypted = true`
* `SELECT Id,Name,Path FROM Documents WHERE Path LIKE 'Documents/Legal/%'`
* `SELECT Name,Tags FROM Pictures WHERE * LIKE '%Project%'`
* `SELECT * FROM Documents WHERE Name LIKE '%Munich%'`
* `UPDATE folders SET Name = 'Archived' WHERE Id = '1.2'`
* `UPDATE folders SET Encrypted = false, Tags = 'public' WHERE Name = 'Public Docs'`

In the `SELECT` clause you can use any of the `Folder` field names: `Id, Name, Path, Tags, Template` or the `*` character to select *all* fields.

The `FROM` clause indicates the name of the *top root node* as specified in your
source JSON file that contains the File Directory Structure that is replicated on
the physical file system. You can have one file for Documents and another for 
Pictures or both together, so there is no limitation of the number of top root nodes.
The logical structure is hiearchical, as in a tree.

The `WHERE` clause has built-in operators like `=`, `CONTAINS` and `LIKE`. The 1st
argument in the `WHERE` clause is the field to examine, or you can use `*` to specify
that it look at `Name, Path, Id, Tags` for the match. 

The `LIKE` operator works like in SQL, so you can use `_` to signify a single
character or `%` to signify several characters.

***

Copyright &copy;2025-2026 Lord of Scripts&tm;
