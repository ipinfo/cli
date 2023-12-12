package main

var DetailedHelp = `ipinfo - IP information lookup tool.

Usage:
    ipinfo <command> [options] [arguments]

Commands:
    ipinfo <ip>       Lookup details for a specific IP address.
        Options:
            --field <field>, -f <field>    
               Lookup only specific fields in the output.
            --nocolor           
               Disable colored output.
            --json, -j          
               Output JSON format.
            --csv, -c           
               Output CSV format.
            --yaml, -y          
               Output YAML format.

        Examples:
            ipinfo 8.8.8.8
            ipinfo 8.8.8.8 --field hostname,city


    ipinfo <asn>      Lookup details for a specific ASN.
        Options:
            --field <field>    
               Lookup only specific fields in the output.
            --nocolor           
               Disable colored output.
            --json, -j          
               Output JSON format.
            --csv, -c           
               Output CSV format.
            --yaml, -y          
               Output YAML format.

        Examples:
            ipinfo <enter-asn>
            ipinfo <enter-asn> --field hostname,city


    ipinfo myip        Get details for your IP.
        Usage:
            ipinfo myip [<opts>]

        Options:
            --field <field>    
               Lookup only specific fields in the output.
            --nocolor           
               Disable colored output.
            --json, -j          
               Output JSON format.
            --csv, -c           
               Output CSV format.
            --yaml, -y          
               Output YAML format.

        Examples:
            ipinfo myip
            ipinfo myip --field hostname,city


    ipinfo bulk        Get details for multiple IPs in bulk.
        Usage:
            ipinfo bulk [<opts>] <ip | ip-range | cidr | filepath>

        Options:
            --field <field>    Lookup only specific fields in the output.
            --json, -j          Output JSON format.
            --csv, -c           Output CSV format.
            --yaml, -y          Output YAML format.

        Examples:
            # Lookup all IPs from stdin ('-' can be implied).
            $ ipinfo prips 8.8.8.0/24 | ipinfo bulk
            $ ipinfo prips 8.8.8.0/24 | ipinfo bulk -

            # Lookup all IPs in 2 files.
            $ ipinfo bulk /path/to/iplist1.txt /path/to/iplist2.txt

            # Lookup all IPs from CIDR.
            $ ipinfo bulk 8.8.8.0/24

            # Lookup all IPs from multiple CIDRs.
            $ ipinfo bulk 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

            # Lookup all IPs in an IP range.
            $ ipinfo bulk 8.8.8.0-8.8.8.255

            # Lookup all IPs from multiple sources simultaneously.
            $ ipinfo bulk 8.8.8.0-8.8.8.255 1.1.1.0/30 123.123.123.123 ips.txt


    ipinfo asn         Tools related to ASNs.
        Usage:
            ipinfo asn <cmd> [<opts>]

        Commands:
            bulk       lookup ASNs in bulk


    ipinfo summarize   Get summarized data for a group of IPs.
        Usage:
            ipinfo summarize [<opts>] <ip | ip-range | cidr | filepath>

        Options:
            --json, -j
               Output JSON format.
            --csv, -c           
               Output CSV format.
            --yaml, -y          
               Output YAML format.

        Examples:
            # Summarize all IPs from stdin ('-' can be implied).
            $ ipinfo prips 8.8.8.0/24 | ipinfo summarize
            $ ipinfo prips 8.8.8.0/24 | ipinfo summarize -

            # Summarize all IPs in 2 files.
            $ ipinfo summarize /path/to/iplist1.txt /path/to/iplist2.txt

            # Summarize all IPs from CIDR.
            $ ipinfo summarize 8.8.8.0/24

            # Summarize all IPs from multiple CIDRs.
            $ ipinfo summarize 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

            # Summarize all IPs in an IP range.
            $ ipinfo summarize 8.8.8.0-8.8.8.255

            # Summarize all IPs from multiple sources simultaneously.
            $ ipinfo summarize 8.8.8.0-8.8.8.255 1.1.1.0/30 123.123.123.123 ips.txt


    ipinfo map         Open a URL to a map showing the locations of a group of IPs.
        Usage:
            ipinfo map [<opts>] <ip | ip-range | cidr | filepath>

        Options:
            --no-browser 
               don't open the map link in the browser
			   default is false.

        Examples:
            # Map all IPs from stdin ('-' can be implied).
            $ ipinfo prips 8.8.8.0/24 | ipinfo map
            $ ipinfo prips 8.8.8.0/24 | ipinfo map -

            # Map all IPs in 2 files.
            $ ipinfo map /path/to/iplist1.txt /path/to/iplist2.txt

            # Map all IPs from CIDR.
            $ ipinfo map 8.8.8.0/24

            # Map all IPs from multiple CIDRs.
            $ ipinfo map 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

            # Map all IPs in an IP range.
            $ ipinfo map 8.8.8.0-8.8.8.255

            # Map all IPs from multiple sources simultaneously.
            $ ipinfo map 8.8.8.0-8.8.8.255 1.1.1.0/30 123.123.123.123 ips.txt


    ipinfo prips       Print IP list from CIDR or range.
        Usage:
            ipinfo prips [<opts>] <ip | ip-range | cidr | file>

        Examples:
            # List all IPs in a CIDR.
            $ ipinfo prips 8.8.8.0/24

            # List all IPs in multiple CIDRs.
            $ ipinfo prips 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

            # List all IPs in an IP range.
            $ ipinfo prips 8.8.8.0-8.8.8.255

            # List all IPs in multiple CIDRs and IP ranges.
            $ ipinfo prips 1.1.1.0/30 8.8.8.0-8.8.8.255 2.2.2.0/30 7.7.7.0,7.7.7.10

            # List all IPs from stdin input (newline-separated).
            $ echo '1.1.1.0/30\n8.8.8.0-8.8.8.255\n7.7.7.0,7.7.7.10' | ipinfo prips


    ipinfo grepip      Grep for IPs matching criteria from any source.
        Usage:
            ipinfo grepip [<opts>]

        Options:
            --only-matching, -o    
               print only matched IP in result line, excluding surrounding content.
            --no-filename, -h      
               don't print source of match in result lines when more than 1 source.
            --no-recurse           
               don't recurse into more directories in directory sources.

        Filters:
            --ipv4, -4               
               match only IPv4 addresses.
            --ipv6, -6               
               match only IPv6 addresses.
            --exclude-reserved, -x   
               exclude reserved/bogon IPs.
               full list can be found at https://ipinfo.io/bogon.

        Examples:
            ipinfo grepip -o <filename>
    

    ipinfo matchip      Print the overlapping IPs and subnets.
        Usage:
            ipinfo matchip [flags] <expression(s)> <file(s) | stdin | cidr | ip | ip-range>

        Flags:
           --expression, -e       CIDRs, ip-ranges to to check overlap with. Can be used multiple times.

        Examples:
            # Single expression + single file
            $ ipinfo matchip 127.0.0.1/24 file1.txt

            # Single expression + multiple files
            $ ipinfo matchip 127.0.0.1/24 file1.txt file2.txt file3.txt

            # Multi-expression + any files
            $ cat expression-list1.txt | ipinfo matchip -e 127.0.0.1/24 -e 8.8.8.8-8.8.9.10 -e - -e expression-list2.txt file1.txt file2.txt file3.txt

            # First arg is expression
            $ ipinfo matchip 8.8.8.8-8.8.9.10 8.8.0.0/16 8.8.0.10
    

    ipinfo gredomain      Grep for domains matching criteria from any source.
        Usage:
            ipinfo grepdomain [<opts>]

        Options:
            --only-matching, -o    
               print only matched IP in result line, excluding surrounding content.
            --no-filename, -h      
               don't print source of match in result lines when more than 1 source.
            --no-recurse           
               don't recurse into more directories in directory sources.

        Filters:
            --no-punycode, -n      
               do not convert domains to punycode.

        Examples:
            ipinfo grepdomain -o <filename>


    ipinfo cidr2range  Convert CIDRs to IP ranges.
        Usage:
            ipinfo cidr2range [<opts>] <cidr | filepath>

        Examples:
            # Get the range for CIDR 1.1.1.0/30.
            $ ipinfo cidr2range 1.1.1.0/30

            # Convert CIDR entries to IP ranges in 2 files.
            $ ipinfo cidr2range /path/to/file1.txt /path/to/file2.txt

            # Convert CIDR entries to IP ranges from stdin.
            $ cat /path/to/file1.txt | ipinfo cidr2range

            # Convert CIDR entries to IP ranges from stdin and a file.
            $ cat /path/to/file1.txt | ipinfo cidr2range /path/to/file2.txt


    ipinfo cidr2ip     Convert CIDRs to individual IPs within those CIDRs.
        Usage:
            ipinfo cidr2ip [<opts>] <cidrs | filepath>

        Examples:
            ipinfo cidr2ip 8.8.8.0/24


    ipinfo range2cidr  Convert IP ranges to CIDRs.
        Usage:
            ipinfo range2cidr [<opts>] <ip-range | filepath>

        Description:
            Accepts IP ranges and file paths to files containing IP ranges, converting
            them all to CIDRs (and multiple CIDRs if required).

            If a file is input, it is assumed that the IP range to convert is the first
            entry of each line. Other data is allowed and copied transparently.

            If multiple CIDRs are needed to represent an IP range on a line with other
            data, the data is copied per CIDR required. For example:

            in[0]: "1.1.1.0,1.1.1.2,other-data"
            out[0]: "1.1.1.0/31,other-data"
            out[1]: "1.1.1.2/32,other-data"

            IP ranges can have the form "<start><sep><end>" where "<sep>" can be "," or
            "-", and "<start>" and "<end>" can be any 2 IPs; order does not matter, but
            the resulting CIDRs are printed in the order they cover the range.

        Examples:
            # Get all CIDRs for range 1.1.1.0-1.1.1.2.
            $ ipinfo range2cidr 1.1.1.0-1.1.1.2
            $ ipinfo range2cidr 1.1.1.0,1.1.1.2

            # Convert all range entries to CIDRs in 2 files.
            $ ipinfo range2cidr /path/to/file1.txt /path/to/file2.txt

            # Convert all range entries to CIDRs from stdin.
            $ cat /path/to/file1.txt | ipinfo range2cidr

            # Convert all range entries to CIDRs from stdin and a file.
            $ cat /path/to/file1.txt | ipinfo range2cidr /path/to/file2.txt


    ipinfo range2ip    Convert IP ranges to individual IPs within those ranges.
        Usage: ipinfo range2ip [<opts>] <ip-range | filepath>

        Description:
            Accepts IP ranges and file paths to files containing IP ranges, converting
            them all to individual IPs within those ranges.

            $ ipinfo range2ip 8.8.8.0-8.8.8.255

            IP ranges can be of the form "<start><sep><end>" where "<sep>" can be "," or
            "-", and "<start>" and "<end>" can be any 2 IPs; order does not matter.


    ipinfo randip      Generates random IPs.
        Usage: 
            ipinfo randip [<opts>]

        Description:
            Generate random IPs.

            By default, generates 1 random IPv4 address with starting range 0.0.0.0 and 
            ending range 255.255.255.255, but can be configured to generate any number of 
            a combination of IPv4/IPv6 addresses within any range.

            Using --ipv6 or -6 without any starting or ending range will generate a IP 
            between range of :: to ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff.

            Note that only IPv4 or IPv6 IPs can be generated, but not both.

            When generating unique IPs, the range size must not be less than the desired 
            number of random IPs.

        Options:
            --num, -n
               number of IPs to generate.
            --ipv4, -4              
               generate IPv4 IPs.
            --ipv6, -6              
               generate IPv6 IPs.
            --start, -s             
               starting range of IPs. default: minimum IP possible for IP type selected.
            --end, -e               
               ending range of IPs. default: maximum IP possible for IP type selected.
            --exclude-reserved -x   
               exclude reserved/bogon IPs. full list can be found at https://ipinfo.io/bogon.
            --unique, -u            
               generate unique IPs.

        Examples:
            ipinfo randip --count 5


    ipinfo splitcidr   Splits a larger CIDR into smaller CIDRs.
        Usage:
            ipinfo splitcidr <cidr> <split>

        Examples:
            ipinfo splitcidr 8.8.8.0/24 25


    ipinfo mmdb        Read, import, and export mmdb files.
        Usage:
            ipinfo mmdb <cmd> [<opts>] [<args>]

        Commands:
            read        read data for IPs in an mmdb file.
                Usage:
                    ipinfo mmdb read [<opts>] <ip | ip-range | cidr | filepath> <mmdb>

                Options:
                    -f <format>, --format <format>   
                    the output format.
                    can be "json", "json-compact", "json-pretty", "tsv" or "csv". Note that "json" is short for "json-compact".
                    default: json.

            import      import data in non-mmdb format into mmdb.
                Usage:
                    ipinfo mmdb import [<opts>] [<input>] [<output>]
                
                Options:
                    Input/Output:
                        -i <fname>, --in <fname>
                            input file name. (e.g. data.csv or - for stdin)
                            must be in CSV, TSV or JSON.
                            default: stdin.
                        -o <fname>, --out <fname>
                            output file name. (e.g. sample.mmdb)
                            default: stdout.
                        -c, --csv
                            interpret input file as CSV.
                            by default, the .csv extension will turn this on.
                        -t, --tsv
                            interpret input file as TSV.
                            by default, the .tsv extension will turn this on.
                        -j, --json
                            interpret input file as JSON.
                            by default, the .json extension will turn this on.
                    
                    Fields:
                      One of the following fields flags, or other flags that implicitly specify
                      these, must be used, otherwise --fields-from-header is assumed.

                      The first field is always implicitly the network field, unless
                      --range-multicol is used, in which case the first 2 fields are considered
                      to be start_ip,end_ip.

                      When specifying --fields, do not specify the network column(s).

                      -f, --fields <comma-separated-fields>
                       explicitly specify the fields to assume exist in the input file.
                       example: col1,col2,col3
                       default: N/A.
                      --fields-from-header
                        assume first line of input file is a header, and set the fields as that.
                        default: true if no other field source is used, false otherwise.
                      --range-multicol
                        assume that the network field is actually two columns start_ip,end_ip.
                        default: false.
                      --joinkey-col
                        assume --range-multicol and that the 3rd column is join_key, and ignore
                        this column when converting to JSON.
                        default: false.
                      --no-fields
                        specify that no fields exist except the implicit network field.
                        when enabled, --no-network has no effect; the network field is written.
                        default: false.
                      --no-network
                        if --fields-from-header is set, then don't write the network field, which
                        is assumed to be the *first* field in the header.
                        default: false.

                    Meta:
                    --ip <4 | 6>
                        output file's ip version.
                        default: 6.
                    -s, --size <24 | 28 | 32>
                        size of records in the mmdb tree.
                        default: 32.
                    -m, --merge <none | toplevel | recurse>
                        the merge strategy to use when inserting entries that conflict.
                        none     => no merge; only replace conflicts.
                        toplevel => merge only top-level keys.
                        recurse  => recursively merge.
                        default: none.
                    --ignore-empty-values
                        if enabled, write into /0 with empty values for all fields, and for any
                        entry, don't write out a field whose value is the empty string.
                        default: false.
                    --disallow-reserved
                        disallow reserved networks to be added to the tree.
                        default: false.
                    --alias-6to4
                        enable the mapping of some IPv6 networks into the IPv4 network, e.g.
                        ::ffff:0:0/96, 2001::/32 & 2002::/16.
                        default: false.
                    --disable-metadata-pointers
                        some mmdb readers fail to properly read pointers within metadata. this
                        allows turning off such pointers.
                        default: true.
                
                Example:
                # Imports an input file and outputs an mmdb file with default configurations. 
                $ ipinfo mmdb import input.csv output.mmdb
                

            export      export data from mmdb format into non-mmdb format.
                Usage:
                    ipinfo mmdb export [<opts>] <mmdb_file> [<out_file>]
                
                Options:
                    Input/Output:
                     -o <fname>, --out <fname>
                        output file name. (e.g. out.csv)
                        default: <out_file> if specified, otherwise stdout.

                    Format:
                     -f <format>, --format <format>
                        the output file format.
                        can be "csv", "tsv" or "json".
                        default: csv if output file ends in ".csv", tsv if ".tsv",
                        json if ".json", otherwise csv.
                    --no-header
                        don't output the header for file formats that include one, like
                        CSV/TSV/JSON.
                        default: false.


            diff        see the difference between two mmdb files.
                Usage:
                    ipinfo mmdb diff [<opts>] <old> <new>
                
                Description:
                    Print subnet and record differences between two mmdb files (i.e. do set
                    difference "(new - old) U (old - new))".
                
                Options:
                    --subnets, -s
                        show subnets difference.
                    --records, -r
                        show records difference.


            metadata    print metadata from the mmdb file.
                Usage:
                    ipinfo mmdb metadata [<opts>] <mmdb_file>
                
                Options:
                    --nocolor
                       disable colored output.
                    
                  Format:
                    -f <format>, --format <format>
                      the metadata output format.
                      can be "pretty" or "json".
                      default: pretty.


            verify      check that the mmdb file is not corrupted or invalid.
                Usage:
                    ipinfo mmdb verify [<opts>] <mmdb_file>


            completion  install or output shell auto-completion script.


        Options:
            --nocolor  disable colored output.


    ipinfo calc        Evaluates a mathematical expression that may contain IP addresses.
        Usage:
            ipinfo calc <expression> [<opts>]
            
        Examples:
            ipinfo calc "2*2828-1"
            ipinfo calc "190.87.89.1*2"
            ipinfo calc "2001:0db8:85a3:0000:0000:8a2e:0370:7334*6"


    ipinfo tool        Misc. tools related to IPs, IP ranges, and CIDRs.
        Usage:
            ipinfo tool <cmd> [<opts>] [<args>]
        
        Commands:
            aggregate                    aggregate IPs, IP ranges, and CIDRs.
                Usage:
                    ipinfo tool aggregate [<opts>] <cidr | ip | ip-range | filepath>
                
                Description:
                    Accepts IPv4 IPs and CIDRs, aggregating them efficiently.

                    If input contains single IPs, it tries to merge them into the input CIDRs,
                    otherwise they are printed to the output as they are.

                    IP range can be of format <start-ip><SEP><end-ip>, where <SEP> can either
                    be a ',' or a '-'.


                Examples:
                    # Aggregate two CIDRs.
                    $ ipinfo tool aggregate 1.1.1.0/30 1.1.1.0/28

                    # Aggregate enteries from 2 files.
                    $ ipinfo tool aggregate /path/to/file1.txt /path/to/file2.txt

                    # Aggregate enteries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool aggregate

                    # Aggregate enteries from stdin and a file.
                    $ cat /path/to/file1.txt | ipinfo tool aggregate /path/to/file2.txt

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            next                         get the next IP of the input IP
                Usage:
                    ipinfo tool next [<opts>] <ip | filepath>
                

                Examples:
                    # Find the next IP for the given inputs 
                    $ ipinfo tool next 1.1.1.0

                    # Find next IP from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool next
                

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.
                

            prev                         get the previous IP of the input IP
                Usage:
                    ipinfo tool prev [<opts>] <ip | filepath>
                

                Examples:
                    # Find the previous IP for the given inputs 
                    $ ipinfo tool prev 1.1.1.0

                    # Find prev IP from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool prev

                
                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.

                
            is_v4                        reports whether input is an IPv4 address.
                Usage:
                    ipinfo tool is_v4 [<opts>] <cidr | ip | ip-range | filepath>


                Examples:
                    # Check CIDR.
                    $ ipinfo tool is_v4 1.1.1.0/30

                    # Check IP range.
                    $ ipinfo tool is_v4 1.1.1.0-1.1.1.244

                    # Check for file.
                    $ ipinfo tool is_v4 /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_v4

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_v6                        reports whether input is an IPv6 address.
                Usage:
                    ipinfo tool is_v6 [<opts>] <cidr | ip | ip-range | filepath>

                Examples:
                    # Check CIDR.
                    $ ipinfo tool is_v6 2001:db8::/32

                    # Check IP range.
                    $ ipinfo tool is_v6 2001:db8::1-2001:db8::10

                    # Check for file.
                    $ ipinfo tool is_v6 /path/to/file1.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_v6

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_valid                     reports whether an IP is valid.
                Usage:
                    ipinfo tool is_valid <ip>

                Examples:
                    ipinfo is_valid "190.87.89.1"
                    ipinfo is_valid "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
                    ipinfo is_valid "::"
                    ipinfo is_valid "0"
                    ipinfo is_valid ""

                
            is_one_ip                    checks whether a CIDR or IP Range contains exactly one IP.
                Usage:
                    ipinfo tool is_one_ip [<opts>] <cidr | ip | ip-range | filepath>

                Examples:
                    # Check CIDR.
                    $ ipinfo tool is_one_ip 1.1.1.0/30

                    # Check IP.
                    $ ipinfo tool is_one_ip 1.1.1.1

                    # Check IP range.
                    $ ipinfo tool is_one_ip 1.1.1.1-2.2.2.2
                    
                    # Check for file.
                    $ ipinfo tool is_one_ip /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_one_ip


            unmap                        returns ip with any IPv4-mapped IPv6 address prefix removed.
                Usage:
                    ipinfo tool unmap [<opts>] <ip>

                Examples:
                    ipinfo tool unmap "::ffff:8.8.8.8"
                    ipinfo tool unmap "192.180.32.1"
                    ipinfo tool unmap "::ffff:192.168.1.1"


            lower                        get start IP of IPs, IP ranges, and CIDRs.
                Usage:
                    ipinfo tool lower [<opts>] <cidr | ip | ip-range | filepath>

                Examples:
                    # Finds lower IP for CIDR.
                    $ ipinfo tool lower 192.168.1.0/24

                    # Finds lower IP for IP range.
                    $ ipinfo tool lower 1.1.1.0-1.1.1.244

                    # Finds lower IPs from stdin.
                    $ cat /path/to/file.txt | ipinfo tool lower

                    # Find lower IPs from file.
                    $ ipinfo tool lower /path/to/file1.txt
            upper                        get end IP of IPs, IP ranges, and CIDRs.
                Usage:
                    ipinfo tool upper [<opts>] <cidr | ip | ip-range | filepath>
                
                Examples:
                    # Finds upper IP for CIDR.
                    $ ipinfo tool upper 192.168.1.0/24

                    # Finds upper IP for IP range.
                    $ ipinfo tool upper 1.1.1.0-1.1.1.244

                    # Finds upper IPs from stdin.
                    $ cat /path/to/file.txt | ipinfo tool upper

                    # Find upper IPs from file.
                    $ ipinfo tool upper /path/to/file1.txt


            is_v4in6                     get whether the IP is an IPv4-mapped IPv6 address.
                Usage:
                    ipinfo tool is_v4in6 [<opts>] <ips>
                
                Examples:
                    ipinfo is_v4in6 "::7f00:1"
                    ipinfo is_v4in6 "::ffff"
                    ipinfo is_v4in6 "::ffff:8.8.8.8"
                    ipinfo is_v4in6 "::ffff:192.0.2.1

                
            ip2n                         converts an IPv4 or IPv6 address to its decimal representation.
                Usage:
                    ipinfo tool ip2n <ip>

                Examples:
                    ipinfo ip2n "190.87.89.1"
                    ipinfo ip2n "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
                    ipinfo ip2n "2001:0db8:85a3::8a2e:0370:7334"
                    ipinfo ip2n "::7334"
                    ipinfo ip2n "7334::"


            n2ip	                     evaluates a mathematical expression and converts it to an IPv4 or IPv6.
                Usage:
                    ipinfo tool n2ip [<opts>] <number>
                

                Examples:
                    ipinfo n2ip "4294967295 + 87"
                    ipinfo n2ip "4294967295" --ipv6
                    ipinfo n2ip -6 "201523715"
                    ipinfo n2ip "51922968585348276285304963292200960"
                    ipinfo n2ip "a:: - 4294967295"
                
                Options:
                    --ipv6, -6
                    force conversion to IPv6 address


            n2ip6	                     evaluates a mathematical expression and converts it to an IPv6.
                Usage:
                    ipinfo tool n2ip6 [<opts>] <number>
                
                
                Examples:
                    ipinfo n2ip6 "4294967295 + 87"
                    ipinfo n2ip6 "4294967295"
                    ipinfo n2ip6 "201523715"
                    ipinfo n2ip6 "51922968585348276285304963292200960"
                    ipinfo n2ip6 "a:: - 4294967295"
                

            prefix                       misc. prefix tools related to CIDRs.
                Usage:
                    ipinfo tool prefix <cmd> [<opts>] [<args>]

                Commands:
                    addr      returns the base IP address of a prefix.
                        Usage:
                            ipinfo tool prefix addr <cidr>
                        
                        Examples:
                            # CIDR Valid Examples.
                            $ ipinfo tool prefix addr 192.168.0.0/16
                            $ ipinfo tool prefix addr 10.0.0.0/8
                            $ ipinfo tool prefix addr 2001:0db8:1234::/48
                            $ ipinfo tool prefix addr 2606:2800:220:1::/64

                            # CIDR Invalid Examples.
                            $ ipinfo tool prefix addr 192.168.0.0/40
                            $ ipinfo tool prefix addr 2001:0db8:1234::/129


                    bits      returns the length of a prefix and reports -1 if invalid.
                        Usage:
                            ipinfo tool prefix bits <cidr>
                        
                        Examples:
                            # CIDR Valid Examples.
                            $ ipinfo tool prefix bits 192.168.0.0/16
                            $ ipinfo tool prefix bits 10.0.0.0/8
                            $ ipinfo tool prefix bits 2001:0db8:1234::/48
                            $ ipinfo tool prefix bits 2606:2800:220:1::/64

                            # CIDR Invalid Examples.
                            $ ipinfo tool prefix bits 192.168.0.0/40
                            $ ipinfo tool prefix bits 2001:0db8:1234::/129


                    masked    returns canonical form of a prefix, masking off non-high bits, and returns the zero if invalid.
                        Usage:
                            ipinfo tool prefix masked <cidr>
                        
                        Examples:
                            # CIDR Valid Examples.
                            $ ipinfo tool prefix masked 192.168.0.0/16
                            $ ipinfo tool prefix masked 10.0.0.0/8
                            $ ipinfo tool prefix masked 2001:0db8:1234::/48
                            $ ipinfo tool prefix masked 2606:2800:220:1::/64

                            # CIDR Invalid Examples.
                            $ ipinfo tool prefix masked 192.168.0.0/40
                            $ ipinfo tool prefix masked 2001:0db8:1234::/129


                    is_valid  reports whether a prefix is valid.
                        Usage:
                            ipinfo tool prefix is_valid <cidr>

                        Examples:
                            # CIDR Valid Examples.
                            $ ipinfo tool prefix is_valid 192.168.0.0/16
                            $ ipinfo tool prefix is_valid 10.0.0.0/8
                            $ ipinfo tool prefix is_valid 2001:0db8:1234::/48
                            $ ipinfo tool prefix is_valid 2606:2800:220:1::/64

                            # CIDR Invalid Examples.
                            $ ipinfo tool prefix is_valid 192.168.0.0/40
                            $ ipinfo tool prefix is_valid 2001:0db8:1234::/129

            is_loopback                  reports whether an IP is a valid loopback address.
                Usage:
                    ipinfo tool is_loopback [<opts>] <cidr | ip | ip-range | filepath>
                
                Examples:
                    $ ipinfo tool is_loopback 127.0.0.0
                    $ ipinfo tool is_loopback 160.0.0.0
                    $ ipinfo tool is_loopback ::1
                    $ ipinfo tool is_loopback fe08::2

                    # Check CIDR.
                    $ ipinfo tool is_loopback 127.0.0.0/32
                    $ ipinfo tool is_loopback 128.0.0.0/32
                    $ ipinfo tool is_loopback ::1/128
                    $ ipinfo tool is_loopback fe08::2/64

                    # Check IP range.
                    $ ipinfo tool is_loopback 127.0.0.1-127.20.1.244
                    $ ipinfo tool is_loopback 128.0.0.1-128.30.1.125

                    # Check for file.
                    $ ipinfo tool is_loopback /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_loopback

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_multicast                 reports whether an IP is a valid multicast address.
                Usage:
                    ipinfo tool is_multicast [<opts>] <cidr | ip | ip-range | filepath>
                
                Examples:
                    $ ipinfo tool is_multicast 239.0.0.0
                    $ ipinfo tool is_multicast 127.0.0.0
                    $ ipinfo tool is_multicast ff00::
                    $ ipinfo tool is_multicast ::1

                    # Check CIDR.
                    $ ipinfo tool is_multicast 239.0.0.0/32
                    $ ipinfo tool is_multicast 139.0.0.0/32
                    $ ipinfo tool is_multicast ff00::1/64
                    $ ipinfo tool is_multicast ::1/64

                    # Check IP range.
                    $ ipinfo tool is_multicast 239.0.0.0-239.255.255.1
                    $ ipinfo tool is_multicast 240.0.0.0-240.255.255.1
                    $ ipinfo tool is_multicast ff00::1-ff00::ffff
                    $ ipinfo tool is_multicast ::1-::ffff

                    # Check for file.
                    $ ipinfo tool is_multicast /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_multicast

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_unspecified               reports whether an IP is an unspecified address.
                Usage:
                    ipinfo tool is_unspecified [<opts>] <ip | filepath>
                
                Examples:
                    $ ipinfo tool is_unspecified 0.0.0.0
                    $ ipinfo tool is_unspecified 124.198.16.8
                    $ ipinfo tool is_unspecified ::
                    $ ipinfo tool is_unspecified fe08::1

                    # Check for file.
                    $ ipinfo tool is_unspecified /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_unspecified

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_global_unicast            reports whether an IP is a global unicast address.
                Usage:
                    ipinfo tool is_global_unicast [<opts>] <cidr | ip | ip-range | filepath>
                
                Examples:
                    $ ipinfo tool is_global_unicast 10.255.0.0
                    $ ipinfo tool is_global_unicast 255.255.255.255
                    $ ipinfo tool is_global_unicast 2000::1
                    $ ipinfo tool is_global_unicast ff00::1

                    # Check CIDR.
                    $ ipinfo tool is_global_unicast 10.255.0.0/32
                    $ ipinfo tool is_global_unicast 255.255.255.255/32
                    $ ipinfo tool is_global_unicast 2000::1/64
                    $ ipinfo tool is_global_unicast ff00::1/64

                    # Check IP range.
                    $ ipinfo tool is_global_unicast 10.0.0.1-10.8.95.6
                    $ ipinfo tool is_global_unicast 0.0.0.0-0.255.95.6
                    $ ipinfo tool is_global_unicast 2000::1-2000::ffff
                    $ ipinfo tool is_global_unicast ff00::1-ff00::ffff

                    # Check for file.
                    $ ipinfo tool is_global_unicast /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_global_unicast
                
                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_link_local_unicast        reports whether IP is a link local unicast.
                Usage:
                    ipinfo tool is_link_local_unicast [<opts>] <cidr | ip | ip-range | filepath>
                
                Examples:
                    $ ipinfo tool is_link_local_unicast 169.254.0.0
                    $ ipinfo tool is_link_local_unicast 127.0.0.0
                    $ ipinfo tool is_link_local_unicast fe80::1
                    $ ipinfo tool is_link_local_unicast ::1

                    # Check CIDR.
                    $ ipinfo tool is_link_local_unicast 169.254.0.0/32
                    $ ipinfo tool is_link_local_unicast 139.0.0.0/32
                    $ ipinfo tool is_link_local_unicast fe80::1/64
                    $ ipinfo tool is_link_local_unicast ::1/64

                    # Check IP range.
                    $ ipinfo tool is_link_local_unicast 169.254.0.0-169.254.255.1
                    $ ipinfo tool is_link_local_unicast 240.0.0.0-240.255.255.1
                    $ ipinfo tool is_link_local_unicast fe80::1-feb0::1
                    $ ipinfo tool is_link_local_unicast ::1-::ffff

                    # Check for file.
                    $ ipinfo tool is_link_local_unicast /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_link_local_unicast
                
                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_link_local_multicast      reports whether IP is a link local multicast address.
                Usage:
                    ipinfo tool is_link_local_multicast [<opts>] <cidr | ip | ip-range | filepath>
                
                Examples:
                    $ ipinfo tool is_link_local_multicast 224.0.0.0
                    $ ipinfo tool is_link_local_multicast 169.200.0.0
                    $ ipinfo tool is_link_local_multicast ff02::2
                    $ ipinfo tool is_link_local_multicast fe80::

                    # Check CIDR.
                    $ ipinfo tool is_link_local_multicast 224.0.0.0/32
                    $ ipinfo tool is_link_local_multicast 169.200.0.0/32
                    $ ipinfo tool is_link_local_multicast ff02::1/64
                    $ ipinfo tool is_link_local_multicast fe80::1/64

                    # Check IP range.
                    $ ipinfo tool is_link_local_multicast 224.0.0.1-224.255.255.255
                    $ ipinfo tool is_link_local_multicast 169.254.0.1-169.254.255.0
                    $ ipinfo tool is_link_local_multicast ff02::1-ff02::ffff
                    $ ipinfo tool is_link_local_multicast fe80::1-fe80::ffff

                    # Check for file.
                    $ ipinfo tool is_link_local_multicast /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_link_local_multicast
                
                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


            is_interface_local_multicast reports whether IP is an interface local multicast.
                Usage:
                    ipinfo tool is_interface_local_multicast [<opts>] <cidr | ip | ip-range | filepath>
                
                Examples:
                    $ ipinfo tool is_interface_local_multicast ff01::1
                    $ ipinfo tool is_interface_local_multicast ::1

                    # Check CIDR.
                    $ ipinfo tool is_interface_local_multicast ff01::ffff/32
                    $ ipinfo tool is_interface_local_multicast ff03::ffff/32

                    # Check IP range.
                    $ ipinfo tool is_interface_local_multicast ff01::1-ff01:ffff::1
                    $ ipinfo tool is_interface_local_multicast ff03::1-ff03:ffff::1

                    # Check for file.
                    $ ipinfo tool is_interface_local_multicast /path/to/file.txt 

                    # Check entries from stdin.
                    $ cat /path/to/file1.txt | ipinfo tool is_interface_local_multicast

                Options:
                    --quiet, -q
                        quiet mode; suppress additional output.


        Examples:
            ipinfo tool


    ipinfo download    Download free ipinfo database files.
        Usage:
            ipinfo download [<opts>] <database> [<output>]
        
        Examples:
            # Download country database in csv format.
            $ ipinfo download country -f csv > country.csv
            $ ipinfo download country-asn country_asn.mmdb

        Databases:
            asn            free ipinfo asn database.
            country        free ipinfo country database.
            country-asn    free ipinfo country-asn database.

        Options:
            --token <tok>, -t <tok>
                use <tok> as API token.

        Outputs:
            --compress, -c
                save the file in compressed format.
                default: false.
            --format, -f <mmdb | json | csv>
                output format of the database file.
                default: mmdb.


    ipinfo cache       Manage the cache.
        Usage:
            ipinfo cache [<opts>] [clear]

        Examples:
            # Clear all data currently in the cache.
            $ ipinfo cache clear


    ipinfo config      Manage the configuration.
        Usage:
            ipinfo config [<key>=<value>...]

        Examples:
            $ ipinfo config cache=disable
            $ ipinfo config token=testtoken cache=enable

        Configurations:
            cache=<enable | disable>
                Control whether the cache is enabled or disabled.
            open_browser=<enable | disable>
                Control whether the links should open the browser or not.
            token=<tok>
                Save a token for use when querying the API.
                (Token will not be validated).


    ipinfo quota       Print the request quota of your account.
        Usage:
            ipinfo quota [<opts>]
        
        Options:
            --detailed, -d
                show a detailed view of all available limits.
                default: false.


    ipinfo init        Login or signup for an account.
        Usage:
            ipinfo init [<opts>] [<token>]
        

        Examples:
            # Login command with token flag.
            $ ipinfo init --token <token>

            # Authentication without token flag.
            $ ipinfo init

        Options:
            --token <tok>, -t <tok>
                token to login with.
                (this is potentially unsafe; let the CLI prompt you instead).
            --no-check
                disable checking if the token is valid or not.
                default: false.


    ipinfo logout      Delete your current API token session.
        Usage:
            ipinfo logout


    ipinfo completion  Install or output shell auto-completion script.
        Usage:
            ipinfo completion [<opts>] [install | bash | zsh | fish]
        
        Description:
            Install or print out the code needed to do shell auto-completion.

            The current explicitly supported shells are:
            - bash
            - zsh
            - fish

        Examples:
            # Attempt auto-installation on any of the supported shells.
            $ ipinfo completion install

            # Output auto-completion script for bash for manual installation.
            $ ipinfo completion bash

            # Output auto-completion script for zsh for manual installation.
            $ ipinfo completion zsh

            # Output auto-completion script for fish for manual installation.
            $ ipinfo completion fish


    ipinfo version     Show the current version.
        Examples:
            ipinfo version


Options:
    --token <token>    Use <token> as the API token.
    --nocache          Do not use the cache.
    --version          Show the version number.
    --help, -h         Show this help message.
    --pretty, -p       Output pretty format.`
