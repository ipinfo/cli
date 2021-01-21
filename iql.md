### Key-Value Filter

Key-value pairs of the form `<key><op><value>` are used as filters.

Nested keys can be retrieved by joining keys in the path to the target key with
a dot, e.g. `<key1>.<key2>.<key3><op><value>`.

The supported operators are:

`=`
    checks for data with equal values.
`!=`
    checks for data with different values.
`>`
    checks for data with greater values.
`>=`
    checks for data with greater or equal values.
`<`
    checks for data with lesser values.
`<=`
    checks for data with lesser or equal values.

Values must be wrapped in quotation marks if they contain spaces, which would
otherwise indicate the start of a new token.

#### Examples

`ip<8.8.8.255`
    filter for data with IPs less than `8.8.8.255`, e.g. `8.8.8.{0..254}`.

`anycast=false`
    filter for data that isn't categorized as "anycast".

`anycast="false"`
    INCORRECT: the `anycast` key will never have a string value.

`country=US`
`country="US"` (redundant quotation marks)
    filter for data whose country key has the value `"US"`.

`asn.domain=google.com`
`asn.domain="google.com"` (redundant quotation marks)
    filter for data whose ASN's domain is `"google.com"`.

`privacy.vpn!=false`
    filter for data which is considered to be coming from a VPN.

`domains.total<1000`
    filter for data which has less than 1000 associated domains.

### Full IQL Examples

#### Example 1

**IN**
```
@ip=8.8.8.0/24 anycast=false country="US" ip>=8.8.8.253 @sort(desc)=ip @out(csv)=ip,city
```

**OUT**
```csv
ip,city
8.8.8.255,Mountain View
8.8.8.254,Mountain View
8.8.8.253,Mountain View
```

#### Example 2

**IN**
```
@ip=8.8.8.0/24 anycast=false country="US" ip<=8.8.8.2 @sort(asc)=ip @out(json)=ip,city
```

**OUT**
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
    }
]
```
