--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dpaola2
--

SELECT pg_catalog.setval('transactions_id_seq', 3, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dpaola2
--

SELECT pg_catalog.setval('users_id_seq', 3, true);


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: dpaola2
--

INSERT INTO transactions (id, owner_id, amount, description) VALUES (1, 1, 5000, 'dinner with friends');
INSERT INTO transactions (id, owner_id, amount, description) VALUES (2, 1, 10000, 'vegas vacation');
INSERT INTO transactions (id, owner_id, amount, description) VALUES (3, 1, 0, 'transaction via facebook id, woot!');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: dpaola2
--

INSERT INTO users (id, facebook_id) VALUES (1, 1932106);
INSERT INTO users (id, facebook_id) VALUES (3, 2000);


--
-- PostgreSQL database dump complete
--

