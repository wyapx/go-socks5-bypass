package src

import (
	"errors"
	"net"
	"net/http"

	"github.com/babolivier/go-doh-client"
	"golang.org/x/net/context"
)

// NameResolver is used to implement custom name resolution
type NameResolver interface {
	Resolve(ctx context.Context, name string) (context.Context, net.IP, error)
}

// DNSResolver uses the system DNS to resolve host names
type DNSResolver struct{}

//func (d DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
//	addr, err := net.ResolveIPAddr("ip", name)
//	if err != nil {
//		return ctx, nil, err
//	}
//	return ctx, addr.IP, err
//}

func Resolver(name string) (net.IP, error) {
	ip := net.ParseIP(name)
	if ip != nil {
		return ip, nil
	}
	client := &http.Client{}
	resolver := doh.Resolver{
		"208.67.220.220",
		doh.IN,
		client,
	}

	a, _, err := resolver.LookupA(name)
	if err != nil {
		return nil, err
	}
	if len(a) == 0 {
		return nil, errors.New("no IP addresses found")
	}

	return net.ParseIP(a[0].IP4), nil
}

func (d DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	addr, err := Resolver(name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, addr, err
}
