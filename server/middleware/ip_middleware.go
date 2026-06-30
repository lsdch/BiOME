package middleware

import (
	"net"
	"net/netip"

	"github.com/danielgtaylor/huma/v2"
)

func extractIP(ctx huma.Context) *netip.Addr {
	headers := []string{
		"CF-Connecting-IP",
		"X-Real-IP",
		"X-Forwarded-For",
	}

	for _, h := range headers {
		if v := ctx.Header(h); v != "" {
			if ip, err := netip.ParseAddr(v); err == nil {
				return &ip
			}
		}
	}

	// fallback
	host, _, err := net.SplitHostPort(ctx.RemoteAddr())
	if err == nil {
		if ip, err := netip.ParseAddr(host); err == nil {
			return &ip
		}
	}

	return nil
}

func IPMiddleware(ctx huma.Context, next func(huma.Context)) {
	ip := extractIP(ctx)
	if ip != nil {
		ctx.SetHeader("client-ip", ip.String())
	}
	next(ctx)
}
