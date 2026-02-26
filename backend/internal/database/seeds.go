package database

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"

// 	"github.com/Samu-Amy/Shokora/internal/store"
// )

// var first_names = []string{"Michael", "Gianni", "Anna", "Sam", "Maria", "Franco", "Matteo", "Bobby"}
// var last_names = []string{"Gibbs", "Ruiz", "Rivera", "Lutz", "Tanner", "Coleman", "Villa"}
// var email_suffix = []string{"example", "gmail", "hotmail"}
// var email_domain = []string{"com", "it"}

// var product_names = []string{
// 	"Tiramisù Classico",
// 	"Cannolo Siciliano",
// 	"Babà al Rum",
// 	"Crostatina ai Frutti di Bosco",
// 	"Millefoglie alla Crema",
// 	"Sfogliatella Riccia",
// 	"Zeppola di San Giuseppe",
// 	"Panna Cotta alla Vaniglia",
// 	"Cheesecake al Limone",
// 	"Torta Caprese",
// 	"Delizia al Limone",
// 	"Tortino al Cuore Caldo",
// 	"Croissant al Burro",
// 	"Croissant alla Crema",
// 	"Croissant al Pistacchio",
// 	"Plumcake allo Yogurt",
// 	"Pralina Fondente 70%",
// 	"Pralina al Pistacchio",
// 	"Pralina alla Nocciola",
// 	"Gianduia Artigianale",
// 	"Tavoletta Extra Fondente",
// 	"Tavoletta al Latte",
// 	"Tavoletta Bianca al Cocco",
// 	"Cioccolatino al Rum",
// 	"Cuore di Cioccolato",
// 	"Scorze d'Arancia Candite",
// 	"Crema Spalmabile al Cacao",
// 	"Dragées alle Mandorle",
// 	"Caffè Espresso",
// 	"Cappuccino Classico",
// 	"Cappuccino al Cacao",
// 	"Latte Macchiato",
// 	"Caffè Marocchino",
// 	"Caffè Shakerato",
// 	"Cioccolata Calda Fondente",
// 	"Cioccolata Calda Bianca",
// 	"The al Limone",
// 	"The alla Pesca",
// 	"Spremuta d'Arancia",
// 	"Succo di Frutta Artigianale",
// }

// var product_descriptions = []string{
// 	"Dolce tradizionale con savoiardi imbevuti di caffè e crema al mascarpone.",
// 	"Croccante cialda fritta ripiena di ricotta dolce.",
// 	"Dolce soffice lievitato, imbevuto di rum.",
// 	"Base di pasta frolla con crema e frutti di bosco freschi.",
// 	"Strati di sfoglia croccante farciti con crema pasticcera.",
// 	"Dolce napoletano di pasta sfoglia ripieno di ricotta e semolino.",
// 	"Frittella farcita con crema pasticcera.",
// 	"Crema delicata alla vaniglia con coulis leggero.",
// 	"Torta fredda con crema al formaggio e limone.",
// 	"Torta al cioccolato fondente e mandorle.",
// 	"Dessert soffice al profumo di limone della Costiera.",
// 	"Piccolo tortino al cioccolato con cuore morbido.",
// 	"Croissant fragrante preparato con burro di alta qualità.",
// 	"Croissant ripieno di crema pasticcera.",
// 	"Croissant farcito con crema al pistacchio.",
// 	"",
// 	"Cioccolatino fondente dal gusto intenso.",
// 	"Pralina ripiena di crema al pistacchio.",
// 	"Pralina ripiena di crema alla nocciola.",
// 	"Classico impasto di cioccolato e nocciole.",
// 	"Tavoletta di cioccolato fondente ad alta percentuale di cacao.",
// 	"Tavoletta di cioccolato al latte morbida e vellutata.",
// 	"Tavoletta di cioccolato bianco aromatizzata al cocco.",
// 	"Cioccolatino aromatizzato al rum.",
// 	"",
// 	"Scorze d’arancia ricoperte di cioccolato fondente.",
// 	"Crema spalmabile al cacao ideale per colazioni e dessert.",
// 	"Mandorle ricoperte di cioccolato croccante.",
// 	"Caffè espresso italiano dal gusto intenso.",
// 	"Cappuccino con latte montato e schiuma cremosa.",
// 	"Cappuccino arricchito con cacao in polvere.",
// 	"Latte caldo con caffè espresso.",
// 	"Caffè con schiuma di latte e cacao.",
// 	"Caffè freddo shakerato con ghiaccio.",
// 	"Bevanda calda al cioccolato fondente.",
// 	"Bevanda calda a base di cioccolato bianco.",
// 	"",
// 	"",
// 	"Spremuta fresca di arance selezionate.",
// 	"Succo ottenuto da frutta selezionata.",
// }

// func Seed(store store.Storage, db *sql.DB) {
// 	ctx := context.Background()

// 	// TODO: aggiorna (per le tabelle aggiunte dopo)

// 	// Users
// 	users := generateUsers(100)
// 	transaction, _ := db.BeginTx(ctx, nil)

// 	for _, user := range users {
// 		if err := store.User.Create(ctx, user); err != nil {
// 			_ = transaction.Rollback()
// 			log.Println("Error creating user: ", err)
// 			return
// 		}
// 	}

// 	transaction.Commit()

// 	// Products
// 	products := generateProducts(40)

// 	for _, product := range products {
// 		if err := store.Product.Create(ctx, product); err != nil {
// 			log.Println("Error creating product: ", err)
// 			return
// 		}
// 	}
// }

// func generateUsers(num int) []*store.User {
// 	users := make([]*store.User, num)

// 	for i := 0; i < num; i++ {
// 		first_name := first_names[i%len(first_names)] + fmt.Sprintf("_%d", i)
// 		last_name := last_names[i%len(last_names)] + fmt.Sprintf("_%d", i)

// 		users[i] = &store.User{
// 			FirstName: first_name,
// 			LastName:  last_name,
// 			Email:     fmt.Sprintf("%s@%s.%s", first_name, email_suffix[i%len(email_suffix)], email_domain[i%len(email_domain)]),
// 		}
// 	}

// 	return users
// }

// func generateProducts(num int) []*store.Product {
// 	product := make([]*store.Product, num)

// 	for i := 0; i < num; i++ {
// 		product_name := product_names[i%len(product_names)]
// 		product_description := product_descriptions[i%len(product_descriptions)]

// 		product[i] = &store.Product{
// 			Name:        product_name,
// 			Description: product_description,
// 			ImageURL:    "",
// 			Price:       rand.Float64() * 25,
// 			Discount:    rand.Float64(),
// 		}
// 	}

// 	return product
// }

// // func generateMenuOrders(num int, products []*store.Product, users []*store.User) []*store.User {
// // 	orders := make([]*store.MenuOrder, num)

// // 	for i := 0; i < num; i++ {
// // 		// Get random product and user
// // 		product := products[rand.Intn(len(products))]
// // 		user := users[rand.Intn(len(users))]

// // 		menuOrder[i] = &store.User{
// // 			//...
// // 		}
// // 	}

// // 	return orders
// // }
