# Adster Forecasting Service
This project provides a forecasting service for an ad server platform. It includes an API for uploading CSV files and processing ad requests. This project is deployed on free AWS EC2 VM.

### Table of Contents
* Requirements
* Setup
* API Endpoints
* Forecast API
* CSV Upload API
* Database Schema
* Normalization and Denormalization
* Scaling Considerations
* Requirements
* Go 1.16+
* PostgreSQL (Free version from Render.com)
* Redis (Free version from cloud.redis.io)
* Curl or Postman for testing the API


### Requirements
Go 1.16+
PostgreSQL (using Free version from Render.com)
Redis (using Free version from cloud.redis.io)
Curl or Postman for testing the API

### Setup
1. Clone the repository 
2. Install dependencies
3. Update Environment variables if needed
4. Run the Application

### API Endpoints
1. Forecast API
Please not all the filters are not implemented. Need more understanding of domain knowledge like target_type, target_id etc

``` curl -X POST http://localhost:8080/forecast \
-H "Content-Type: application/json" \
-d '{
  "geo_target": {
    "included": [
      {
        "target_type": "COUNTRY",
        "target_id": 2840,
        "name": "United States",
        "country_code": "US"
      }
    ]
  },
  "device_type": {
    "included": [1, 2]
  },
  "inventory_url": {
    "included": ["example.com", "sample.org"]
  }
}'
```


2. CSV Upload API

Used batching of 1000 records per batch. 

```
curl -X PUT -F 'file=@sample.csv' http://localhost:8080/upload-csv
```


### Database Schema
Normalized Schema:
The data is stored in a normalized form across three tables (Up for discussion):
need to Denormalize of based on the frequency and scale of the filtering query

DDL sql queries are in sql directory.


##### Scaling Considerations
At very large scale, while the normalized schema offers better data integrity and reduced redundancy, querying across multiple tables might become slow due to frequent joins. In such cases:

##### Denormalization: 
Storing pre-joined data in a single table may speed up queries at the cost of increased redundancy.
Sharding: Sharding the data across multiple database instances could help distribute the load.
Caching: Redis can be further leveraged to cache frequently accessed data and reduce database load.


## Data Processing Pipeline Architecture
The data processing pipeline is designed to handle CSV file uploads in a scalable and event-driven manner, ensuring that the data is processed efficiently and stored in the database for forecasting.

##### CSV Upload to S3:
When a CSV file containing ad request data is uploaded to Amazon S3, it serves as the initial entry point of the data pipeline. S3 provides a scalable storage solution for handling large volumes of file uploads.

##### Event-Driven Trigger (AWS Lambda):
Once a file is uploaded to S3, an S3 Event Notification triggers an AWS Lambda function. This function acts as the orchestrator, responsible for initiating the downstream data processing. It ensures that each upload event triggers the necessary actions to process the CSV data.

##### Change Data Capture (CDC) with Kafka:
The Lambda function pushes a message containing metadata about the CSV file (e.g., file path, timestamp) to a Kafka topic. Kafka acts as the event bus for distributing messages across various consumers. The use of CDC ensures that every CSV upload is captured and processed reliably, and Kafka provides a scalable mechanism for processing these events asynchronously.

##### Kafka Consumer (Data Processor):
A Kafka consumer listens to the topic where the CSV metadata is published. Upon receiving a message, the consumer downloads the CSV file from S3 and begins processing the data. This processing involves:

##### Parsing the CSV data.
Validating and transforming the ad request data.
Batch inserting the data into the PostgreSQL database.
Batch Insertion into PostgreSQL:
The parsed data is batch-inserted into a normalized PostgreSQL schema (users, ad_details, request_logs) to ensure efficient storage and query performance. Batch processing reduces the number of database transactions, improving performance, especially when dealing with large files.



## Forecasting Algorithm Explanation
The forecasting service uses historical ad request data to predict daily impressions and reach based on the specified targeting criteria (such as geographic location, device type, and inventory URL). Here's an overview of how the algorithm works:

Historical Data Aggregation:

The system stores historical ad request data in a normalized PostgreSQL database across three main tables: users, ad_details, and request_logs.
Each record in request_logs links a user, ad details (e.g., domain, ad position, size), and a timestamp. This structure enables efficient querying of past ad requests based on targeting criteria like geographic location, device type, and inventory URL.
Targeting Criteria:

When a forecast request is made, the targeting criteria specify constraints such as:
Geo-targeting: Filtering by country, region, or city.
Device type: Filtering by mobile, tablet, or desktop.
Inventory URL: Filtering by the domain or URL where the ad is displayed.
These criteria are used to filter the historical data to match the desired audience or ad placement characteristics.


###### Querying the Database:
The forecast algorithm queries the request_logs table, joined with the users and ad_details tables, to find historical records that match the specified targeting criteria.
For example, if the criteria include users from "United States" using "Mobile" devices, the query will count all the records where these conditions are met.
Predicting Daily Impressions:

The impressions are calculated by counting the total number of ad requests (from request_logs) that match the filtering conditions over the past X days (where X is the configurable time window for the forecast).
The result is divided by the number of days to get the average daily impressions.
Example:
```
SELECT COUNT(*) / X AS daily_impressions
FROM request_logs
JOIN users ON request_logs.user_id = users.id
JOIN ad_details ON request_logs.ad_id = ad_details.id
WHERE users.geo_country = 'US'
  AND users.device_type = 1
  AND ad_details.domain IN ('example.com', 'sample.org')
  AND request_logs.timestamp >= NOW() - INTERVAL 'X DAYS';
  ```


###### Predicting Reach:


Reach is the number of unique users who would see the ad. To calculate this, the algorithm counts the distinct users (user_id) that match the targeting criteria over the same historical window.
Like impressions, the reach is also averaged over the number of days to provide a daily reach prediction.
Example:
```
SELECT COUNT(DISTINCT request_logs.user_id) / X AS daily_reach
FROM request_logs
JOIN users ON request_logs.user_id = users.id
JOIN ad_details ON request_logs.ad_id = ad_details.id
WHERE users.geo_country = 'US'
  AND users.device_type = 1
  AND ad_details.domain IN ('example.com', 'sample.org')
  AND request_logs.timestamp >= NOW() - INTERVAL 'X DAYS';
```

###### Handling Multiple Targeting Criteria:
The algorithm allows for multiple criteria to be included simultaneously. For example, a forecast request might target "Mobile users in the US visiting example.com." The forecast query handles such filtering by adding multiple WHERE conditions to match the request.

###### Caching with Redis:
To improve performance for frequently requested forecasts, the results of the database queries can be cached in Redis. If a similar forecast request is made again within a short time frame, the cached result is returned, avoiding the need to query the database repeatedly.


## Test Cases
Skipping test case for now, check test folder 