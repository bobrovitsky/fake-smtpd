
package main

import (
    "github.com/bobrovitsky/smtpd"
    "log"
    "net"
    "io"
    "io/ioutil"
)

var (
    listen = "0.0.0.0:25"
)

type Session struct {
	Conn net.Conn
}

func (s *Session) Connect(source string) error {
    return nil
}

func (s *Session) Hello(hostname string) error {
    return nil
}

func (s *Session) AuthUser(identity, username string) (password string, err error) {
    return "", nil
}

func (s *Session) AuthSuccess(username string) {
}

func (s *Session) Sender(address string) error {
    return nil
}

func (s *Session) Recipient(address string) error {
    return nil
}

func (s *Session) Message(reader io.Reader) error {
    io.Copy(ioutil.Discard, reader)
	return nil
}

func main() {

    srv := smtpd.Server {
        Hostname: "fake-smtpd",
        Pipelining: true,
    }

    listener, err := net.Listen("tcp", listen)

    log.Println("start listen ", listen)

    if err != nil {
        log.Fatal(err)
    }

    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }

        go func() {

            session := &Session{
                Conn: conn,
            }

            srv.ServeSMTP(conn, session)
            conn.Close()

        }()
    }
}
