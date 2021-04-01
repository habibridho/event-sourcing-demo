package repository

type InsufficientBalance struct{}

func (i InsufficientBalance) Error() string {
	return "insufficient balance"
}
