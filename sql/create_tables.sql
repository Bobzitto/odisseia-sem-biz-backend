--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

-- 
-- Name: turmas; type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.turmas (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    school character varying(255) NOT NULL,
    year character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NULL
);

--
-- 
--
ALTER TABLE public.turmas ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.turmas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

-- 
-- Name: aulas; type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.aulas (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    size character varying(255) NOT NULL,
    active boolean NOT NULL,
    review numeric NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NULL
);

COPY public.aulas (id, name, size, active, review, created_at, updated_at) FROM stdin;
1	Regra de 3 Simples	4	true	4.5	2022-09-23 00:00:00	2022-09-23 00:00:00 
2	Termodinâmica	7	true	5	2022-09-23 00:00:00	2022-09-23 00:00:00
3	Sinônimos	5	true	4	2022-09-23 00:00:00	2022-09-23 00:00:00
4	Segunda Guerra Mundial	5	true	2.8	2022-09-23 00:00:00	2022-09-23 00:00:00
5	Tigres Asiáticas	4	false	3.8	2022-09-23 00:00:00	2022-09-23 00:00:00
6	Verb To Be	7	true	4.3	2022-09-23 00:00:00	2022-09-23 00:00:00
\.

--
-- 
--
ALTER TABLE public.aulas ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.aulas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY public.aulas
    ADD CONSTRAINT aulas_pkey PRIMARY KEY (id);
--
-- Name: materias; type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.materias (
    id integer NOT NULL,
    materia character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

--
-- Fix duplicate id issue: remove or change the duplicate entry
-- You can use the following command to remove the duplicate with id = 8
-- DELETE FROM public.materias WHERE id = 8;

-- 
-- Ensure the following COPY command has unique IDs
COPY public.materias (id, materia, created_at, updated_at) FROM stdin;
1	Física	2022-09-23 00:00:00	2022-09-23 00:00:00
2	Matemática	2022-09-23 00:00:00	2022-09-23 00:00:00
3	Português	2022-09-23 00:00:00	2022-09-23 00:00:00
4	História	2022-09-23 00:00:00	2022-09-23 00:00:00
5	Geografia	2022-09-23 00:00:00	2022-09-23 00:00:00
6	Sociologia	2022-09-23 00:00:00	2022-09-23 00:00:00
7	Inglês	2022-09-23 00:00:00	2022-09-23 00:00:00
8	Filosofia	2022-09-23 00:00:00	2022-09-23 00:00:00
9	Artes	2022-09-23 00:00:00	2022-09-23 00:00:00
\.

ALTER TABLE ONLY public.materias
    ADD CONSTRAINT materias_pkey PRIMARY KEY (id);

CREATE TABLE public.aulas_materias (
    id integer NOT NULL,
    aula_id integer,
    materia_id integer
);

COPY public.aulas_materias (id, aula_id, materia_id) FROM stdin;
1	1	2
2	2	1
3	3	3
4	4	4
5	5	5
6	6	7
\.

ALTER TABLE public.aulas_materias ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.aulas_materias_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY public.aulas_materias
    ADD CONSTRAINT aulas_materias_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.aulas_materias
    ADD CONSTRAINT aulas_materias_materia_id_fkey FOREIGN KEY (materia_id) REFERENCES public.materias(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.aulas_materias
    ADD CONSTRAINT aulas_materias_aula_id_fkey FOREIGN KEY (aula_id) REFERENCES public.aulas(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- 
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

--
-- 
--
ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- 
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, first_name, last_name, email, password, created_at, updated_at) FROM stdin;
1	Admin	User	admin@example.com	$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOIEHpyh2DAMmso1wy	2022-09-23 00:00:00	2022-09-23 00:00:00
\.

-- 
-- Data for Name: turmas; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.turmas (id, name, school, year, created_at, updated_at) FROM stdin;
1	405	E.M. Porra Flávio	4o ano	2022-09-23 00:00:00	2022-09-23 00:00:00
2	712	E.M. Porra Flávio	7o ano	2022-09-23 00:00:00	2022-09-23 00:00:00
3	5o ano do barulho	CEPT Treinamento Presencial	5o ano	2022-09-23 00:00:00	2022-09-23 00:00:00
4	Turminha 201	CEPT Treinamento Presencial	2o ano	2022-09-23 00:00:00	2022-09-23 00:00:00
5	9D	CEPT Treinamento Presencial	9o ano	2022-09-23 00:00:00	2022-09-23 00:00:00
6	313	Escola Implantação	3o ano	2022-09-23 00:00:00	2022-09-23 00:00:00
\.

