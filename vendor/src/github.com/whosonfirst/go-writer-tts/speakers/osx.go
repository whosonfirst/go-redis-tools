package speakers

import (
       "github.com/everdev/mack"
       "strings"       
)

type OSXSpeaker struct {
     Speaker
}

func NewOSXSpeaker () (*OSXSpeaker, error) {

     s := OSXSpeaker{}
     return &s, nil
}

func (s *OSXSpeaker) Speak (text string) error {

     _, err := s.WriteString(text)
     return err
}

func (s *OSXSpeaker) WriteString (text string) (int64, error) {
     r := strings.NewReader(text)
     return r.WriteTo(s)
}

func (s *OSXSpeaker) Write(p []byte) (int, error) {

     var text string
     text = string(p[:])
    
     mack.Say(text)

     count := len(text)
     return count, nil
}

func (s *OSXSpeaker) Close() error {
     return nil
}
