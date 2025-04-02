sudo systemctl start docker
sudo systemctl status docker
sudo docker compose up --build 
sudo docker-compose down -v


# Config Service

## Overview  
The **Config Service** provides APIs for managing alerts and resources. It includes:  
- **Alert Management**: Create, retrieve, and get alerts by ID.  
- **Resource Management**: Create, retrieve, update, and delete resources.  
- **Swagger Documentation**: Available for testing API endpoints.  

---

## Endpoints  

### 1. Alerts API  
- `POST /alerts/` → Create an alert  
- `GET /alerts/` → Get all alerts  
- `GET /alerts/{id}/` → Get an alert by ID  

### 2. Resources API  
- `POST /resources/` → Create a resource  
- `GET /resources/` → Get all resources  
- `GET /resources/{id}/` → Get a resource by ID  
- `PUT /resources/{id}/` → Update a resource  
- `DELETE /resources/{id}/` → Delete a resource  

### 3. Swagger Docs  
- **Standalone**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)  
- **Dockerized**: [http://config:8080/swagger/index.html](http://config:8080/swagger/index.html)  

---

## Running the Service  

### Standalone  
```sh
go run cmd/main.go



# Scheduler Service

## Overview  
The **Scheduler Service** is responsible for:  
- Fetching resources from the **Config Service**  
- Retrieving events from the **UCA server**  
- Publishing events to **NATS JetStream**, where they will be consumed by the **Timetable Service**  

---

## How It Works  
1. **Fetch Resources**: Scheduler calls the `Config` service to get the list of resources.  
2. **Fetch Events from UCA**: Using the retrieved resources, it queries the **UCA server** for event data.  
3. **Process Events**: Groups events by UID and merges resource IDs.  
4. **Publish to NATS**: Sends processed events to **NATS JetStream** under the `EVENTS` stream.  

---

## Endpoints & Functionalities  

### 1. Fetching Resources  
- Calls the **Config Service** at:  
  - `http://config:8080/resources/` (when running in Docker)  
  - `http://localhost:8080/resources/` (when running standalone)  
- Parses and prints all available resources.  

### 2. Fetching Events from UCA  
- Uses the retrieved resource IDs to fetch events from the UCA server.  
- Groups events by UID for better processing.  

### 3. Publishing to NATS JetStream  
- Sends events to `EVENTS.stream` in **NATS JetStream**.  
- Uses **NATS Server** running at `nats://nats-server:4222`.  

---

## Running the Service  

### Standalone  
```sh
go run main.go



# Timetable Service

## Overview
The **Timetable Service** is responsible for:
- Receiving event data from the **Scheduler Service** through **NATS JetStream**.
- Comparing incoming events with its database.
- Storing new events and updating modified events.
- Publishing changes to **NATS JetStream** for the **Alerter Service**.
- Exposing API endpoints for retrieving event information.

---

## How It Works

1. **Receiving Events**: The service subscribes to `EVENTS.stream` from the **Scheduler Service**.
2. **Processing Events**:
   - Checks if events exist in the database.
   - Adds new events to the database.
   - Updates existing events if changes are detected.
   - Sends changes to the **Alerter Service**.
3. **Publishing Updates**:
   - Changes are published to `ALERTS.stream` in **NATS JetStream**.
   - The **Alerter Service** then notifies users of changes.

---

## API Documentation

### Base URL
- **Local:** `http://localhost:8090`
- **Docker:** `http://timetable:8090`

### Endpoints

#### 1. Retrieve All Events
```http
GET /events/
```
- **Description:** Fetch all events stored in the database.
- **Response:**
```json
[
  {
    "uid": "event123",
    "title": "Math Lecture",
    "resourceIDs": ["roomA", "teacherB"]
  }
]
```

#### 2. Search Event by UID
```http
GET /events/search?uid={event_uid}
```
- **Description:** Fetch a specific event using its unique identifier.
- **Response:**
```json
{
  "uid": "event123",
  "title": "Math Lecture",
  "resourceIDs": ["roomA", "teacherB"]
}
```

#### 3. Swagger Documentation
- **Local:** [Swagger UI](http://localhost:8090/swagger/index.html)
- **Docker:** [Swagger UI](http://timetable:8090/swagger/index.html)

---

## NATS Message Queue

### Consumer (Receiving Events from Scheduler)
- **Stream:** `EVENTS.stream`
- **Functionality:**
  - Listens for events from the **Scheduler Service**.
  - Parses event data and checks for changes.
  - Saves new events and updates modified events.
  - Publishes changes to the `ALERTS.stream`.

### Producer (Publishing to Alerter)
- **Stream:** `ALERTS.stream`
- **Functionality:**
  - Sends `UID` and **changes** of modified events to **Alerter Service**.

---

## Running the Service

### Standalone
```sh
go run main.go
```

### Dockerized

#### Dockerfile
```dockerfile
FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o timetable ./main.go
CMD ["./timetable"]
```

#### Docker Compose
```yaml
version: "3.8"
services:
  timetable:
    build: ./timetable
    depends_on:
      - nats
    networks:
      - microservices_network
```

#### Start Everything
```sh
docker-compose up --build
```

---

## Database Schema

### Events Table
| Column       | Type         | Description |
|-------------|-------------|-------------|
| `id`        | UUID        | Unique event ID |
| `uid`       | VARCHAR(255)| Unique event UID |
| `title`     | TEXT        | Event title |
| `resourceIDs` | JSON       | List of associated resources |

### Resources Table
| Column     | Type         | Description |
|-----------|-------------|-------------|
| `id`      | UUID        | Unique resource ID |
| `event_uid` | VARCHAR(255)| Associated event UID |
| `resource_name` | TEXT  | Name of the resource |

---

## Contributors
- **Your Name** - [GitHub Profile](https://github.com/yourgithub)

## License
This project is licensed under the **MIT License**.

---




package mail

import "middleware/example/internal/models"

import (
	"fmt"
	"log"
	"strings"
)

func PreparingMail(mail string, events []models.Event, eventChanges map[string]map[string]interface{}, all bool) {
	// Loop through each event and format the changes
	for _, event := range events {

		var mailBody strings.Builder

		mailBody.WriteString(fmt.Sprintf("Bonjour.\n\nÉvénement %s a été modifié.\n\n", event.Name))

		eventData, exists := eventChanges[event.UID]
		if !exists {
			continue
		}

		changes, hasChanges := eventData["changes"].(map[string]interface{})
		if !hasChanges {
			continue
		}

		mailBody.WriteString("Les changements apportés : \n\n")
		for field, change := range changes {
			changeMap, ok := change.(map[string]interface{})
			if !ok {
				continue
			}

			oldValue, oldExists := changeMap["old"].(string)
			newValue, newExists := changeMap["new"].(string)

			if oldExists && newExists {
				mailBody.WriteString(fmt.Sprintf("- %s: \n    - Avant: \"%s\"\n    - Après: \"%s\"\n\n", field, oldValue, newValue))
			}
		}

		// Add event details
		mailBody.WriteString("\nDétails : \n")
		mailBody.WriteString(fmt.Sprintf("- Nom : %s\n", event.Name))
		mailBody.WriteString(fmt.Sprintf("- Description : %s\n", event.Description))
		mailBody.WriteString(fmt.Sprintf("- Début : %s\n", event.Start.Format("2006-01-02 15:04:05")))
		mailBody.WriteString(fmt.Sprintf("- Fin : %s\n", event.End.Format("2006-01-02 15:04:05")))
		mailBody.WriteString(fmt.Sprintf("- Lieu : %s\n", event.Location))
		mailBody.WriteString(fmt.Sprintf("- Dernière mise à jour : %s\n\n", event.LastUpdate.Format("2006-01-02 15:04:05")))

		mailBody.WriteString("Cordialement,\nAbdelali ichou du L'équipe P&G Innovations\n")
		//mailBody.WriteString("------------------------------------------------------\n")

		// Send mail
		// fmt.Printf("📧 Sending mail to %s:\n%s\n", mail, mailBody.String())
		sendMail(mail, mailBody.String())
	}
}

func sendMail(mail string, content string) {
	// Token required for the API
	token := "PueiQkxDnrLjMHlFzfVVUCojDPTlZchQeRWecXTk"

	// Example event data
	event := struct {
		EventContent string
	}{
		EventContent: content,
	}

	// Get email html shape from template
	emailContent, err := models.GetEmailContent("mail.html", event)
	if err != nil {
		log.Fatalf("Failed to generate email content: %s", err)
	}

	// Send email
	// here im working with abdelali.ichou@etu.uca.fr instead of the mail argument just for testing purposes
	err = models.SendEmail("abdelali.ichou@etu.uca.fr", emailContent.Subject, emailContent.Body, token)
	if err != nil {
		log.Fatalf("Failed to send email: %s", err)
	}
}

