
package main

import (
    "github.com/bobrovitsky/smtpd"
    "log"
    "net"
    "io"
    "io/ioutil"
    "math/rand"
    "time"
    "fmt"
)

var (
    listen = "0.0.0.0:25"
    hard_percent = 5
    soft_percent = 15
)

var hard_bounce = []string{
    "554 rejected due to spam content",
    "550 deliver error: dd This user doesn't have yahoo.com account",
    "554 RLY:B1",
    "554 Invalid recipient",
    "554 Message permanently rejected",
}

var soft_bounce = []string{
    "421 Service not available",
    "421 Mailbox unavailable",
    "452 Insufficient system storage",
    "452 Mailbox full",
}

const (
    step_connect = iota
    step_hello
    step_sender
    step_recipient
    step_message
)

type Answer struct {
    step int
    resp error
    num int
}

type Session struct {
	Conn net.Conn
    Answer *Answer
}

func GetAnswer() *Answer {
    var num = rand.Intn(100)

    if hard_percent >= num {
        return &Answer {
            step: rand.Intn(5),
            resp: fmt.Errorf(hard_bounce[rand.Intn(len(hard_bounce))]),
            num: num,
        }
    }

    if num > hard_percent && num <= (soft_percent + hard_percent) {
        return &Answer {
            step: rand.Intn(5),
            resp: fmt.Errorf(soft_bounce[rand.Intn(len(soft_bounce))]),
            num: num,
        }
    }

    return nil
}

func (s *Session) Connect(source string) error {
    if s.Answer != nil && s.Answer.step == 0 {
        return s.Answer.resp
    }

    return nil
}

func (s *Session) Hello(hostname string) error {
    if s.Answer != nil && s.Answer.step == 1 {
        return s.Answer.resp
    }

    return nil
}

func (s *Session) AuthUser(identity, username string) (password string, err error) {
    return "", nil
}

func (s *Session) AuthSuccess(username string) {
}

func (s *Session) Sender(address string) error {
    if s.Answer != nil && s.Answer.step == 2 {
        return s.Answer.resp
    }

    return nil
}

func (s *Session) Recipient(address string) error {
    if s.Answer != nil && s.Answer.step == 3 {
        return s.Answer.resp
    }

    return nil
}

func (s *Session) Message(reader io.Reader) error {
    io.Copy(ioutil.Discard, reader)

    if s.Answer != nil && s.Answer.step == 4 {
        return s.Answer.resp
    }

	return nil
}

func main() {

    rand.Seed(time.Now().UTC().UnixNano())

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
                Answer: GetAnswer(),
            }

            srv.ServeSMTP(conn, session)
            conn.Close()

        }()
    }
}
