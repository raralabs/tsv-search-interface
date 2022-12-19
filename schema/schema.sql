CREATE TABLE IF NOT EXISTS search_indices(
    id  VARCHAR(36),
    table_info TEXT NOT NULL,
    action_info JSON NOT NULL,
    tsv_text tsvector,
    search_field JSONB,
    PRIMARY KEY(id, table_info)
);
CREATE INDEX "index_table_info" on public.search_indices USING btree (table_info);

CREATE SCHEMA rara1;
CREATE TABLE rara1.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara2;
CREATE TABLE rara2.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara3;
CREATE TABLE rara3.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara4;
CREATE TABLE rara4.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara5;
CREATE TABLE rara5.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara6;
CREATE TABLE rara6.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara7;
CREATE TABLE rara7.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara8;
CREATE TABLE rara8.search_indices (LIKE public.search_indices INCLUDING ALL);

CREATE SCHEMA rara9;
CREATE TABLE rara9.search_indices (LIKE public.search_indices INCLUDING ALL);