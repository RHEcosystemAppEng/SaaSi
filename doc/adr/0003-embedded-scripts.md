# 3. Embedded scripts

Date: 2023-03-10

## Status

Proposed: 2023-03-10

## Context

The project includes executable scripts and other resources that are referenced in the Go modules, and are executed 
using relative path. This can't work when the module is imported as part of another module, as the script lookup would fail.

## Decision

Use the Goland [embed](https://pkg.go.dev/embed) directive to access the files from the program and copy all the files
in a temporary folder at startup time (so, if they refer to each other using relative paths it would still work).

This works seamlessly in both CLI and REST environments.

## Consequences

We can integrate the file resources in the Go application and continue to refer to them using the embed directive.
