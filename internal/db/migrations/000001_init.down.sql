-- DropForeignKey
ALTER TABLE public.user DROP CONSTRAINT "user_role_id_fkey";

-- DropForeignKey
ALTER TABLE public.article DROP CONSTRAINT "article_author_id_fkey";

-- DropForeignKey
ALTER TABLE public.article DROP CONSTRAINT "article_category_id_fkey";

-- DropForeignKey
ALTER TABLE public.article DROP CONSTRAINT "article_domain_id_fkey";

-- DropTable
DROP TABLE public.user;

-- DropTable
DROP TABLE public.role;

-- DropTable
DROP TABLE public.domain;

-- DropTable
DROP TABLE public.category;

-- DropTable
DROP TABLE public.author;

-- DropTable
DROP TABLE public.article;


