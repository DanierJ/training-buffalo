--
-- PostgreSQL database dump
--

-- Dumped from database version 10.10
-- Dumped by pg_dump version 10.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: devices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devices (
    id uuid NOT NULL,
    manufacturer character varying(255) NOT NULL,
    make character varying(255) NOT NULL,
    model character varying(255) NOT NULL,
    storage character varying(255) NOT NULL,
    cost integer NOT NULL,
    operating_system character varying(255) NOT NULL,
    image_url character varying(255) NOT NULL,
    is_new boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    user_id uuid
);


ALTER TABLE public.devices OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    lastname character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: devices devices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devices
    ADD CONSTRAINT devices_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

