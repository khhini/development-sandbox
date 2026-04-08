package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/khhini/development-sandbox/golang/vector-search-quickstart/internals/domains"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
	"google.golang.org/genai"
)

var embeddingOutputDimention = int32(768)

var foodSeeds = []domains.FoodItem{
	{"Sate Ayam", "Potongan daging ayam yang ditusuk dan dibakar, disajikan dengan bumbu kacang gurih dan kecap manis.", "Jawa"},
	{"Gado-Gado", "Campuran sayuran rebus, tahu, tempe, dan telur yang disiram dengan saus kacang kental.", "Jawa"},
	{"Rendang", "Daging sapi yang dimasak perlahan dengan santan dan rempah-rempah hingga bumbu meresap dan menghitam.", "Sumatera"},
	{"Soto Betawi", "Sup daging sapi dengan kuah santan dan susu yang gurih, biasanya disajikan dengan emping.", "Jakarta"},
	{"Coto Makassar", "Sup daging dan jeroan sapi dengan kuah kental berempah, dimakan bersama ketupat atau buras.", "Sulawesi"},
	{"Pempek", "Olahan ikan dan sagu yang digoreng, disajikan dengan kuah cuka (cuko) yang asam, manis, dan pedas.", "Sumatera"},
	{"Soto Ayam Lamongan", "Soto ayam kuning dengan ciri khas bubuk koya gurih dan kuah bening aromatik.", "Jawa Timur"},
	{"Soto Banjar", "Soto ayam khas Kalimantan dengan kuah bening berempah kayu manis dan cengkeh, disajikan dengan perkedel.", "Kalimantan Selatan"},
	{"Soto Padang", "Soto daging sapi goreng (dendeng) dengan kuah bening kaldu sapi dan soun, disajikan dengan kerupuk merah.", "Sumatera Barat"},
	{"Pempek Kapal Selam", "Adonan ikan dan sagu kenyal berisi telur utuh, disajikan dengan kuah cuko hitam pedas manis.", "Sumatera Selatan"},
	{"Batagor", "Bakso tahu goreng renyah disiram saus kacang pedas dan sedikit perasan jeruk limau.", "Jawa Barat"},
	{"Martabak Manis", "Pancake tebal berpori dengan topping mentega, cokelat, kacang, dan keju yang lumer.", "Nasional"},
	{"Bakwan Sayur", "Gorengan tepung renyah berisi irisan kubis dan wortel, biasanya dimakan dengan cabai rawit.", "Nasional"},
	{"Cilok", "Aci dicolok; camilan kenyal dari tepung tapioka dengan bumbu kacang atau saus sambal.", "Jawa Barat"},
	{"Nasi Goreng Kampung", "Nasi yang ditumis dengan terasi, cabai, dan telur, memberikan aroma asap yang kuat.", "Nasional"},
	{"Ayam Taliwang", "Ayam bakar pedas dengan bumbu cabai kering dan terasi, khas suku Sasak.", "NTB"},
	{"Papeda", "Bubur sagu putih yang kental dan lengket, biasanya disandingkan dengan ikan kuah kuning.", "Papua & Maluku"},
	{"Gudeg", "Nangka muda yang dimasak lama dengan santan dan gula merah hingga berwarna cokelat manis.", "Yogyakarta"},
	{"Babi Guling", "Babi panggang utuh dengan bumbu base genep yang kaya rempah dan kulit yang renyah.", "Bali"},
}

func seedDatabase(ctx context.Context, client *genai.Client, conn *pgx.Conn) {
	for _, food := range foodSeeds {
		if err := domains.AddFood(ctx, conn, client, food); err != nil {
			log.Printf("failed to save %s to db: %v", food.Name, err)
		} else {
			fmt.Println("indexing success: %s \n", food.Name)
		}
	}
}

func main() {
	ctx := context.Background()

	client, _ := genai.NewClient(ctx, nil)

	conn, _ := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	defer conn.Close(ctx)
	_ = pgxvector.RegisterTypes(ctx, conn)

	seedDatabase(ctx, client, conn)
}
