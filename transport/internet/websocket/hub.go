// +build !confonly

package websocket

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	http_proto "v2ray.com/core/common/protocol/http"
	"v2ray.com/core/common/session"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
)

type requestHandler struct {
	path string
	ln   *Listener
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:   4 * 1024,
	WriteBufferSize:  4 * 1024,
	HandshakeTimeout: time.Second * 4,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *requestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != h.path {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		newError("failed to convert to WebSocket connection").Base(err).WriteToLog()
		return
	}

	forwardedAddrs := http_proto.ParseXForwardedFor(request.Header)
	remoteAddr := conn.RemoteAddr()
	if len(forwardedAddrs) > 0 && forwardedAddrs[0].Family().IsIP() {
		remoteAddr = &net.TCPAddr{
			IP:   forwardedAddrs[0].IP(),
			Port: int(0),
		}
	}

	h.ln.addConn(newConnection(conn, remoteAddr))
}

type Listener struct {
	sync.Mutex
	server   http.Server
	listener net.Listener
	config   *Config
	addConn  internet.ConnHandler
	locker   *internet.FileLocker // for unix domain socket
}

func ListenWS(ctx context.Context, address net.Address, port net.Port, streamSettings *internet.MemoryStreamConfig, addConn internet.ConnHandler) (internet.Listener, error) {
	l := &Listener{
		addConn: addConn,
	}
	wsSettings := streamSettings.ProtocolSettings.(*Config)
	l.config = wsSettings
	if l.config != nil && streamSettings.SocketSettings != nil {
		streamSettings.SocketSettings.AcceptProxyProtocol = l.config.AcceptProxyProtocol
	}
	var listener net.Listener
	var err error
	if port == net.Port(0) { //unix
		listener, err = internet.ListenSystem(ctx, &net.UnixAddr{
			Name: address.Domain(),
			Net:  "unix",
		}, streamSettings.SocketSettings)
		if err != nil {
			return nil, newError("failed to listen unix domain socket(for WS) on ", address).Base(err)
		}
		newError("listening unix domain socket(for WS) on ", address).WriteToLog(session.ExportIDToError(ctx))
		locker := ctx.Value(address.Domain())
		if locker != nil {
			l.locker = locker.(*internet.FileLocker)
		}
	} else { //tcp
		listener, err = internet.ListenSystem(ctx, &net.TCPAddr{
			IP:   address.IP(),
			Port: int(port),
		}, streamSettings.SocketSettings)
		if err != nil {
			return nil, newError("failed to listen TCP(for WS) on ", address, ":", port).Base(err)
		}
		newError("listening TCP(for WS) on ", address, ":", port).WriteToLog(session.ExportIDToError(ctx))
	}

	// if wsSettings.AcceptProxyProtocol {
	// 	policyFunc := func(upstream net.Addr) (proxyproto.Policy, error) { return proxyproto.REQUIRE, nil }
	// 	listener = &proxyproto.Listener{Listener: listener, Policy: policyFunc}
	// 	newError("accepting PROXY protocol").AtWarning().WriteToLog(session.ExportIDToError(ctx))
	// }

	if config := v2tls.ConfigFromStreamSettings(streamSettings); config != nil {
		if tlsConfig := config.GetTLSConfig(); tlsConfig != nil {
			listener = tls.NewListener(listener, tlsConfig)
		}
	}
	l.listener = listener
	// l := &Listener{
	// 	config:   wsSettings,
	// 	addConn:  addConn,
	// 	listener: listener,
	// }

	l.server = http.Server{
		Handler: &requestHandler{
			path: wsSettings.GetNormalizedPath(),
			ln:   l,
		},
		ReadHeaderTimeout: time.Second * 4,
		MaxHeaderBytes:    2048,
	}

	go func() {
		if err := l.server.Serve(l.listener); err != nil {
			newError("failed to serve http for WebSocket").Base(err).AtWarning().WriteToLog(session.ExportIDToError(ctx))
		}
	}()

	return l, err
}

// Addr implements net.Listener.Addr().
func (ln *Listener) Addr() net.Addr {
	return ln.listener.Addr()
}

// Close implements net.Listener.Close().
func (ln *Listener) Close() error {
	if ln.locker != nil {
		ln.locker.Release()
	}
	return ln.listener.Close()
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, ListenWS))
}
