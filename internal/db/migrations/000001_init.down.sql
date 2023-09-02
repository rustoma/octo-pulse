-- DropForeignKey
ALTER TABLE "User" DROP CONSTRAINT "User_role_id_fkey";

-- DropForeignKey
ALTER TABLE "Article" DROP CONSTRAINT "Article_author_id_fkey";

-- DropForeignKey
ALTER TABLE "Article" DROP CONSTRAINT "Article_category_id_fkey";

-- DropForeignKey
ALTER TABLE "Article" DROP CONSTRAINT "Article_domain_id_fkey";

-- DropTable
DROP TABLE "User";

-- DropTable
DROP TABLE "Role";

-- DropTable
DROP TABLE "Domain";

-- DropTable
DROP TABLE "Category";

-- DropTable
DROP TABLE "Author";

-- DropTable
DROP TABLE "Article";


