# automated-message-sender

## Overview

The Automated Message Sender API is designed to manage and send messages automatically. It features endpoints to start and stop message sending, retrieve sent messages, and send new messages. The API uses Fiber for the web framework, GORM for database interactions, and Redis for caching.

## Features

- **Start and Stop Message Sending:** Control the automated message sending process.
- **Retrieve Sent Messages:** Get all messages that have been sent.
- **Send a Message:** Send messages to specified recipients.

## Technologies

- **Fiber:** Fast and lightweight web framework for Go.
- **GORM:** ORM for Go, used for database interactions.
- **PostgreSQL:** Database used for storing messages.
- **Redis:** Caching layer for message IDs.
- **Cron:** Scheduling library for periodic tasks.
- **Swagger:** API documentation and testing.

## Prerequisites

- Docker
- Docker Compose

## Setup

1. **Clone the Repository**

   ```sh
   git clone https://github.com/iremkzlrmk/automated-message-sender.git
   cd automated-message-sender
   ```
2. **Build the App**

   ```sh
   docker-compose build
   ```
3. **Run the App**

    ```sh
   docker-compose up
   ```

## Documentation

Go to `http://localhost:8080/swagger/index.html` for the documentation and further details of the API implementation.
