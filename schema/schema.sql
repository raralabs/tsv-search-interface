CREATE TABLE IF NOT EXISTS search_indices(
    id  VARCHAR(36),
    table_info TEXT NOT NULL,
    action_info JSON NOT NULL,
    tsv_text tsvector,
    search_field JSONB,
    PRIMARY KEY(id, table_info)
);
CREATE INDEX IF NOT EXISTS "index_table_info" on public.search_indices USING btree (table_info);

CREATE TABLE IF NOT EXISTS internal_search_indices(
    id  VARCHAR(36),
    table_info TEXT NOT NULL,
    tsv_text tsvector,
    search_field JSONB,
    PRIMARY KEY(id, table_info)
);
CREATE INDEX IF NOT EXISTS "index_internal_table_info" on public.internal_search_indices USING btree (table_info);

CREATE TABLE IF NOT EXISTS related_infos(
    table_info TEXT NOT NULL,
    related_table TEXT NOT NULL,
    foreign_field  TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS "index_related_table_info" on public.related_infos USING btree (table_info);
CREATE INDEX IF NOT EXISTS "index_related_table" on public.related_infos USING btree (related_table);

CREATE SCHEMA IF NOT EXISTS myra;
CREATE TABLE IF NOT EXISTS myra.search_indices (LIKE public.search_indices INCLUDING ALL);
CREATE TABLE IF NOT EXISTS myra.internal_search_indices (LIKE public.internal_search_indices INCLUDING ALL);
CREATE TABLE IF NOT EXISTS myra.related_info (LIKE public.related_infos INCLUDING ALL);

-- CREATE SCHEMA IF NOT EXISTS rara2;
-- CREATE TABLE IF NOT EXISTS rara2.search_indices (LIKE public.search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara2.internal_search_indices (LIKE public.internal_search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara2.related_info (LIKE public.related_info INCLUDING ALL);
--
-- CREATE SCHEMA IF NOT EXISTS rara3;
-- CREATE TABLE IF NOT EXISTS rara3.search_indices (LIKE public.search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara3.internal_search_indices (LIKE public.internal_search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara3.related_info (LIKE public.related_info INCLUDING ALL);
--
-- CREATE SCHEMA IF NOT EXISTS rara4;
-- CREATE TABLE IF NOT EXISTS rara4.search_indices (LIKE public.search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara4.internal_search_indices (LIKE public.internal_search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara4.related_info (LIKE public.related_info INCLUDING ALL);
--
-- CREATE SCHEMA IF NOT EXISTS rara5;
-- CREATE TABLE IF NOT EXISTS rara5.search_indices (LIKE public.search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara5.internal_search_indices (LIKE public.internal_search_indices INCLUDING ALL);
-- CREATE TABLE IF NOT EXISTS rara5.related_info (LIKE public.related_info INCLUDING ALL);