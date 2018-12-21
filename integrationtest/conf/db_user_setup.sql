CREATE ROLE newssearch WITH LOGIN PASSWORD 'password';
GRANT CONNECT ON DATABASE newsranker TO newssearch;
GRANT USAGE ON SCHEMA public TO newssearch;
GRANT SELECT ON article TO newssearch;
GRANT SELECT ON twitter_references TO newssearch;
GRANT SELECT ON subject TO newssearch;
GRANT SELECT ON article_cluster TO newssearch;
GRANT SELECT ON cluster_member TO newssearch;
