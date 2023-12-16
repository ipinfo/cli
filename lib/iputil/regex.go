package iputil

import "regexp"

var IpV4Regex *regexp.Regexp
var IpV6Regex *regexp.Regexp
var IpRegex *regexp.Regexp

var V4IpCidrRegex *regexp.Regexp
var V6IpCidrRegex *regexp.Regexp
var IpCidrRegex *regexp.Regexp

var V4IpRangeRegex *regexp.Regexp
var V6IpRangeRegex *regexp.Regexp
var IpRangeRegex *regexp.Regexp

var V4IpSubnetRegex *regexp.Regexp
var V6IpSubnetRegex *regexp.Regexp
var IpSubnetRegex *regexp.Regexp

var V4CidrRegex *regexp.Regexp
var V6CidrRegex *regexp.Regexp
var CidrRegex *regexp.Regexp

var V4RangeRegex *regexp.Regexp
var V6RangeRegex *regexp.Regexp
var RangeRegex *regexp.Regexp

var V4SubnetRegex *regexp.Regexp
var V6SubnetRegex *regexp.Regexp
var SubnetRegex *regexp.Regexp

const V4Octet = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
const IPv4RegexPattern = `(` + V4Octet + `\.){3}` + V4Octet
const IPv4RangeRegexPattern = IPv4RegexPattern + `[-,]` + IPv4RegexPattern
const IPv4CIDRRegexPattern = IPv4RegexPattern + `\/([1-2]?[0-9]|3[0-2])`

const IPv6RegexPattern = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
const IPv6RangeRegexPattern = IPv6RegexPattern + `[-,]` + IPv6RegexPattern
const IPv6CIDRRegexPattern = IPv6RegexPattern + `\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8])`
