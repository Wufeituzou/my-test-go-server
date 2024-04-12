
# create user

```bash
curl -X POST -H "Content-Type application/json" -d '{"id": "1", "name": "lisi"}' http://localhost:8080/users
```

# get user

```bash
curl -i localhost:8080/users
```