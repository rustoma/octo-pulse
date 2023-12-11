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
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "domain_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE IF NOT EXISTS public.category (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
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
    "description" TEXT,
    "thumbnail" INTEGER,
    "publication_date" TIMESTAMP(3),
    "is_published" BOOLEAN NOT NULL DEFAULT false,
    "author_id" INTEGER,
    "category_id" INTEGER NOT NULL,
    "domain_id" INTEGER NOT NULL,
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
    "upload_date" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "type" TEXT,
    "width" INTEGER,
    "height" INTEGER,
    "category_id" INTEGER,

    CONSTRAINT "ImageStorage_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE public.image_category (
     "category_id" SERIAL NOT NULL,
     "category_name" TEXT NOT NULL,

     CONSTRAINT "ImageCategory_pkey" PRIMARY KEY ("category_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "user_email_key" ON public.user("email");

-- CreateIndex
CREATE UNIQUE INDEX "domain_domain_name_key" ON public.domain("name");

-- CreateIndex
CREATE UNIQUE INDEX "category_category_name_key" ON public.category("name");

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
ALTER TABLE public.image_storage ADD CONSTRAINT "image_storage_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES public.image_category("category_id") ON DELETE SET NULL ON UPDATE CASCADE;

