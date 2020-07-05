BEGIN;

CREATE TABLE locs (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  name character varying(255) NOT NULL,
  code character varying(1024) NOT NULL,
  parent_id INTEGER,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL
);

COMMENT ON TABLE locs IS 'GEO location table';

COMMENT ON COLUMN locs.id IS 'ID';
COMMENT ON COLUMN locs.name IS 'location name';
COMMENT ON COLUMN locs.code IS 'location full name code';
COMMENT ON COLUMN locs.parent_id IS 'parent location id';
COMMENT ON COLUMN locs.latitude IS 'point latitude';
COMMENT ON COLUMN locs.longitude IS 'point longitude';

CREATE INDEX locs_parent_id_idx ON locs(parent_id);
CREATE INDEX logcs_name_idx ON locs( LOWER(name) );
CREATE INDEX logcs_code_idx ON locs( LOWER(code) );
CREATE INDEX locs_latitude_idx ON locs(latitude);
CREATE INDEX locs_longitude_idx ON locs(longitude);

CREATE TABLE forms (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    loc_id BIGSERIAL NOT NULL REFERENCES locs(id) ON DELETE CASCADE,
    name character varying(255) NOT NULL
);

COMMENT ON TABLE forms IS 'Location name forms';

COMMENT ON COLUMN forms.id IS 'Location form id';
COMMENT ON COLUMN forms.loc_id IS 'Base location id';
COMMENT ON COLUMN forms.name IS 'Location form name';

CREATE INDEX forms_loc_id_idx ON forms(loc_id);
CREATE INDEX forms_loc_name_idx ON forms(LOWER(name));

END;
