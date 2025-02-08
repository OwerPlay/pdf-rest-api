# pdf-rest-api
## Installaion
1. Create ``.env`` file in root dir example: 
```
APP_PORT=8080
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=pdf-rest-api-db
MYSQL_USER=myuser
MYSQL_PASSWORD=mypassword
```

2. Run ``docker-compose up --build``


## Database ER Diagram
```mermaid
erDiagram
    USER {
        int id PK
        string name
        string email
    }
    USER_FILES {
        string filename
        int user_id PK, FK
        int file_id PK, FK
    }
    FILES {
        int id PK
        string status
        longblob pdf_file
        longblob parsed_file
    }
    PARSE_QUEUE {
        int file_id PK, FK
        date file_uploaded
    }
    USER ||--o{ USER_FILES : has
    USER_FILES ||--o{ FILES : link_to
    PARSE_QUEUE ||--o| FILES : processes
```
