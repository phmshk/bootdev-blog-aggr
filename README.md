# Gator

Gator is a command-line tool for aggregating RSS feeds. It fetches posts from blogs and stores them in a database for local browsing.

## Prerequisites

Before running Gator, you must have the following software installed on your system:

- Go: Version 1.26 or higher is required to compile the source code.
- PostgreSQL: A running database instance is required to store users, feeds, and posts.

## Installation

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Compile the application using the Go toolchain:

   go build -o gator

## Configuration

Gator requires a configuration file located in your home directory named ".gatorconfig.json". This file tells the application how to connect to your database and which user is currently logged in.

Create the file at ~/.gatorconfig.json with the following structure:

{
"db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
"current_user_name": ""
}

Replace "username", "password", and "gator" with your actual PostgreSQL credentials and database name.

## Usage

Once configured, you can run the program using the gator command.

### User Management

- register <username>: Create a new user account.
- login <username>: Switch the active user.
- users: List all users in the system.
- reset: Completely clear all data from the database.

### Feed Management

- addfeed <name> <url>: Add a new RSS feed to the system.
- feeds: View all feeds registered by all users.
- follow <url>: Subscribe to a feed for the current user.
- unfollow <url>: Remove a subscription.
- following: List feeds the current user is following.

### Aggregation and Browsing

- agg <time_between_reqs>: Start the background process that fetches posts. Example: "agg 1m" for 1-minute intervals.
- browse [limit]: Display the most recent posts from your followed feeds.
