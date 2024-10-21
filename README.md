# Event Service

Core feature is in **server** folder and **proto** folder, **client** folder is only an adaptor if you want to call this service in another service and not up-to-date, you need to implement some func. like getAllEventsBy__

To run this service
### 1. Create a Shared Docker Network

Both RabbitMQ and the Event Service need to communicate over a shared Docker network. Before starting the containers, create this network:
```bash
docker network create shared-network
```

### 2. Set up RabbitMQ Server

The Event Service relies on a RabbitMQ instance for message brokering. You will need to create your own `docker-compose.yml` file for RabbitMQ. Below is an example configuration:
```yaml
version: "3.8"

networks:
  shared-network:
    external: true

services:
  rabbitmq:
    image: rabbitmq:3-management
    container_name: shared-rabbitmq
    ports:
      - "5672:5672" # AMQP port for producer/consumer
      - "15672:15672" # Management UI
    networks:
      - shared-network
```

Next, run the RabbitMQ container:
   ```bash
   docker-compose -f rabbitmq.yml up -d
  ```

You can verify that RabbitMQ is running by accessing the management UI in your browser at http://localhost:15672 (login with guest:guest).
### Step 3: Run Event Service
1. Ensure that RabbitMQ is running and is accessible on the `shared-network`.

2. Run the Event Service using Docker Compose:
   ```bash
   docker-compose -f docker-compose.yaml up -d
   ```
