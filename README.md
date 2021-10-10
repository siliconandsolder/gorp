# Gorp
Almost a portmanteau of Go and Grep. 

## What is this?
This is a tiny grep-like application that searches files for strings that match a regular expression. The purpose of this project was to teach myself about Go's channels and goroutines. This is not a practical tool.

## Usage
There are two modes for this application: regular and verbose.

### Regular Mode
In regular mode, gorp with accumulate all matching strings, but only print them after it has finished searching the provided directory.

To search for matches in regular mode, execute gorp with the following command:
```
gorp [directory_to_search] [string_to_find] [file_extensions_to_open]
```
**Example**:
```
gorp ./dir/ "hello world!" .txt,.pdf
```
### Verbose Mode
In verbose mode, gorp prints matching strings as it finds them.

To search for matches in verbose mode, execute gorp with the following command:
```
gorp -v [directory_to_search] [string_to_find] [file_extensions_to_open]
```
**Example**:
```
gorp -v ./dir/ "hello world!" .txt,.pdf
```
