package demo

func ChanCloseDemo() {
	ch := make(chan any)

	close(ch)
	// close(ch) // 会报错。
}
