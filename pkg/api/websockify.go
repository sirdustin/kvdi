package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tinyzimmer/kvdi/pkg/apis/kvdi/v1alpha1"
	"github.com/tinyzimmer/kvdi/pkg/util/tlsutil"

	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
)

var clientTLSConfig *tls.Config

func init() {
	var err error
	clientTLSConfig, err = tlsutil.NewClientTLSConfig()
	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (d *desktopAPI) mtlsWebsockify(w http.ResponseWriter, r *http.Request) {
	endpointURL := getEndpointURL(r)
	apiLogger.Info(fmt.Sprintf("Starting new mTLS websocket proxy with %s", endpointURL))
	proxy := websocketproxy.NewProxy(endpointURL)
	proxy.Dialer = &websocket.Dialer{
		TLSClientConfig: clientTLSConfig,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	proxy.Upgrader = &upgrader
	proxy.ServeHTTP(w, r)
}

func getEndpointURL(r *http.Request) *url.URL {
	nn := getNamespacedNameFromRequest(r)
	url, _ := url.Parse(fmt.Sprintf("wss://%s.%s.%s:%d", nn.Name, nn.Name, nn.Namespace, v1alpha1.WebPort))
	return url
}