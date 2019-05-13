--
-- PostgreSQL database dump
--

-- Dumped from database version 11.2 (Debian 11.2-1.pgdg90+1)
-- Dumped by pg_dump version 11.2 (Debian 11.2-1.pgdg90+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: sensor; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sensor (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    serial_no character varying(50) NOT NULL,
    unit_type character varying(50) NOT NULL,
    max_safe_value double precision NOT NULL,
    min_safe_value double precision NOT NULL
);


--
-- Name: sensor_reading; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sensor_reading (
    id bigint NOT NULL,
    value double precision NOT NULL,
    sensor_id integer NOT NULL,
    taken_on timestamp(6) without time zone NOT NULL
);


--
-- Name: sensor_reading_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.sensor_reading ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.sensor_reading_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: sensor sensor_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sensor
    ADD CONSTRAINT sensor_pkey PRIMARY KEY (id);


--
-- Name: sensor_reading sensor_reading_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sensor_reading
    ADD CONSTRAINT sensor_reading_pkey PRIMARY KEY (id);


--
-- Name: fki_fk_sensor_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_sensor_id ON public.sensor_reading USING btree (sensor_id);


--
-- Name: sensor_reading fk_sensor_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sensor_reading
    ADD CONSTRAINT fk_sensor_id FOREIGN KEY (sensor_id) REFERENCES public.sensor(id);


--
-- PostgreSQL database dump complete
--

