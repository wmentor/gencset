BEGIN;

CREATE TABLE locs (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  name character varying(255) NOT NULL,
  code character varying(1024) NOT NULL,
  parent_id INTEGER NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  forms TEXT NOT NULL
);

COMMENT ON TABLE locs IS 'GEO location table';

COMMENT ON COLUMN locs.id IS 'ID';
COMMENT ON COLUMN locs.name IS 'location name';
COMMENT ON COLUMN locs.code IS 'location full name code';
COMMENT ON COLUMN locs.parent_id IS 'parent location id';
COMMENT ON COLUMN locs.latitude IS 'point latitude';
COMMENT ON COLUMN locs.longitude IS 'point longitude';

CREATE INDEX locs_parent_id_idx ON locs(parent_id);
CREATE INDEX locs_name_idx ON locs( LOWER(name) );
CREATE INDEX locs_code_idx ON locs( LOWER(code) );
CREATE INDEX locs_latitude_idx ON locs(latitude);
CREATE INDEX locs_longitude_idx ON locs(longitude);

CREATE FUNCTION on_insert_loc() RETURNS trigger AS $BODY$
BEGIN
  IF NEW.parent_id = 0 THEN
    NEW.code = NEW.name;
  ELSE
    NEW.code = ( SELECT NEW.name || ' / ' || code FROM locs WHERE id = NEW.parent_id LIMIT 1 );
  END IF;
  IF NEW.forms IS NULL THEN
    NEW.forms = NEW.name;
  END IF;
  RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER on_insert_loc_trg BEFORE INSERT ON locs FOR EACH ROW EXECUTE PROCEDURE on_insert_loc();

INSERT INTO locs(name,parent_id,latitude,longitude)
VALUES ('Россия',0,61.698653, 99.505405);

INSERT INTO locs(name,parent_id,latitude,longitude) VALUES
('Москва',(SELECT id FROM locs WHERE name = 'Россия' LIMIT 1),55.753215, 37.622504);

INSERT INTO locs(name,parent_id,latitude,longitude) VALUES
('Красная площадь',(SELECT id FROM locs WHERE name = 'Москва' LIMIT 1),55.753215, 37.622504);

SELECT * FROM locs;

END;
