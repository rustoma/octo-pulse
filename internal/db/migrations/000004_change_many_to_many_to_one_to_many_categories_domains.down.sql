CREATE TABLE categories_domains (
    domain_id INT NOT NULL,
    category_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

ALTER TABLE public.category ALTER COLUMN domain_id DROP NOT NULL;
UPDATE public.category SET domain_id = NULL;

ALTER TABLE public.category DROP COLUMN domain_id;






