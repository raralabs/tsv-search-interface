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
    foreign_field  TEXT NOT NULL,
    mapping_field Text NOT NULL,
    PRIMARY KEY(table_info,related_table,foreign_field,mapping_field)
);

CREATE INDEX IF NOT EXISTS "index_related_table_info" on public.related_infos USING btree (table_info);
CREATE INDEX IF NOT EXISTS "index_related_table" on public.related_infos USING btree (related_table);

CREATE OR REPLACE FUNCTION update_changes()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$$
DECLARE t record;
        r record;
BEGIN
    for t in SELECT * FROM "public".related_infos WHERE related_table = NEW.table_info loop
            for r in select * from "public".internal_search_indices where search_field::jsonb ->> t.foreign_field::Text = OLD.id loop
                    UPDATE "public".internal_search_indices SET tsv_text=concat(jsonb_to_tsvector(r.search_field,'["all"]')::text,' ',jsonb_to_tsvector(NEW.search_field,'["all"]')::text)::tsvector
                    WHERE id = r.id;
                end loop;
        end loop;
    RETURN NEW;
END
$$;

DROP TRIGGER update_data_trigger ON "public".internal_search_indices;
CREATE TRIGGER update_data_trigger AFTER UPDATE ON "public".internal_search_indices
    FOR EACH ROW
    WHEN (pg_trigger_depth() < 1)
EXECUTE PROCEDURE update_changes();



CREATE SCHEMA IF NOT EXISTS myra;
CREATE TABLE IF NOT EXISTS myra.search_indices (LIKE public.search_indices INCLUDING ALL);
CREATE TABLE IF NOT EXISTS myra.internal_search_indices (LIKE public.internal_search_indices INCLUDING ALL);
CREATE TABLE IF NOT EXISTS myra.related_infos (LIKE public.related_infos INCLUDING ALL);

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