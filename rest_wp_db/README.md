
### Service launch

```bash
go run cmd/main.go
```

The server will start at http://localhost:8080

### Create a task

```bash
curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" \
     -d '{"type": "email", "payload": {"recipient": "user@example.com"}}'
```

### Get all tasks

```bash
curl http://localhost:8080/tasks
```

### Get task by ID

```bash
curl http://localhost:8080/tasks/{id}
```

### Cancel an active task

```bash
curl -X POST http://localhost:8080/tasks/{id}/cancel
```

### Delete a task record

```bash
curl -X DELETE http://localhost:8080/tasks/{id}
```