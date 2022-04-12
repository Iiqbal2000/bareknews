# Bareknews
This is a news management system

## Endpoint List

### Add a Tag

```bash
curl -X "POST" "http://localhost:3333/tags" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
  "title": "I Am Ozzy",
  "author": "Ozzy Osbourne",
  "pages": 294,
  "quantity":10
}'
```

## Form Error Response
```json
{
  "error": {
    "message" : "Something bad happened",
  }
}
```