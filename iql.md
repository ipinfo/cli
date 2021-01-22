# IPinfo Query Language (IQL) Specification

The IPinfo Query Language (IQL) is used to query, transform, and output data
from publicly available IPinfo APIs.

## Quickstart

To get a quick idea of what IQL is able to do, consider this query:

```
@data=8.8.8.0/24 @sort=ip(desc) @out(csv)=ip,city
anycast=false country="US" ip>=8.8.8.253
```

Which will output the following CSV to `stdout`:

```csv
ip,city
8.8.8.255,Mountain View
8.8.8.254,Mountain View
8.8.8.253,Mountain View
```

## Data Source

Data sources specify which API to use and which data to gather against that
API.

Data sources are specified by:

```
@data=<values>
```

The supported data `<values>` are:

`<ip>`: Individual IP value, e.g. `@data=8.8.8.8`.

`<cidr>`: IP range using CIDR syntax, e.g. `@data=8.8.0.0/16`.

`<value>,<value>,...,<value>`: multiple values separated by a comma (`,`), e.g.
`@data=1.1.1.1,8.8.8.8,9.9.9.9,8.8.0.0/16`.

In the future, `<asn>` and other data sources will be supported as well.

## Post-Processing

Post-processing happens on data that is gathered and filtered.

### Sorting

Sorting data in ascending or descending order on multiple fields is possible
with the following syntax:

```
@sort=<field>(<?order>)
```

#### Sort Fields

The supported `<field>` values are all those available for the specified data
source.

Multiple fields may be specified by separating each field with a comma.

Nested fields are specified by combining all keys in the path to the target key
with a dot (`.`).

Fields may be specified as `$<n>` where `n` is an index starting from `1` that
refers to the position of a field in the `@out` pre-processing directive.

#### Sort Orders

The supported orders `<order>` are:

`asc`: ascending order; this is the default.

`desc`: descending order.

The default order, if not specified, is ascending order.

#### Example 1

Sort data in ascending order by IP.

```
@sort=ip
```

#### Example 2

Sort data in ascending order by city, and then by IP for conflicting cities.

```
@sort=city(asc),ip
```

#### Example 3

Sort data in descending order by IP.

```
@sort=ip(desc)
```

#### Example 4

Sort data in descending order by city, and then descending order by IP for
conflicting cities.

```
@sort=city(desc),ip(desc)
```

#### Example 5

Sort data in descending order by the 1st field of `@out`, and then descending
order by IP for conflicting values of the 1st field.

```
@sort=$1(desc),ip(desc)
```

### Output

IQL can output specific fields in different formats.

The following syntax is used to specify an output format:

```
@out(<?format>)=<fields>
```

#### Output Fields

The supported `<fields>` are all those available for the specified data source.

Multiple fields may be specified by separating each field with a comma.

Nested fields are specified by combining all keys in the path to the target key
with a dot (`.`).

#### Output Format

The supported formats are:

`csv`: CSV format; the header will consist of the specified fields.

`json`: JSON format; an array holding JSON objects, each of whose key/value
pairs map to a field and its value, respectively.

`fmt`: generic format allowing each row to be generically specified by the
user. See the "Generic Output Format" section.

The default format, if not specified, is `csv`.

##### Generic Output Format

When `fmt` is specified as the output format, the value is no longer
interpreted as a list of fields. Instead, it must be a quoted string containing
arbitrary text along with placeholders for field values, e.g.:

```
@out(fmt)="<ip> was geolocated to <city>."
```

The above example will output each row as the above text, with the fields
between `<` and `>` replaced by their corresponding value.

When using this format, each specified field, even if repeated, is counted as a
unique increasing number when using `$<n>` format in e.g. `@sort` directives.
So consider the following:

```
@sort=$3(desc),$1
@out(fmt)="<ip> | <city> | <country>"
```

This will do a descending sort on the country, and then an ascending sort on
the IP when countries conflict.

## Key-Value Filter

Key-value pairs of the form `<key><op><value>` are used as filters.

Nested keys can be specified by joining keys in the path to the target key with
a dot, e.g. `<key1>.<key2>.<key3><op><value>`.

The supported operators `<op>` are:

`=`: Checks for data with equal values.

`!=`: Checks for data with different values.

`>`: Checks for data with greater values.

`>=`: Checks for data with greater or equal values.

`<`: Checks for data with lesser values.

`<=`: Checks for data with lesser or equal values.

Values must be wrapped in quotation marks if they contain spaces, which would
otherwise indicate the start of a new token.

### Boolean Operators

Key-value filters can be combined using boolean operators as
`<kv_filter><bool_op><kv_filter>` for binary boolean operators and
`<bool_op><kv_filter>` for unary boolean operators.

The supported boolean operators `<bool_op>` are:

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

The following examples show full IQL examples that may appear in the real
world.

### Example 1

```
@data=8.8.8.0/24 @sort=ip(desc) @out(csv)=ip,city
anycast=false country="US" ip>=8.8.8.253
```

```csv
ip,city
8.8.8.255,Mountain View
8.8.8.254,Mountain View
8.8.8.253,Mountain View
```

### Example 2

```
@data=8.8.8.0/24 @sort=$1(asc) @out(json)=ip,city
anycast=false country="US" ip<=8.8.8.2
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
