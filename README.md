# Blogator

Blogator is a command-line RSS feed aggregator written in Go. It allows users to register accounts, follow RSS feeds, continuously fetch posts from those feeds, and browse saved posts from the terminal.

## Prerequisites

Before installing Blogator, make sure you have the following installed:

- [Go](https://go.dev/) — required to build and run the application
- [PostgreSQL](https://www.postgresql.org/) — required for storing users, feeds, follows, and posts

PostgreSQL must be running before using Blogator.

## Installation

Install Blogator using:

```bash
go install github.com/Ha0cH/blogator@latest
```

Make sure your Go binary directory is included in your system's `PATH` so that you can run `blogator` directly from the terminal.

You can check your Go path with:

```bash
go env GOPATH
```

Go-installed binaries are typically stored in:

```text
$GOPATH/bin
```

## Database Setup

Create a PostgreSQL database for Blogator:

```sql
CREATE DATABASE blogator;
```

Make sure your PostgreSQL server is running and that you know the connection URL for the database.

For example:

```text
postgres://username:password@localhost:5432/blogator?sslmode=disable
```

## Configuration

Blogator reads its configuration from:

```text
~/.gatorconfig.json
```

Create the file if it does not already exist.

Example configuration:

```json
{
  "db_url": "postgres://username:password@localhost:5432/blogator?sslmode=disable",
  "current_user_name": ""
}
```

Replace the database username and password with your own PostgreSQL credentials.

The `current_user_name` field stores the user who is currently logged in and will be updated by Blogator when you register or log in.

## Usage

The general command format is:

```bash
blogator <command> [arguments]
```

### Register a User

Create a new user:

```bash
blogator register <name>
```

Example:

```bash
blogator register hao
```

### Log In

Switch to an existing user:

```bash
blogator login <name>
```

Example:

```bash
blogator login hao
```

### Add an RSS Feed

Add a feed to Blogator:

```bash
blogator addfeed "<feed-name>" "<feed-url>"
```

Example:

```bash
blogator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
```

### View Available Feeds

```bash
blogator feeds
```

### Follow a Feed

```bash
blogator follow <feed-url>
```

Example:

```bash
blogator follow "https://blog.boot.dev/index.xml"
```

### View Followed Feeds

```bash
blogator following
```

### Unfollow a Feed

```bash
blogator unfollow <feed-url>
```

### Run the Aggregator

Start the RSS aggregation process:

```bash
blogator agg <time-between-requests>
```

For example:

```bash
blogator agg 1m
```

This continuously checks followed RSS feeds for new posts every minute and saves them to the database.

You can use other Go duration values such as:

```text
30s
5m
1h
```

Stop the aggregator with `Ctrl+C`.

### Browse Posts

View saved posts:

```bash
blogator browse
```

You can optionally specify how many posts to display:

```bash
blogator browse 10
```

## Example Workflow

A typical first-time setup might look like:

```bash
blogator register hao

blogator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"

blogator follow "https://blog.boot.dev/index.xml"

blogator agg 1m
```

Then, in another terminal:

```bash
blogator browse 10
```