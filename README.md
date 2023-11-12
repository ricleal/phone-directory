# phone-directory
A phone directory build in Go using GORM and GIN with a DDD approach


## Create a new contact
```bash
http POST localhost:8080/v1/users <<< '{
    "name": "John Doe",
    "phones": ["123-123-1234", "123-123-1235"],
    "addresses": ["123 Main St", "456 Main St"]
}'

http POST localhost:8080/v1/users <<< '{
    "name": "Jane Doe",
    "phones": ["123-123-1236", "123-123-1237"],
    "addresses": ["123 Main St", "456 Main St"]
}'
```

## Get a contact
```bash
http localhost:8080/v1/users/1
```