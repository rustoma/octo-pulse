-- CreateTable
CREATE TABLE IF NOT EXISTS public.user (
    "id" SERIAL NOT NULL,
    "email" TEXT NOT NULL,
    "refresh_token" TEXT,
    "password_hash" TEXT NOT NULL,
    "role_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE IF NOT EXISTS public.role (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "role_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE IF NOT EXISTS public.domain (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "domain_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE IF NOT EXISTS public.category (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "slug" TEXT NOT NULL,
    "weight" INTEGER NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "category_pkey" PRIMARY KEY ("id")
);


-- CreateTable
CREATE TABLE IF NOT EXISTS public.categories_domains (
    "domain_id" INTEGER NOT NULL,
    "category_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "categories_domains_pkey" PRIMARY KEY ("domain_id","category_id")
);

-- CreateTable
CREATE TABLE IF NOT EXISTS public.author (
    "id" SERIAL NOT NULL,
    "first_name" TEXT NOT NULL,
    "last_name" TEXT NOT NULL,
    "description" TEXT,
    "image_url" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "author_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE IF NOT EXISTS public.article (
    "id" SERIAL NOT NULL,
    "title" TEXT NOT NULL,
    "slug" TEXT NOT NULL,
    "body" TEXT,
    "thumbnail" INTEGER,
    "publication_date" TIMESTAMP(3),
    "is_published" BOOLEAN NOT NULL DEFAULT false,
    "author_id" INTEGER,
    "category_id" INTEGER NOT NULL,
    "domain_id" INTEGER NOT NULL,
    "featured" BOOLEAN NOT NULL DEFAULT false,
    "reading_time" INTEGER,
    "is_sponsored" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "article_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE public.image_storage (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "path" TEXT NOT NULL,
    "size" BIGINT NOT NULL,
    "type" TEXT,
    "width" INTEGER,
    "height" INTEGER,
    "alt" TEXT,
    "category_id" INTEGER,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "image_storage_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE public.image_category (
     "id" SERIAL NOT NULL,
     "name" TEXT NOT NULL,
     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
     "updated_at" TIMESTAMP(3) NOT NULL,
     
     CONSTRAINT "image_category_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE public.basic_page (
     "id" SERIAL NOT NULL,
     "title" TEXT NOT NULL,
     "slug" TEXT NOT NULL,
     "body" TEXT,
     "domain" INTEGER NOT NULL,
     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
     "updated_at" TIMESTAMP(3) NOT NULL,

     CONSTRAINT "basic_page_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "article_slug_domain_id_key" ON public.article("slug", "domain_id");

-- CreateIndex
CREATE UNIQUE INDEX "user_email_key" ON public.user("email");

-- CreateIndex
CREATE UNIQUE INDEX "domain_domain_name_key" ON public.domain("name");

-- CreateIndex
CREATE UNIQUE INDEX "category_category_name_key" ON public.category("name");
CREATE UNIQUE INDEX "category_category_slug_key" ON public.category("slug");

-- CreateIndex
CREATE UNIQUE INDEX "basic_page_slug_domain_key" ON public.basic_page("slug", "domain");

-- CreateIndex
CREATE UNIQUE INDEX "imageStorage_path_key" ON public.image_storage("path");

-- AddForeignKey
ALTER TABLE public.user ADD CONSTRAINT "user_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES public.role("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.article ADD CONSTRAINT "article_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES public.author("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.article ADD CONSTRAINT "article_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES public.category("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.article ADD CONSTRAINT "article_domain_id_fkey" FOREIGN KEY ("domain_id") REFERENCES public.domain("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.article ADD CONSTRAINT "article_thumbnail_fkey" FOREIGN KEY ("thumbnail") REFERENCES public.image_storage("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.categories_domains ADD CONSTRAINT "categories_domains_domain_id_fkey" FOREIGN KEY ("domain_id") REFERENCES public.domain("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.categories_domains ADD CONSTRAINT "categories_domains_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES public.category("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.image_storage ADD CONSTRAINT "image_storage_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES public.image_category("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE public.basic_page ADD CONSTRAINT "basic_page_domain_fkey" FOREIGN KEY ("domain") REFERENCES public.domain("id") ON DELETE RESTRICT ON UPDATE CASCADE;

