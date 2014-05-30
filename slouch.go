package slouch

import (
  "net/http"
  "strings"
  "time"
  "fmt"
)

type Conn interface {
  ReadMessage() (string, error)
  WriteMessage(m string) error
  Dial() error
  Close() error
}

type StompClient struct {
  Conn Conn
  On bool
  FrameChan chan string
}

func (s *StompClient) Dial() error {
  return s.Conn.Dial()
}

func (s *StompClient) ReadMessage() (string, error) {
  return s.Conn.ReadMessage()
}

func (s *StompClient) WriteMessage(m string) error {
  return s.Conn.WriteMessage(m)
}

func (s *StompClient) Connect(headers http.Header) error {
  headerString := ""
  for k, v := range headers {
    headerString += strings.ToLower(k) + ":" + strings.Join(v, ",") + "\n"
  }

  var connect = "CONNECT\n" +
  headerString +
  "\n" +
  "\u0000" +
  "\n"

  return s.Conn.WriteMessage(connect)
}

func (s *StompClient) SendMessage(message string) {
  s.FrameChan <- message
}

func (s *StompClient) StartHeartbeat() {
  go func() {
    for s.On {
      time.Sleep(5 * time.Second)
      s.FrameChan <- "\n"
    }
  }()
}

func (s *StompClient) StartMessageHandler() {
  go func() {
    for s.On {
      msg, err := s.ReadMessage()
      if err != nil {
        s.On = false
        s.Conn.Close()
        return
      }
      fmt.Println(msg)
    }
  }()
}

func (s *StompClient) StartSender() {
  go func() {
    for s.On {
      frame := <-s.FrameChan
      err := s.WriteMessage(frame)
      if err != nil {
        s.On = false
        s.Conn.Close()
        return
      }
    }
  }()
}
