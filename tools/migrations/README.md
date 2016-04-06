# Migrations

A set of scripts or small services to change the MongoDB or Elasticsearch data structures. Simple stuff, e.g. change a default value or add a new column. These scripts avoid that all datasets have to be re-created.

#### Migration 001

Rebalance the crawler schedule and decrease the crawling frequency (from every 12h to every 24h).
