# Kafka Stock Exchange

## Summary

Kafka Stock Exchange is a full-stack project designed to provide hands-on experience with microservices using Apache Kafka to simulate real-world scenarios at scale. This project comprises five different services, each serving a specific role in the stock exchange ecosystem:

1. **Platform APIs** - This service is responsible for providing a user-friendly interface and utilizing backend APIs to enable user interaction with the stock exchange.
Running on port 8080

2. **UserAnalytics** - The UserAnalytics service is dedicated to analyzing user behavior and generating valuable insights from user interactions with the stock exchange platform. Running on port 8081

3. **Stock Ingestor** - This service focuses on ingesting real-time stock data, ensuring that the exchange has access to the latest market information to make informed decisions. We are going to mock the data to stimulate the real-time data. Running on port 8082

4. **Stock Aggregator** - The Stock Aggregator service plays a pivotal role in aggregating and consolidating stock data from Stock Ingestor, making it available for other parts of the system to consume. Running on port 8083

5. **Order Processor** - The Order Processor service handles the critical task of processing stock orders placed by users. It ensures the seamless execution of orders within the stock exchange. Running on port 8084
