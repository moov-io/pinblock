package formats

type ISO interface {
	Decode(pinBlock, account string) (string, error)
	Encode(pin, account string) (string, error)
	Format() string
	Padding(pin string) string
}

var _ ISO = &ISO0{}
var _ ISO = &ISO4{}
