
# create user

```bash
curl -X POST -H "Content-Type application/json" -d '{"id": "1", "name": "lisi"}' http://localhost:8080/users
```

# get user

```bash
curl -i 35.232.156.79:8080/users
curl -i 127.0.0.1:8080/users
```