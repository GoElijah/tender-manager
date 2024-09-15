package entity

type Employee struct {
	ID             string `pg:"id"`
	Username       string `pg:"username"`
	FirstName      string `pg:"first_name"`
	LastName       string `pg:"last_name"`
	CreatedAt      string `pg:"created_at"`
	UpdatedAt      string `pg:"updated_at"`
	OrganizationID string `pg:"organization_id"`
}

type Organization struct {
	ID   string `pg:"id"`
	Name string `pg:"name"`
}

type Tender struct {
	ID           string `pg:"id"`
	Name         string `pg:"name"`
	Description  string `pg:"description, omitempty"`
	Status       string `pg:"status"`
	ServiceType  string `pg:"serviceType"`
	Version      string `pg:"version"`
	CreatedAt    string `pg:"created_at"`
	UpdatedAt    string `pg:"updated_at"`
	CreatedBy    string `pg:"created_at"`
	UpdatedBy    string `pg:"updated_at"`
	Organization string `pg:"organization"`
}

type Bid struct {
	ID                 string `pg:"id"`
	Name               string `pg:"name"`
	Description        string `pg:"description"`
	Status             string `pg:"status"`
	TenderOrganization string `pg:"tender_organization"`
	BidOrganization    string `pg:"bid_organization"`
	TenderID           string `pg:"tender_id"`
	Version            string `pg:"version"`
	Votes              int    `pg:"votes"`
	CreatedAt          string `pg:"created_at"`
	UpdatedAt          string `pg:"updated_at"`
	CreatedBy          string `pg:"created_at"`
	UpdatedBy          string `pg:"updated_at"`
	Feedback           string
}
