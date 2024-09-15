-- +migrate Up

CREATE TABLE tenders_snapshots
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id   UUID REFERENCES tenders (id),
    name        VARCHAR(50),
    description VARCHAR(50),
    version     INTEGER,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bids_snapshots
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id      UUID REFERENCES bids (id),
    name        VARCHAR(50),
    description VARCHAR(50),
    version     INTEGER,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);


-- +migrate Down
drop table if exists tenders_snapshots;
drop table if exists bids_snapshots;