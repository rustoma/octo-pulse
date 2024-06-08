DROP INDEX IF EXISTS category_category_name_key;
CREATE UNIQUE INDEX category_domain_id_name_key ON public.category(domain_id, name);

DROP INDEX IF EXISTS category_category_slug_key;
CREATE UNIQUE INDEX category_domain_id_slug_key ON public.category(domain_id, slug);
