module github.com/sdslabs/portkey

go 1.14

require (
	github.com/DataDog/zstd v1.4.5
	github.com/pion/quic v0.1.4
	github.com/pion/webrtc/v3 v3.0.0-beta.15
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
)

replace github.com/pion/quic => github.com/vrngupta53/quic v0.1.5-0.20201226123440-2e529071b048
