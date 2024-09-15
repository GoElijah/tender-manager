-- +migrate Up
CREATE TYPE bid_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CANCELED',
    'APPROVED'
    );

CREATE TABLE bids
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(50),
    status      bid_status,
    version     INTEGER,
    tender_id   UUID REFERENCES tenders (id),
    tender_organization UUID REFERENCES organization (id),
    bid_organization UUID REFERENCES organization (id),
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR REFERENCES employee (username),
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR REFERENCES employee (username),
    votes INTEGER,
    feedback varchar(50)

);

-- +migrate Down
drop table if exists bids;
