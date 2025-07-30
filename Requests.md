curl -X PUT http://localhost:8080/update/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john",
    "surname": "smith",
    "email": "newjohn@example.com"
}'

curl -X GET http://localhost:8080/users/1


curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john",
    "surname": "smith",
    "email": "newjohn@example.com"
}'