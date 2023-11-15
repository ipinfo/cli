package lib

import "regexp"

var ipV4Regex *regexp.Regexp
var ipV6Regex *regexp.Regexp
var ipRegex *regexp.Regexp

var v4IpCidrRegex *regexp.Regexp
var v6IpCidrRegex *regexp.Regexp
var ipCidrRegex *regexp.Regexp

var v4IpRangeRegex *regexp.Regexp
var v6IpRangeRegex *regexp.Regexp
var ipRangeRegex *regexp.Regexp

var v4IpSubnetRegex *regexp.Regexp
var v6IpSubnetRegex *regexp.Regexp
var ipSubnetRegex *regexp.Regexp

var v4CidrRegex *regexp.Regexp
var v6CidrRegex *regexp.Regexp
var cidrRegex *regexp.Regexp

var v4RangeRegex *regexp.Regexp
var v6RangeRegex *regexp.Regexp
var rangeRegex *regexp.Regexp

var v4SubnetRegex *regexp.Regexp
var v6SubnetRegex *regexp.Regexp
var subnetRegex *regexp.Regexp

const v4octet = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
const IPv4RegexPattern = `(` + v4octet + `\.){3}` + v4octet
const IPv4RangeRegexPattern = IPv4RegexPattern + `[:space:]*[-,][:space:]*` + IPv4RegexPattern
const IPv4CIDRRegexPattern = IPv4RegexPattern + `\/([1-2]?[0-9]|3[0-2])`

const IPv6RegexPattern = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
const IPv6RangeRegexPattern = IPv6RegexPattern + `[:space:]*[-,][:space:]*` + IPv6RegexPattern
const IPv6CIDRRegexPattern = IPv6RegexPattern + `\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8])`
