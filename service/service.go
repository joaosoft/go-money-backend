package goaccount

// NewGoAccount ...
func NewGoAccount(options ...GoAccountOption) *GoAccount {
	account := &GoAccount{}
	account.Reconfigure(options...)

	return account
}
