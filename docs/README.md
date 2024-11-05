# Commitea

**Commitea** is a Go-based tool designed for managing and automating tasks related to Git repositories. Built to
streamline workflow and improve productivity, this tool is ideal for software engineers who frequently work with Git.

### Development

To run the program without compiling:

```bash
go run ./cmd/commitea [command]
```

### Installation

Clone this repository:

```bash
git clone https://github.com/BillyUdders/commitea.git
```

Navigate into the project directory:

```bash
cd commitea
```

Build the project:

```bash
go build -o commitea ./cmd/commitea
```

### Commands

#### Log

```bash
./commitea log 
```

Displays the most recent commits in compact and dense form

#### Commit

```bash
./commitea commit 
```

Takes you through a form to create a well formed commit message

#### Status

```bash
./commitea status
```

Like log but with dirty files and branch names