CREATE KEYSPACE IF NOT EXISTS application
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

USE application;

CREATE TABLE IF NOT EXISTS short_urls (
    code TEXT PRIMARY KEY,
    original_url TEXT,
    created_at TIMESTAMP
);