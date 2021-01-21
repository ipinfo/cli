### 

### Example 1

```
ip:8.8.8.0/24 anycast:false country:"US" @sort(desc):ip @out(csv):ip,city
```

Output all IPs & their cities sorted by IP in descending order from the
8.8.8.0/24 range found in the US and not identified to be anycast.

Example output:

```csv
ip,city
8.8.8.255,Mountain View
8.8.8.254,Mountain View
8.8.8.253,Mountain View
...
```

### Example 2

```
ip:8.8.8.0/24 anycast:false country:"US" @sort(asc):ip @out(json):ip,city
```

Output all IPs & their cities sorted by IP in ascending order from the
8.8.8.0/24 range found in the US and not identified to be anycast.

Example output:

```json
[
    {
        "ip": "8.8.8.0",
        "city": "Mountain View"
    },
    {
        "ip": "8.8.8.1",
        "city": "Mountain View"
    },
    {
        "ip": "8.8.8.2",
        "city": "Mountain View"
    },
    ...
]
```
