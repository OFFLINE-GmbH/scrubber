# Scrubber

[![Build Status](https://travis-ci.org/OFFLINE-GmbH/scrubber.svg?branch=master)](https://travis-ci.org/OFFLINE-GmbH/scrubber)

Scrubber provides an easy way to clean up old files in a directory.

It is especially useful on platforms where `logrotate` is not easily available or disk space is sparse.

## Configuration

You can specify directories to clean up in a `toml` configuration file. You can define one or more strategies used for
each directory.

```toml
title = "Log scrubber example config"

[[directory]]
name = "Apache Logs"
path = "/var/logs/apache"
exclude = ["zip"]

    [[directory.strategy]]
    type = "size"
    action = "delete"
    limit = "100M"

    [[directory.strategy]]
    type = "age"
    action = "delete"
    limit = "1y"

    [[directory.strategy]]
    type = "age"
    action = "zip"
    limit = "1d"

[[directory]]
name = "Backups"
path = "/var/backups/yourapp"
include = ["tar.gz"]

    [[directory.strategy]]
    type = "age"
    action = "delete"
    limit = "1y"
```

### Directory

The following options are available for each `directory`:

| Option  | Description                                                                                              |
|---------|----------------------------------------------------------------------------------------------------------|
| name    | A descriptive name for this directory.                                                                   |
| path    | The full path to the directory.                                                                          |
| include | (Optional) Define what files should be included. All files without a matching extension will be ignored. |
| exclude | (Optional) Define what files should be excluded. All files with matching extension will be ignored.      |

You can either specify a `include` or a `exclude` rule but never both.

### Strategy

The following options are available for each `strategy`:

| Option | Possible values    | Description                                                                                                                                                                                                      |
|--------|--------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| type   | `age` and `size`   | If the files should be selected by their `age` (last modified) or their `size`.                                                                                                                                  |
| action | `delete` and `zip` | If matching files should be deleted or zipped. The `zip` action will remove the original file. Make sure to also exclude `zip` files from this rule so created zip files won't be cleaned up on subsequent runs. |
| limit  | A file size or age | Define the max. age as `1y`, `1d`, `2h` or the file size as `1M`, `1GB`, `1000B`. Supported units for the age are `m`, `h`, `d`, `w`, `y`. Supported units for the size are `B`, `KB`, `MB`, `GB`, `TB`, `PB`.   |

## Run

You can run `scrubber` from the command line. The following options are available:


| Param    | Default              | Description                                                   |
|----------|----------------------|---------------------------------------------------------------|
| -config  | scrubber.config.toml | The path to your configuration file.                          |
| -pretend | false                | If specified, scrubber will log actions but not execute them. |

```bash
# Check your config and see what will be done
./scrubber -config scrubber.config.toml -pretend
# Execute the action
./scrubber -config scrubber.config.toml
```
