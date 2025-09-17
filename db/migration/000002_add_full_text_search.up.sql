-- Add a new column to store the searchable vector
ALTER TABLE articles ADD COLUMN content_tsv tsvector;

-- Create a function to update the tsvector column from the title and content
CREATE OR REPLACE FUNCTION articles_tsvector_update() RETURNS trigger AS $$
BEGIN
    NEW.content_tsv :=
        setweight(to_tsvector('simple', coalesce(NEW.title,'')), 'A') ||
        setweight(to_tsvector('simple', coalesce(NEW.content,'')), 'B');
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically update the tsvector column when a row is inserted or updated
CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
ON articles FOR EACH ROW EXECUTE FUNCTION articles_tsvector_update();

-- Create a GIN index for performance on the tsvector column
CREATE INDEX content_tsv_idx ON articles USING GIN (content_tsv);

-- To populate the new column for existing data, run this UPDATE statement
UPDATE articles SET content_tsv = 
    setweight(to_tsvector('simple', coalesce(title,'')), 'A') ||
    setweight(to_tsvector('simple', coalesce(content,'')), 'B');