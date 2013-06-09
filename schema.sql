--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: dpaola2; Tablespace: 
--

CREATE TABLE transactions (
    id integer NOT NULL,
    owner_id integer NOT NULL,
    amount integer NOT NULL,
    description text
);


ALTER TABLE public.transactions OWNER TO dpaola2;

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: dpaola2
--

CREATE SEQUENCE transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.transactions_id_seq OWNER TO dpaola2;

--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dpaola2
--

ALTER SEQUENCE transactions_id_seq OWNED BY transactions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: dpaola2; Tablespace: 
--

CREATE TABLE users (
    id integer NOT NULL,
    facebook_id integer NOT NULL
);


ALTER TABLE public.users OWNER TO dpaola2;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: dpaola2
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO dpaola2;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dpaola2
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: dpaola2
--

ALTER TABLE ONLY transactions ALTER COLUMN id SET DEFAULT nextval('transactions_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: dpaola2
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: users_facebook_id_key; Type: CONSTRAINT; Schema: public; Owner: dpaola2; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_facebook_id_key UNIQUE (facebook_id);


--
-- Name: public; Type: ACL; Schema: -; Owner: dpaola2
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM dpaola2;
GRANT ALL ON SCHEMA public TO dpaola2;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

