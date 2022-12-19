# GO CLIENT SEARCH TOOL

## Initialize the client
``` go
    client := client.NewClient()
    defer client.CloseConnection()
```

## Exposed Interface Methods
``` go
	Search(slug, search string) ([]models.ResponseSearchIndex, error)
	SearchByField(slug string, fieldSearch map[string]interface{}) ([]models.ResponseSearchIndex, error)
	Index(slug string, uid string, table_info string, action map[string]interface{}, search_value map[string]interface{}) (string, error)
	Delete(slug, uid, table_info string) (string, error)
    CloseConnection()
```

### Example Index
``` go
id, err := client.Index("slug", "uid", "table_name", map[string]interface{}{"id": 1}, map[string]interface{}{"first_name": "Ram", "last_name": "Sharma"})
```

### Example Global Search
``` go
data, err := client.Search("slug_name", "search_text")
```

### Example Delete Record
``` go
data, err := client.Delete("slug_name", "uid", "table_name")
```

