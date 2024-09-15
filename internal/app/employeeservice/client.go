package employeeservice

type Client struct {
	storage employeeStorage
}

func New(s employeeStorage) *Client {
	return &Client{
		storage: s,
	}
}
