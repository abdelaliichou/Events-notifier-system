# **Project Documentation: UCA Events Notifier System**

## **Table of Contents**  
1. [Getting Started](#getting-started)  
2. [Architecture Overview](#architecture-overview)  
3. [Services](#services)  
   - [Config Service](#config-service)  
   - [Scheduler Service](#scheduler-service)  
   - [Timetable Service](#timetable-service)  
   - [Alerter Service](#alerter-service)  
4. [Message Queue (NATS JetStream)](#message-queue-nats-jetstream)  
5. [Running the Project](#running-the-project)  
6. [Screenshots exemples](#Screenshots-exemples)  
7. [License & Contributors](#license--contributors)  

---

## **Getting Started**  

### **1. Install Docker**  
Ensure Docker is installed and running:  
```sh
sudo systemctl start docker
sudo systemctl status docker
```

### **2. Run the Project**  
To build and run all services in **one command**:  
```sh
sudo docker compose up --build
```
To **stop** the project:  
```sh
sudo docker-compose down -v
```

---

## **Architecture Overview**  
The system consists of **four main services**, which communicate using **NATS JetStream**:  

- **Config Service** → Manages alerts & resources  
- **Scheduler Service** → Fetches resources, retrieves events, and publishes them to NATS  
- **Timetable Service** → Processes and stores events, publishes updates to NATS  
- **Alerter Service** → Listens for changes and sends email notifications  

Each service is **Dockerized** and interacts via REST APIs and message queues.

---

## **Services**  

### **1. Config Service**  
#### **Overview**  
The **Config Service** manages:  
- Alerts (`/alerts`)  
- Resources (`/resources`)  

#### **Endpoints**  

##### **Alerts API**  
- `POST /alerts/` → Create an alert  
- `GET /alerts/` → Get all alerts  
- `GET /alerts/{id}/` → Get an alert by ID  

##### **Resources API**  
- `POST /resources/` → Create a resource  
- `GET /resources/` → Get all resources  
- `GET /resources/{id}/` → Get a resource by ID  
- `PUT /resources/{id}/` → Update a resource  
- `DELETE /resources/{id}/` → Delete a resource  

##### **Swagger Documentation**  
- **Local**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)  
- **Docker**: [http://config:8080/swagger/index.html](http://config:8080/swagger/index.html)  

---

### **2. Scheduler Service**  
#### **Overview**  
The **Scheduler Service**:  
- Fetches **resources** from Config Service  
- Retrieves **events** from UCA  
- Publishes events to **NATS JetStream**  

#### **How It Works**  
1. **Fetches resources** from Config Service  
2. **Queries UCA** for event data  
3. **Processes events** (groups by UID, merges resource IDs)  
4. **Publishes to NATS** (`EVENTS.stream`)  

#### **NATS Communication**  
- **Publishes events** → `EVENTS.stream`  

#### **Running the Service**  
```sh
go run main.go
```

---

### **3. Timetable Service**  
#### **Overview**  
The **Timetable Service**:  
- Receives event data from the **Scheduler**  
- Stores new events & updates existing ones  
- Publishes changes to **Alerter Service**  
- Exposes API endpoints for event retrieval  

#### **API Endpoints**  

##### **Retrieve All Events**  
```http
GET /events/
```

##### **Search Event by UID**  
```http
GET /events/search?uid={event_uid}
```

##### **Swagger Docs**  
- **Local:** [http://localhost:8090/swagger/index.html](http://localhost:8090/swagger/index.html)  
- **Docker:** [http://timetable:8090/swagger/index.html](http://timetable:8090/swagger/index.html)  

#### **NATS Communication**  
- **Subscribes to:** `EVENTS.stream` (receives events)  
- **Publishes to:** `ALERTS.stream` (sends changes)  

#### **Running the Service**  
```sh
go run main.go
```

---

### **4. Alerter Service**  
#### **Overview**  
The **Alerter Service**:  
- Listens for event changes (`ALERTS.stream`)  
- Retrieves events from **Timetable Service**  
- Fetches alert subscriptions from **Config Service**  
- Sends email notifications  

#### **How It Works**  
1. **Receives event UIDs** from `ALERTS.stream`  
2. **Fetches event details** from **Timetable Service**  
3. **Gets alert subscriptions** from **Config Service**  
4. **Sends emails** to subscribed users  

#### **Running the Service**  
```sh
go run main.go
```

---

## **Message Queue (NATS JetStream)**  
### **Streams**  
| Stream Name   | Used By  | Description |
|--------------|---------|-------------|
| `EVENTS.stream` | Scheduler → Timetable | Sends event data |
| `ALERTS.stream` | Timetable → Alerter | Sends event changes |

### **NATS Server**  
Running at: `nats://nats-server:4222`  
 
 
---


## **Running the Project**  

#### **Start Everything**  
```sh
docker-compose up --build
```


---



## **Screenshots Exemples**  
![Running image](https://github.com/abdelaliichou/Events-notifier-system/blob/main/screenshots/running.png)
![Started images](https://github.com/abdelaliichou/Events-notifier-system/blob/main/screenshots/started.png)
![Stopped images](https://github.com/abdelaliichou/Events-notifier-system/blob/main/screenshots/stopped.png)




---



## **License & Contributors**  

### **Contributors**  
- **ICHOU Abdelali** - [GitHub Profile](https://github.com/abdelaliichou)  

### **License**  
This project is licensed under the **MIT License**.


