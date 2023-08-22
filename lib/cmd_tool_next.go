package lib

import (
	"fmt"
	"net"
	"github.com/spf13/pflag"
)

type CmdToolNextFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolNextFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.Quiet,
		"quiet", "q", false,
		"quiet mode; suppress additional output.",
	)
}

func CmdToolNext(
	f CmdToolNextFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	increment := 1
	actionFunc := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			UpdateIPAddress(input, increment)
		default:
			return ErrNotIP
		}
		return nil
	}
	err := GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
	}
	
	return nil
}

func UpdateIPAddress(input string, delta int) {
	ip := net.ParseIP(input)
	if ip != nil {
		if ip.To4() != nil {
			ipInt := ipToUint32(ip)
			nextPrevIPInt := ipInt + uint32(delta)
			adjustedIPInt := adjustIPUint32(nextPrevIPInt)
			nextPrevIP := uint32ToIP(adjustedIPInt)
			fmt.Println(nextPrevIP)
		} else {
			ipInt := ipToBigInt(ip)
			deltaBigInt := new(big.Int).SetInt64(int64(delta))
			nextPrevIPInt := new(big.Int).Add(ipInt, deltaBigInt)
			adjustedIPInt := adjustIPBigInt(nextPrevIPInt)
			nextPrevIP := bigIntToIP(adjustedIPInt)
			fmt.Println(nextPrevIP)
		}
	}
}

func ipToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

func uint32ToIP(ipInt uint32) net.IP {
	ip := make(net.IP, net.IPv4len)
	binary.BigEndian.PutUint32(ip, ipInt)
	return ip
}

func adjustIPUint32(ipInt uint32) uint32 {
	if ipInt > math.MaxUint32 {
		return ipInt - math.MaxUint32
	}
	if ipInt < 0 {
		return ipInt + math.MaxUint32
	}
	return ipInt
}

func ipToBigInt(ip net.IP) *big.Int {
	ipInt := new(big.Int)
	ipInt.SetBytes(ip)
	return ipInt
}

func bigIntToIP(ipInt *big.Int) net.IP {
	ip := make(net.IP, net.IPv6len)
	ipIntBytes := ipInt.Bytes()
	if len(ipIntBytes) > net.IPv6len {
		ipIntBytes = ipIntBytes[len(ipIntBytes)-net.IPv6len:]
	}
	copy(ip[net.IPv6len-len(ipIntBytes):], ipIntBytes)
	return ip
}

func adjustIPBigInt(ipInt *big.Int) *big.Int {
	if ipInt.Cmp(maxIPv6BigInt) == 0 {
		return big.NewInt(0)
	}
	if ipInt.Cmp(big.NewInt(0)) == 0 {
		return maxIPv6BigInt
	}
	return ipInt
}

var (
	maxIPv6BigInt, _ = new(big.Int).SetString("340282366920938463463374607431768211455", 10) // 2^128 - 1
)