-- AlterTable
ALTER TABLE public.author ADD COLUMN "domain" INTEGER;

-- AddForeignKey
ALTER TABLE public.author ADD CONSTRAINT "author_domain_fkey" FOREIGN KEY ("domain") REFERENCES public.domain("id") ON DELETE SET NULL ON UPDATE CASCADE;

