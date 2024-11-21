# author
**author** is command line tool for writing books and papers.
Allow users to convert documents from markdown into PDF and HTML. Shipped with predefined template for various type of documents.

> Important! Still under development.

## Installation and requirements

### Requirements

Require pandoc to be installed on system in order to be able to use it.

### Installation

```bash
go install github.com/zivlakmilos/author
```

## Manual

### Create project

```bash
author create
```

### Compile project

```bash
author build
```

### Build on file changes

```bash
author watch
```

### Display help

```bash
author --help
```

## Development

### Changelog

#### 1.0.0

- create
- build
- watch

### ToDo

- [x] Project generator system
- [x] Project compiler
- [x] Basic template
- [x] CLI Help
- [x] Add completion script for CLI
- [x] Watch
- [x] README
