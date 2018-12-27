CREATE TABLE article (
  id VARCHAR(50) PRIMARY KEY,
  url VARCHAR(350) NOT NULL,
  title VARCHAR(350) NOT NULL,
  body TEXT,
  keywords TEXT,
  reference_score NUMERIC(9,5) NOT NULL,
  article_date DATE,
  created_at TIMESTAMP,
  UNIQUE (url)
);

CREATE TABLE twitter_references (
  id VARCHAR(50) PRIMARY KEY,
  twitter_author VARCHAR(50),
  follower_count INT,
  article_id VARCHAR(50) REFERENCES article(id),
  UNIQUE (twitter_author, article_id)
);

CREATE TABLE subject (
  id VARCHAR(50) PRIMARY KEY,
  symbol VARCHAR(15),
  name VARCHAR(50),
  score NUMERIC(9,5) NOT NULL,
  article_id VARCHAR(50) REFERENCES article(id),
  UNIQUE (symbol, article_id)
);

CREATE TABLE article_cluster (
  cluster_hash VARCHAR(64) PRIMARY KEY,
  title VARCHAR(255),
  symbol VARCHAR(15),
  article_date DATE,
  score NUMERIC(9,5),
  lead_article_id VARCHAR(50) REFERENCES article(id)
);

CREATE TABLE cluster_member (
  id VARCHAR(50) PRIMARY KEY,
  reference_score NUMERIC(9,5),
  subject_score NUMERIC(9,5),
  cluster_hash VARCHAR(64) REFERENCES article_cluster(cluster_hash),
  article_id VARCHAR(50) REFERENCES article(id),
  UNIQUE (cluster_hash, article_id)
);

-- Inital data.
INSERT INTO article(id, url, title, keywords, reference_score, article_date)
VALUES 
  ('9dd44a5c-41be-47a8-9bec-cd6398541f1f', 'https://url1.com', 'Title 1', 'w11,w12,w13', 0.15, '2018-12-27'),
  ('cef50e4b-e08c-444e-8ab8-a7092102a62f', 'https://url2.com', 'Title 2', 'w21', 0.2, '2018-12-27'),
  ('b0885e4e-b54a-471d-b34d-d83433d77084', 'https://url3.com', 'Title 3', NULL, 0.1, '2018-12-27');

INSERT INTO article_cluster(cluster_hash, score, lead_article_id, article_date, symbol)
VALUES 
  ('74ee92a6fa51f4ccb156e7ac2e1feeb6979c152ccd00655dd812a17f96afdfc3', 0.5, 'cef50e4b-e08c-444e-8ab8-a7092102a62f', '2018-12-27', 'AAPL'),
  ('9cdc0bf117ce5a4b8f6bbcf2ae08932101af7dc1551688dc7878d48038b090f5', 0.4, '9dd44a5c-41be-47a8-9bec-cd6398541f1f', '2018-12-27', 'GOOG'),
  ('380be4c456e82a5239268a48151c6dfec155fd8d2c8cd2aca23dfa6d757d653a', 0.3, 'b0885e4e-b54a-471d-b34d-d83433d77084', '2018-12-27', 'AAPL');