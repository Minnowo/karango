package pg

const Migration0001 string = `
CREATE TABLE IF NOT EXISTS tbl_version (
    version INTEGER PRIMARY KEY
);
INSERT INTO tbl_version (version) VALUES (0) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS tbl_flags (
    flag INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS tbl_food (
    food_id SERIAL PRIMARY KEY,
    food VARCHAR(255) NOT NULL,
    portion FLOAT  NOT NULL,
    unit VARCHAR(50)  NOT NULL,
    protein FLOAT  NOT NULL,
    carb FLOAT  NOT NULL,
    fibre FLOAT  NOT NULL,
    fat FLOAT  NOT NULL
);

CREATE TABLE IF NOT EXISTS tbl_day (
    day_id          SERIAL PRIMARY KEY,
    day_start       TIMESTAMP NOT NULL,
    day_end         TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tbl_event_master (
    event_id                          SERIAL PRIMARY KEY,
    day_id                            INTEGER NOT NULL REFERENCES tbl_day(day_id),
    event                             VARCHAR(255) NOT NULL,
    event_time                        TIMESTAMP  NOT NULL,
    net_carbs                         FLOAT  NOT NULL,
    blood_glucose                     FLOAT  NOT NULL,
    insulin_sensitivity_factor        FLOAT  NOT NULL,
    insulin_to_carb_ratio             FLOAT  NOT NULL,
    blood_glucose_target              FLOAT  NOT NULL,
    recommended_insulin_amount        FLOAT  NOT NULL,
    actual_insulin_taken              FLOAT  NOT NULL
);

CREATE TABLE IF NOT EXISTS tbl_event_detail (
    event_item_id          SERIAL PRIMARY KEY,
    event_id               INTEGER NOT NULL REFERENCES tbl_event_master(event_id),
    food_id                INTEGER NOT NULL REFERENCES tbl_food(food_id),
    portion                FLOAT  NOT NULL,
    unit                   VARCHAR(255)  NOT NULL,
    protein                FLOAT  NOT NULL,
    carb                   FLOAT  NOT NULL,
    fibre                  FLOAT  NOT NULL,
    fat                    FLOAT  NOT NULL,
    net_carb               FLOAT  NOT NULL
);

CREATE TABLE IF NOT EXISTS tbl_itcr (
    time_period_id SERIAL PRIMARY KEY,
    time_period_start TIME NOT NULL,
    time_period_end TIME NOT NULL,
    insulin_to_carb_ratio FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS ttl_isf (
    period_id SERIAL PRIMARY KEY,
    insulin_sensitivity_factor FLOAT NOT NULL
);

`
