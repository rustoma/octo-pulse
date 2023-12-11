-- DropForeignKey
ALTER TABLE public.user DROP CONSTRAINT "user_role_id_fkey";

-- DropForeignKey
ALTER TABLE public.article DROP CONSTRAINT "article_author_id_fkey";

-- DropForeignKey
ALTER TABLE public.article DROP CONSTRAINT "article_category_id_fkey";

-- DropForeignKey
ALTER TABLE public.article DROP CONSTRAINT "article_domain_id_fkey";

-- DropForeignKey
ALTER TABLE public.article DROP CONSTRAINT "article_thumbnail_fkey";

-- DropForeignKey
ALTER TABLE public.categories_domains DROP CONSTRAINT "categories_domains_pkey";

-- DropForeignKey
ALTER TABLE public.categories_domains DROP CONSTRAINT "categories_domains_category_id_fkey";

-- DropForeignKey
ALTER TABLE public.image_storage DROP CONSTRAINT "image_storage_category_id_fkey";

-- DropTable
DROP TABLE public.user;

-- DropTable
DROP TABLE public.role;

-- DropTable
DROP TABLE public.categories_domains;

-- DropTable
DROP TABLE public.domain;

-- DropTable
DROP TABLE public.category;

-- DropTable
DROP TABLE public.author;

-- DropTable
DROP TABLE public.article;

-- DropTable
DROP TABLE public.image_storage;

-- DropTable
DROP TABLE public.image_category;