# go-servise-entity

Сервис в разработке и постоянно модифицируется, предназначен для изучения.  

Сервис для локального тестирования инструментов, в сервисе реализован http и gRPC клиент  

Документация на /api/v1/docs  

Тестовый эндпоинт c gRPC  
```bash
# Получите JWT токен через auth сервис
# Затем вызовите тестовый эндпоинт:

curl -X GET http://localhost:8010/api/v1/entity/test_grpc \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
