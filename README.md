# go-service-entity

Сервис в разработке и постоянно модифицируется, предназначен для изучения и тестов.  


В сервисе реализован http и gRPC клиент.  
Деплой реализован на тестовый сервер.    

Документация на /api/v1/docs  
Реализованы CRUD операции  

Тестовый эндпоинт c gRPC  
```bash
# Получите JWT токен через auth сервис
# Затем вызовите тестовый эндпоинт:

curl -X GET http://localhost:8010/api/v1/entity/test_grpc \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
