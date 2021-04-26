CREATE TABLE schedule.schedule_events (
    ID UUID NOT NULL PRIMARY KEY,
    match_id UUID NOT NULL,
    newer_page VARCHAR,
    older_page VARCHAR,
    start_time TIMESTAMP,
    state VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    block_name VARCHAR NOT NULL,
    league_name VARCHAR NOT NULL,
    league_slug VARCHAR NOT NULL,
    CONSTRAINT fk_match
      FOREIGN KEY(match_id)
          REFERENCES schedule.matches(id)
);