# Small MJ sample

We want to write some code against `https://jsonplaceholder.typicode.com/todos/1`

Which returns 

```json
{
    "userId": 1,
    "id": 1,
    "title": "delectus aut autem",
    "completed": false
}
```

## What MJ brings

Check out `docker-compose.yaml` and `todo.yaml`
- Fake server for you to write integration tests against
- CDC that runs against the real service
- From one config; so your tests are always working against something representative of the real API.
