# IPinfo Query Language (IQL) Specification

The IPinfo Query Language (IQL) is used to query, transform, and output data
from publicly available IPinfo APIs.

## Quickstart

To get a quick idea of what IQL is able to do, consider this query:

```
@ip=8.8.8.0/24 anycast=false country="US" ip>=8.8.8.253 @sort(desc)=ip @out(csv)=ip,city
```

Which will output the following CSV to `stdout`:

```csv
ip,city
8.8.8.255,Mountain View
8.8.8.254,Mountain View
8.8.8.253,Mountain View
```

## Data Source

TODO

## Post-Processing

TODO

## Key-Value Filter

Key-value pairs of the form `<key><op><value>` are used as filters.

Nested keys can be retrieved by joining keys in the path to the target key with
a dot, e.g. `<key1>.<key2>.<key3><op><value>`.

The supported operators are:

`=`: Checks for data with equal values.

`!=`: Checks for data with different values.

`>`: Checks for data with greater values.

`>=`: Checks for data with greater or equal values.

`<`: Checks for data with lesser values.

`<=`: Checks for data with lesser or equal values.

Values must be wrapped in quotation marks if they contain spaces, which would
otherwise indicate the start of a new token.

### Boolean Operators

Key-value filters can be combined using boolean operators.

The supported boolean operators are:

`AND`: Checks that the predicate of two filters are both met.

`OR`: Checks that at least one of the predicate of two filters are met.

`NOT`: Checks that the following predicate of a filter isn't met.

When multiple key-value filters appear without any boolean operator, `AND` is
implicitly used.

Using `NOT` on an expression is equivalent to negating all operators within the
expression.

### Operator Negation

Each of these operators has a "negation operator":

`=` negates `!=` and vice versa.

`>` negates `<=` and vice versa.

`<` negates `>=` and vice versa.

`AND` negates `OR` and vice versa.

`NOT` negates itself; `NOT NOT` is a no-op.

### Examples

#### Example 1

Filter for data with IPs `8.8.8.{0..254}`.

```
ip>=8.8.8.0 ip<8.8.8.255
```

#### Example 2

Equivalent to example 1 using `NOT` and `OR`.

```
NOT (ip<8.8.8.0 OR ip>8.8.8.254)
```

#### Example 3

Filter for data that isn't categorized as "anycast".

```
anycast=false
```

#### Example 4

INCORRECT: the `anycast` key will never have a string value.

```
anycast="false"
```

#### Example 5

Filter for data whose country key has the value `"US"` or `"PK"`.

```
country=US OR country=PK
country="US" OR country="PK"
```

#### Example 6

Filter for data whose ASN's domain is `"google.com"`.

```
asn.domain=google.com
asn.domain="google.com"
```

#### Example 7

Filter for data which is considered to be coming from a VPN.

```
privacy.vpn!=false
```

#### Example 8

Filter for data which has less than 1000 associated domains.

```
domains.total<1000
```

## Full IQL Examples

### Example 1

```
@ip=8.8.8.0/24 anycast=false country="US" ip>=8.8.8.253 @sort(desc)=ip @out(csv)=ip,city
```

```csv
ip,city
8.8.8.255,Mountain View
8.8.8.254,Mountain View
8.8.8.253,Mountain View
```

### Example 2

```
@ip=8.8.8.0/24 anycast=false country="US" ip<=8.8.8.2 @sort(asc)=ip @out(json)=ip,city
```

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
