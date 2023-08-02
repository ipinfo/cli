package lib

import (
	"errors"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
	"net"
	"strings"
)

// CmdBulkIPFlags are flags expected by CmdBulkIp.
type CmdBulkIPFlags struct {
	Token   string
	nocache bool
	help    bool
	Field   []string
	json    bool
	Csv     bool
	Yaml    bool
}

// Init initializes the common flags available to CmdBulkIP with sensible
// defaults.
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdBulkIPFlags) Init() {
	_h := "see description in --help"
	pflag.StringVarP(
		&f.Token,
		"token", "t", "",
		_h,
	)
	pflag.BoolVarP(
		&f.nocache,
		"nocache", "", false,
		_h,
	)
	pflag.BoolVarP(
		&f.help,
		"help", "h", false,
		_h,
	)
	pflag.StringSliceVarP(
		&f.Field,
		"field", "f", []string{},
		_h,
	)
	pflag.BoolVarP(
		&f.json,
		"json", "j", false,
		_h,
	)
	pflag.BoolVarP(
		&f.Csv,
		"csv", "c", false,
		_h,
	)
	pflag.BoolVarP(
		&f.Yaml,
		"yaml", "y", false,
		_h,
	)
}

func CmdBulkIP(ii *ipinfo.Client, args []string) (ipinfo.BatchCore, error) {
	var ips []net.IP
	if !validateWithFunctions(args, []func(string) bool{StrIsIPStr}) {
		return nil, ErrInvalidInput
	}

	if strings.TrimSpace(ii.Token) == "" {
		return nil, errors.New("bulk lookups require a token; login via `ipinfo init`")
	}

	for _, ipstr := range args {
		ips = append(ips, net.ParseIP(ipstr))
	}

	return ii.GetIPInfoBatch(ips, ipinfo.BatchReqOpts{
		TimeoutPerBatch:              60 * 30, // 30min
		ConcurrentBatchRequestsLimit: 20,
	})
}
