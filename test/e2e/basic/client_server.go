package basic

import (
	"fmt"

	"github.com/fatedier/frp/test/e2e/framework"
	"github.com/fatedier/frp/test/e2e/framework/consts"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("[Feature: Client-Server]", func() {
	f := framework.NewDefaultFramework()

	Describe("Protocol", func() {
		supportProtocols := []string{"tcp", "kcp", "websocket"}
		for _, protocol := range supportProtocols {
			It(protocol, func() {
				framework.Logf("protocol: %s", protocol)
				serverConf := consts.DefaultServerConfig
				clientConf := consts.DefaultClientConfig

				serverConf += fmt.Sprintf(`
				protocol = %s
				`, protocol)

				clientConf += fmt.Sprintf(`
				protocol = %s

				[tcp]
				type = tcp
				local_port = {{ .%s }}
				remote_port = {{ .%s }}

				[udp]
				type = udp
				local_port = {{ .%s }}
				remote_port = {{ .%s }}
				`, protocol,
					framework.TCPEchoServerPort, framework.GenPortName("TCP"),
					framework.UDPEchoServerPort, framework.GenPortName("UDP"),
				)

				f.RunProcesses([]string{serverConf}, []string{clientConf})

				framework.ExpectTCPRequest(f.UsedPorts[framework.GenPortName("TCP")],
					[]byte(consts.TestString), []byte(consts.TestString), connTimeout, "tcp proxy")
				framework.ExpectUDPRequest(f.UsedPorts[framework.GenPortName("UDP")],
					[]byte(consts.TestString), []byte(consts.TestString), connTimeout, "udp proxy")
			})
		}
	})
})
