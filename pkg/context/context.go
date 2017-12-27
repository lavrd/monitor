package context

var context ctx

type ctx struct {
	Verbose bool
	Address string
}

// Get returns context obj
func Get() *ctx {
	return &context
}
