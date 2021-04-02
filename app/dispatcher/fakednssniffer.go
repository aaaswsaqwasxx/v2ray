// +build !confonly

package dispatcher

import (
	"context"

	core "github.com/v2fly/v2ray-core/v4"
	"github.com/v2fly/v2ray-core/v4/common"
	"github.com/v2fly/v2ray-core/v4/common/net"
	"github.com/v2fly/v2ray-core/v4/common/session"
	"github.com/v2fly/v2ray-core/v4/features/dns"
)

// newFakeDNSSniffer Create a Fake DNS metadata sniffer
func newFakeDNSSniffer(ctx context.Context) (snifferWithMetadata, error) {
	var fakeDNSEngine dns.FakeDNSEngine
	err := core.RequireFeatures(ctx, func(fdns dns.FakeDNSEngine) {
		fakeDNSEngine = fdns
	})
	if err != nil {
		return snifferWithMetadata{}, err
	}
	if fakeDNSEngine == nil {
		errNotInit := newError("FakeDNSEngine is not initialized, but such a sniffer is used").AtError()
		return snifferWithMetadata{}, errNotInit
	}
	return snifferWithMetadata{domainSniffer: func(ctx context.Context, bytes []byte) (SniffResult, error) {
		Target := session.OutboundFromContext(ctx).Target
		if Target.Network == net.Network_TCP || Target.Network == net.Network_UDP {
			domainFromFakeDNS := fakeDNSEngine.GetDomainFromFakeDNS(Target.Address)
			if domainFromFakeDNS != "" {
				newError("fake dns got domain: ", domainFromFakeDNS, " for ip: ", Target.Address.String()).WriteToLog(session.ExportIDToError(ctx))
				return &fakeDNSSniffResult{domainName: domainFromFakeDNS}, nil
			}
		}
		return nil, common.ErrNoClue
	}, protocolSniffer: func(ctx context.Context, bytes []byte) (SniffResult, error) {
		return nil, common.ErrNoClue
	}, metadataSniffer: true}, nil
}

type fakeDNSSniffResult struct {
	domainName string
}

func (fakeDNSSniffResult) Protocol() string {
	return "fakedns"
}

func (f fakeDNSSniffResult) Domain() string {
	return f.domainName
}
