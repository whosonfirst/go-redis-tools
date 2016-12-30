package speakers

type Speaker interface {
     WriteString (text string) (int64, error)
     Write(p []byte) (int, error)
     Close() error 
}
