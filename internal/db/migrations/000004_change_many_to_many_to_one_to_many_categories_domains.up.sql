ALTER TABLE public.category ADD COLUMN domain_id INT;

-- Copy data from categories_domains to categories
UPDATE public.category
SET domain_id = cd.domain_id
FROM public.categories_domains cd
WHERE public.category.id = cd.category_id;

-- Set domain_id to NOT NULL
ALTER TABLE public.category ALTER COLUMN domain_id SET NOT NULL;

ALTER TABLE public.category
ADD CONSTRAINT fk_domain
FOREIGN KEY (domain_id) REFERENCES public.domain(id);


DROP TABLE IF EXISTS categories_domains;
