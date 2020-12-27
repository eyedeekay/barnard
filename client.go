package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
	"layeh.com/gumble/gumbleutil"

	"github.com/eyedeekay/goSam"
)

func (b *Barnard) start() {
	b.Config.Attach(gumbleutil.AutoBitrate)
	b.Config.Attach(b)

	var err error
	if strings.HasSuffix(b.Address, ".i2p") {
		sam, err := goSam.NewClientFromOptions(
			goSam.SetAddr(addr),
			goSam.SetInQuantity(5),
			goSam.SetOutQuantity(5),
			goSam.SetInBackups(3),
			goSam.SetOutBackups(3),
			goSam.SetInLength(2),
			goSam.SetOutLength(2),
			goSam.SetInVariance(-1),
			goSam.SetOutVariance(-1),
			goSam.SetCloseIdle(false),
			goSam.SetDebug(false),
			goSam.SetLeaseSetEncType("4"),
			goSam.SetSAMMinVersion(0),
			goSam.SetSAMMaxVersion(1),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		_, err = gumble.DialWithDialer(sam.Dial, b.Address, b.Config, &b.TLSConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		_, err = gumble.DialWithDialer(new(net.Dialer), b.Address, b.Config, &b.TLSConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	// Audio
	if os.Getenv("ALSOFT_LOGLEVEL") == "" {
		os.Setenv("ALSOFT_LOGLEVEL", "0")
	}
	if stream, err := gumbleopenal.New(b.Client); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		b.Stream = stream
	}
}

func (b *Barnard) OnConnect(e *gumble.ConnectEvent) {
	b.Client = e.Client

	b.Ui.SetActive(uiViewInput)
	b.UiTree.Rebuild()
	b.Ui.Refresh()

	b.UpdateInputStatus(fmt.Sprintf("To: %s", e.Client.Self.Channel.Name))
	b.AddOutputLine(fmt.Sprintf("Connected to %s", b.Client.Conn.RemoteAddr()))
	if e.WelcomeMessage != nil {
		b.AddOutputLine(fmt.Sprintf("Welcome message: %s", esc(*e.WelcomeMessage)))
	}
}

func (b *Barnard) OnDisconnect(e *gumble.DisconnectEvent) {
	var reason string
	switch e.Type {
	case gumble.DisconnectError:
		reason = "connection error"
	}
	if reason == "" {
		b.AddOutputLine("Disconnected")
	} else {
		b.AddOutputLine("Disconnected: " + reason)
	}
	b.UiTree.Rebuild()
	b.Ui.Refresh()
}

func (b *Barnard) OnTextMessage(e *gumble.TextMessageEvent) {
	b.AddOutputMessage(e.Sender, e.Message)
}

func (b *Barnard) OnUserChange(e *gumble.UserChangeEvent) {
	if e.Type.Has(gumble.UserChangeChannel) && e.User == b.Client.Self {
		b.UpdateInputStatus(fmt.Sprintf("To: %s", e.User.Channel.Name))
	}
	b.UiTree.Rebuild()
	b.Ui.Refresh()
}

func (b *Barnard) OnChannelChange(e *gumble.ChannelChangeEvent) {
	b.UiTree.Rebuild()
	b.Ui.Refresh()
}

func (b *Barnard) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	var info string
	switch e.Type {
	case gumble.PermissionDeniedOther:
		info = e.String
	case gumble.PermissionDeniedPermission:
		info = "insufficient permissions"
	case gumble.PermissionDeniedSuperUser:
		info = "cannot modify SuperUser"
	case gumble.PermissionDeniedInvalidChannelName:
		info = "invalid channel name"
	case gumble.PermissionDeniedTextTooLong:
		info = "text too long"
	case gumble.PermissionDeniedTemporaryChannel:
		info = "temporary channel"
	case gumble.PermissionDeniedMissingCertificate:
		info = "missing certificate"
	case gumble.PermissionDeniedInvalidUserName:
		info = "invalid user name"
	case gumble.PermissionDeniedChannelFull:
		info = "channel full"
	case gumble.PermissionDeniedNestingLimit:
		info = "nesting limit"
	}
	b.AddOutputLine(fmt.Sprintf("Permission denied: %s", info))
}

func (b *Barnard) OnUserList(e *gumble.UserListEvent) {
}

func (b *Barnard) OnACL(e *gumble.ACLEvent) {
}

func (b *Barnard) OnBanList(e *gumble.BanListEvent) {
}

func (b *Barnard) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (b *Barnard) OnServerConfig(e *gumble.ServerConfigEvent) {
}
