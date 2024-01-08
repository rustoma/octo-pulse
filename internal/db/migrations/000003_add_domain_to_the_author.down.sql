-- DropForeignKey
ALTER TABLE public.author DROP CONSTRAINT "author_domain_fkey";

-- AlterTable
ALTER TABLE public.author DROP COLUMN "domain";