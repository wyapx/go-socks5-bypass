package src

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/babolivier/go-doh-client"
	gocache "github.com/patrickmn/go-cache"
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

var client = &http.Client{}
var resolver = doh.Resolver{
	Host:       "208.67.220.220",
	Class:      doh.IN,
	HTTPClient: client,
}
var cache = gocache.New(5*time.Minute, 10*time.Minute)

func (d DNSResolver) Resolver4(name string) (net.IP, uint32, error) {
	a, t, err := resolver.LookupA(name)
	if err != nil {
		return nil, 0, err
	}
	if len(a) == 0 {
		return nil, 0, errors.New("no IP addresses found")
	}

	return net.ParseIP(a[0].IP4), t[0], nil
}

func (d DNSResolver) Resolver6(name string) (net.IP, uint32, error) {
	a, t, err := resolver.LookupAAAA(name)
	if err != nil {
		return nil, 0, err
	}
	if len(a) == 0 {
		//return nil, errors.New("no IP addresses found")
		return d.Resolver4(name)
	}

	return net.ParseIP(a[0].IP6), t[0], nil
}

func (d DNSResolver) Resolver(name string) (net.IP, error) {
	ip := net.ParseIP(name)
	if ip != nil {
		return ip, nil
	}

	_ip, found := cache.Get(name)
	if found {
		return _ip.(net.IP), nil
	} else {
		ip, ttl, err := d.Resolver4(name)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("cache for %s:%v\n", name, ip)
		cache.Set(name, ip, time.Duration(ttl)*time.Second)
		return ip, nil
	}
}

func (d DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	addr, err := d.Resolver(name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, addr, err
}
