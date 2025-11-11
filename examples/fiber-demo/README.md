# Fiber Demo Application

Ứng dụng demo cho GoErrorKit với Fiber framework.

## Chạy Demo

```bash
cd examples/fiber-demo
go run main.go
```

Server sẽ chạy tại `http://localhost:8081`

## Test Endpoints

### Panic Demos (Tự động recovered)

```bash
# Division by zero panic
curl http://localhost:8081/panic/division

# Index out of range panic
curl http://localhost:8081/panic/index

# Deep call stack panic
curl http://localhost:8081/panic/stack
```

### Custom Error Demos

```bash
# Business error (404)
curl http://localhost:8081/error/business?product_id=123

# System error (500)
curl http://localhost:8081/error/system

# Validation error (400)
curl http://localhost:8081/error/validation?age=15

# Auth error (401)
curl http://localhost:8081/error/auth

# External error (502)
curl http://localhost:8081/error/external?service=payment
```

## Kiểm tra Logs

Sau khi test, check file `logs/errors.log` để xem detailed error logs với:
- Chính xác panic location
- Full stack trace
- Call chain
- Request context

```bash
cat logs/errors.log | jq
```

