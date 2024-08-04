# SCS - Simple Client-Server via RabbitMQ

You can run it with docker-compose and the default .env:
```
docker-compose up -d
```

Uses CSP (Common Sequantial Processes) pattern to distribute the work between logically independant routines.
After messages are received from AMQP, they are being sent to router that will use consistient hashing to distribute them 
between set of pre-configured workers so all operations for same ```key``` are serialized to execute in same worker to avoid potential race conditions.

The ```commands``` directory contains multiple files with commands. Which ones to be loaded by the client during test is specified in ```client/Dockerfile```. Each file will be loaded and executed concurrently to simulate multiple clients.

After client finishes, it won't exit so we can attach to the running container and inspect the saved files.