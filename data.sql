DROP TABLE IF EXISTS station CASCADE;
DROP TABLE IF EXISTS route CASCADE;
DROP TABLE IF EXISTS route_stations CASCADE;
DROP TABLE IF EXISTS route_stations_time CASCADE;

CREATE TABLE public.station (
    id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar (100)
);
CREATE TABLE public.route (
    id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar (100)
);
CREATE TABLE public.route_stations (
    id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    route_id INTEGER,
    station_id INTEGER,
    pos INTEGER,
    CONSTRAINT route_id_fk FOREIGN KEY (route_id) REFERENCES public.route(id),
    CONSTRAINT station_id_fk FOREIGN KEY (station_id) REFERENCES public.station(id),
    CONSTRAINT route_station_unique UNIQUE (route_id, station_id)
);
CREATE TABLE public.route_stations_time (
    id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    route_station_id INTEGER,
    queue INTEGER,
    stop_time TIME NOT NULL,
    CONSTRAINT route_station_id_fk FOREIGN KEY (route_station_id) REFERENCES public.route_stations(id),
    CONSTRAINT route_station_queue_unique UNIQUE (route_station_id, queue)
);

ALTER SEQUENCE route_id_seq RESTART WITH 1;
ALTER SEQUENCE station_id_seq RESTART WITH 1;
ALTER SEQUENCE route_stations_id_seq RESTART WITH 1;
ALTER SEQUENCE route_stations_time_id_seq RESTART WITH 1;

INSERT INTO station (name) VALUES ('м. Купчино');
INSERT INTO station (name) VALUES ('м. Московская');
INSERT INTO station (name) VALUES ('м. Технологический институт');
INSERT INTO station (name) VALUES ('Невский проспект 110');
INSERT INTO station (name) VALUES ('Пулково');
INSERT INTO station (name) VALUES ('Пулковское шоссе 10');

INSERT INTO route (name) VALUES ('Автобус № 1 Купчино-Невский');
INSERT INTO route (name) VALUES ('Автобус № 1 Невский-Купчино');
INSERT INTO route (name) VALUES ('Автобус № 2 Пулково-Московская');
INSERT INTO route (name) VALUES ('Автобус № 2 Московская-Пулково');

INSERT INTO route_stations (route_id, station_id, pos) VALUES (1, 1, 0);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (1, 2, 1);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (1, 3, 2);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (1, 4, 3);

INSERT INTO route_stations (route_id, station_id, pos) VALUES (2, 4, 0);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (2, 3, 1);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (2, 2, 2);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (2, 1, 3);

INSERT INTO route_stations (route_id, station_id, pos) VALUES (3, 5, 0);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (3, 6, 1);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (3, 2, 2);

INSERT INTO route_stations (route_id, station_id, pos) VALUES (4, 2, 0);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (4, 6, 1);
INSERT INTO route_stations (route_id, station_id, pos) VALUES (4, 5, 2);

INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (1, 0, '08:00:00');
INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (2, 0, '08:15:00');
INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (3, 0, '08:30:00');
INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (4, 0, '08:45:00');

INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (1, 1, '20:00:00');
INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (2, 1, '20:15:00');
INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (3, 1, '20:30:00');
INSERT INTO route_stations_time (route_station_id, queue, stop_time) VALUES (4, 1, '20:45:00');