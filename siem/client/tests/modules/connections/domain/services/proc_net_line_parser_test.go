package services

import (
    "testing"

    "siem-client/src/modules/connections/domain/services"
    "siem-client/tests"
)

const (
    tcpListenLine      = "   0: 0100007F:0277 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 22512 1 0000000000000000 100 0 0 10 0"
    tcpEstablishedLine = "   1: 0501A8C0:CD2C 01FAFA8E:01BB 01 00000000:00000000 02:000A7D8C 00000000  1000        0 98765 2 0000000000000000 20 4 30 10 -1"
    tcpTimeWaitLine    = "   2: 0501A8C0:CD2D 01FAFA8E:01BB 06 00000000:00000000 03:00001234 00000000     0        0 0 3 0000000000000000"
    tcp6ListenLine     = "   0: 00000000000000000000000001000000:0277 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 22513 1 0000000000000000 100 0 0 10 0"
    tcp6MappedLine     = "   1: 0000000000000000FFFF00000501A8C0:0016 0000000000000000FFFF00000A01A8C0:D431 01 00000000:00000000 02:00000000 00000000     0        0 4444 1 0000000000000000 20 4 30 10 -1"
    udpUnconnLine      = "  10: 00000000:0044 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 3333 2 0000000000000000 0"
    udpConnectedLine   = "  11: 0501A8C0:9C40 08080808:0035 01 00000000:00000000 00:00000000 00000000  1000        0 5555 2 0000000000000000 0"
    procNetHeaderLine  = "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode"
)

func TestParsesTcpListeningSocket(t *testing.T) {
    socketEntity := services.ProcNetLineParserGetInstance().GetSocketEntityOrNil("tcp", tcpListenLine)
    tests.AssertEquals(t, "tcp", socketEntity.GetProtocol())
    tests.AssertEquals(t, 4, socketEntity.GetIpVersion())
    tests.AssertEquals(t, "127.0.0.1:631", socketEntity.GetLocalEndpoint())
    tests.AssertEquals(t, "*", socketEntity.GetRemoteEndpoint())
    tests.AssertEquals(t, "LISTEN", socketEntity.GetState())
    tests.AssertEquals(t, uint64(22512), socketEntity.GetInode())
    tests.AssertEquals(t, true, socketEntity.IsListening())
    tests.AssertEquals(t, true, socketEntity.IsLoopback())
}

func TestParsesTcpEstablishedConnection(t *testing.T) {
    socketEntity := services.ProcNetLineParserGetInstance().GetSocketEntityOrNil("tcp", tcpEstablishedLine)
    tests.AssertEquals(t, "192.168.1.5:52524", socketEntity.GetLocalEndpoint())
    tests.AssertEquals(t, "142.250.250.1:443", socketEntity.GetRemoteEndpoint())
    tests.AssertEquals(t, "ESTABLISHED", socketEntity.GetState())
    tests.AssertEquals(t, 1000, socketEntity.GetUid())
    tests.AssertEquals(t, false, socketEntity.IsListening())
}

func TestParsesTimeWaitWithInodeZero(t *testing.T) {
    socketEntity := services.ProcNetLineParserGetInstance().GetSocketEntityOrNil("tcp", tcpTimeWaitLine)
    tests.AssertEquals(t, "TIME_WAIT", socketEntity.GetState())
    tests.AssertEquals(t, uint64(0), socketEntity.GetInode())
}

func TestParsesTcp6Loopback(t *testing.T) {
    socketEntity := services.ProcNetLineParserGetInstance().GetSocketEntityOrNil("tcp6", tcp6ListenLine)
    tests.AssertEquals(t, 6, socketEntity.GetIpVersion())
    tests.AssertEquals(t, "[::1]:631", socketEntity.GetLocalEndpoint())
    tests.AssertEquals(t, true, socketEntity.IsLoopback())
}

func TestIpv4MappedInTcp6IsShownAsIpv4(t *testing.T) {
    socketEntity := services.ProcNetLineParserGetInstance().GetSocketEntityOrNil("tcp6", tcp6MappedLine)
    tests.AssertEquals(t, "192.168.1.5:22", socketEntity.GetLocalEndpoint())
    tests.AssertEquals(t, "192.168.1.10:54321", socketEntity.GetRemoteEndpoint())
}

func TestUdpWithoutRemoteIsUnconnListening(t *testing.T) {
    socketEntity := services.ProcNetLineParserGetInstance().GetSocketEntityOrNil("udp", udpUnconnLine)
    tests.AssertEquals(t, "udp", socketEntity.GetProtocol())
    tests.AssertEquals(t, "UNCONN", socketEntity.GetState())
    tests.AssertEquals(t, true, socketEntity.IsListening())
    tests.AssertEquals(t, "0.0.0.0:68", socketEntity.GetLocalEndpoint())
}

func TestUdpWithRemoteIsAConnection(t *testing.T) {
    socketEntity := services.ProcNetLineParserGetInstance().GetSocketEntityOrNil("udp", udpConnectedLine)
    tests.AssertEquals(t, "ESTABLISHED", socketEntity.GetState())
    tests.AssertEquals(t, false, socketEntity.IsListening())
    tests.AssertEquals(t, "8.8.8.8:53", socketEntity.GetRemoteEndpoint())
}

func TestHeaderAndGarbageAreIgnored(t *testing.T) {
    parser := services.ProcNetLineParserGetInstance()
    tests.AssertNil(t, toAny(parser.GetSocketEntityOrNil("tcp", procNetHeaderLine)))
    tests.AssertNil(t, toAny(parser.GetSocketEntityOrNil("tcp", "")))
    tests.AssertNil(t, toAny(parser.GetSocketEntityOrNil("tcp", "   0: ZZZZ:0277 00000000:0000 0A 0:0 0:0 0 0 0 1")))
    tests.AssertNil(t, toAny(parser.GetSocketEntityOrNil("raw", tcpListenLine)))
}

// toAny un puntero nil tipado no es nil como any: AssertNil necesita el nil de verdad.
func toAny[T any](pointer *T) any {
    if pointer == nil {
        return nil
    }
    return pointer
}
