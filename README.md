# Bareknews
This is a news management system

## Endpoint List

### Get a news

## Form error response with single field
```json
{
  "error": {
    "message" : "Something bad happened :(",
    "description" : "More details about the error here"
  }
}
```

## Form error response with multiple fields
```json
{
  "error": {
    "message" : "Validation Failed",
    "items" : [
      {
        "field" : "first_name",
        "message" : "First name cannot have fancy characters"
      },
      {
        "field" : "password",
        "message" : "Password cannot be blank"
      }
    ]
  }
}
```