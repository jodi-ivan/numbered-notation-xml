package canvas

type Writer interface {
	Write(p []byte) (n int, err error)
}
