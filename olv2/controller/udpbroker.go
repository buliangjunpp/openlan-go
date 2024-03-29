package controller

import (
    "log"
    "net"

    "github.com/danieldin95/openlan-go/olv2/openlanv2"
)

type UdpBroker struct {
    Listen string
    Udp *openlanv2.UdpSocket
    UUID string
    // 
    verbose bool
}

func NewUdpBroker(c *Config) (this *UdpBroker) {
    this = &UdpBroker {
        Listen: c.UdpListen,
        Udp: openlanv2.NewUdpSocket(c.UdpListen, c.Verbose),
        verbose: c.Verbose,
        UUID: openlanv2.GenUUID(16),
    }

    this.Udp.MaxSize = c.Ifmtu+int(openlanv2.HSIZE)

    if err := this.Udp.Listen(); err != nil {
        log.Printf("Error| NewUdpBroker.Listen %s\n", err)
    }

    return
}

func (this *UdpBroker) Close() error {
    return this.Udp.Close()
}

func (this *UdpBroker) GoRecv(doRecv func (raddr *net.UDPAddr, data []byte) error) {
    log.Printf("Info| UdpBroker.GoRecv from %s\n", this.Listen)

    for {
        addr, uuid, data, err := this.Udp.RecvMsg()
        if err != nil {
            log.Printf("Error| UdpBroker.GoRecv: %s\n", err)
            return 
        }
        if this.verbose {
            log.Printf("Info| UdpBroker.GoRecv from %s,%s: % x\n", addr, uuid, data)
        }

        if err:= doRecv(addr, data); err != nil {
            log.Printf("Error| UdpBroker.GoRecv from %s when doRecv %s\n", addr, err)
        }
    }
}

func (this *UdpBroker) DoSend(raddr *net.UDPAddr, uuid string, mesg *openlanv2.Message) error{
    return this.Udp.SendResp(raddr, uuid, mesg)
}



