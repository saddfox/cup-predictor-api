package db

import (
	"errors"
	"fmt"

	"github.com/saddfox/cup-predictor/pkg/auth"
	"github.com/saddfox/cup-predictor/pkg/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// populates database with test and admin user
func SeedUsers() error {
	var user models.User
	if errors.Is(DB.First(&user).Error, gorm.ErrRecordNotFound) == false {
		return errors.New("DB already contains users")
	}

	pwd1, _ := auth.HashPassword("admin")
	user1 := models.User{
		Name:     "admin",
		Email:    "admin@admin.com",
		Password: pwd1,
		Admin:    true,
	}

	pwd2, _ := auth.HashPassword("test1pwd")
	user2 := models.User{
		Name:     "test1",
		Email:    "test1@test.com",
		Password: pwd2,
	}

	DB.Create(&user1)
	DB.Create(&user2)

	fmt.Println("Seeded user ID: ", user1.ID, user2.ID)

	return nil
}

// populates database with sample cups
func SeedTeams() error {
	var teams models.Cup
	if errors.Is(DB.First(&teams).Error, gorm.ErrRecordNotFound) == false {
		return errors.New("DB already contains teams")
	}

	teams1 := models.Cup{
		Teams: datatypes.JSON([]byte(
			`["Qatar", "Ecuador", "Senegal", "Netherlands"
		, "England", "Iran", "USA", "Wales"
		, "Argentina", "Saudi Arabia", "Mexico", "Poland"
		, "France", "Australia", "Denmark", "Tunisia"
		, "Spain", "Costa Rica", "Germany", "Japan"
		, "Belgium", "Canada", "Morocco", "Croatia"
		, "Brazil", "Serbia", "Switzerland", "Cameroon"
		, "Portugal", "Ghana", "Uruguay", "Korea Republic"]`)),
		Name:   "FIFA 2022 World Cup",
		Format: 1,
		Active: true,
	}
	teams2 := models.Cup{
		Teams: datatypes.JSON([]byte(
			`["Uruguay", "Russia", "Saudi Arabia", "Egypt",
			"Spain", "Portugal", "Iran", "Morocco",
			"France", "Denmark", "Peru", "Australia",
			"Croatia", "Argentina", "Nigeria", "Iceland",
			"Brazil", "Switzerland", "Serbia", "Costa Rica",
			"Sweden", "Mexico", "South Korea", "Germany",
			"Belgium", "England", "Tunisia", "Panama",
			"Colombia", "Japan", "Senegal", "Poland"]`)),
		Name:   "FIFA 2018 World Cup",
		Format: 1,
		Active: true,
	}
	teams3 := models.Cup{
		Teams: datatypes.JSON([]byte(
			`["R. Nadal","J. Draper","B. Nakashima","M. Mcdonald","J. Munar","D. Svrcina","M. Ymer","Y. Nishioka","K. Khachanov",
			"B. Zapata Miralles","S. Baez","J. Kubler","O. Otte","J. Shang","D. Altmaier","F. Tiafoe","H. Hurkacz","P. Martinez",
			"L. Sonego","N. Borges","E. Escobedo","T. Daniel","D. Lajovic","D. Shapovalov","S. Korda","C. Garin","Y. Watanuki",
			"A. Rinderknech","J. Millman","M. Huesler","M. Giron","D. Medvedev","S. Tsitsipas","Q. Halys","Y. Hanfmann","R. Hijikata",
			"T. Griekspoor","P. Kotov","I. Ivashka","B. van de Zandschulp","L. Musetti","L. Harris","F. Coria","M. Fucsovics",
			"T. Etcheverry","G. Barrere","K. Edmund","J. Sinner","C. Norrie","L. Van Assche","T. Monteiro","C. Lestienne","C. Eubanks",
			"S. Kwon","J. Lehecka","B. Coric","F. Cerundolo","G. Pella","C. Moutet","Y. Wu","A. Molcan","S. Wawrinka","V. Pospisil",
			"F. Auger-Aliassime","A. Rublev","D. Thiem","M. Purcell","E. Ruusuvuori","D. Galan","J. Chardy","F. Bagnis","D. Evans",
			"D. Kudla","R. Safiullin","R. Gasquet","U. Humbert","M. Cressy","A. Ramos-Vinolas","F. Krajinovic","H. Rune","P. Carreno Busta",
			"P. Cachin","M. Bellucci","B. Bonzi","J. Isner","A. Mannarino","Y. Hsu","A. de Minaur","G. Dimitrov","A. Karatsev","Z. Bergs",
			"L. Djere","E. Couacaud","H. Dellien","R. Carballes Baena","N. Djokovic","T. Fritz","N. Basilashvili","C. Tseng","A. Popyrin",
			"B. Shelton","Z. Zhang","N. Jarry","M. Kecmanovic","D. Schwartzman","O. Krutykh","J. Wolf","J. Thompson","M. Mmoh","L. Lokoli",
			"J. Varillas","A. Zverev","M. Berrettini","A. Murray","T. Kokkinakis","F. Fognini","A. Vukic","B. Holt","J. Sousa",
			"R. Bautista Agut","A. Davidovich Fokina","A. Bublik","J. Struff","T. Paul","C. O'Connell","J. Brooksby","T. Machac","C. Ruud"]`)),
		Name:   "2023 Australian Open - Men",
		Format: 2,
		Active: true,
	}
	teams4 := models.Cup{
		Teams: datatypes.JSON([]byte(
			`["I. Swiatek","J. Niemeier","P. Udvardy","C. Osorio","C. Bucsa","E. Lys","B. Andreescu","M. Bouzkova","E. Rybakina",
			"E. Cocciaretto","K. Juvan","S. Janicijevic","L. Tsurenko","K. Muchova","A. Kalinskaya","D. Collins","L. Pigossi","C. McNally",
			"K. Rakhimova","K. Baindl","A. Bondar","A. Bogdan","D. Yastremska","J. Ostapenko","Q. Zheng","D. Galfi","B. Pera","M. Uchijima",
			"E. Raducanu","T. Korpatsch","K. Siniakova","C. Gauff","J. Pegula","J. Cristian","B. Fruhvirtova","A. Sasnovich","P. Kudermetova",
			"O. Gadecki","M. Kostyuk","A. Anisimova","B. Krejcikova","S. Bejlek","T. Gibson","C. Burel","C. Vandeweghe","A. Kalinina",
			"A. Van Uytvanck","P. Kvitova","M. Keys","A. Blinkova","X. Wang","S. Hunter","L. Jeanjean","N. Podoroska","S. Kenin",
			"V. Azarenka","J. Teichmann","H. Dart","L. Zhu","R. Marino","K. Kucova","D. Shnaider","Y. Yuan","M. Sakkari","D. Kasatkina",
			"V. Gracheva","L. Stefanini","T. Maria","S. Cirstea","Y. Putintseva","X. Wang 2","K. Pliskova","S. Zhang","P. Tig","P. Martic",
			"V. Golubic","K. Volynets","E. Rodina","M. Zanevska","V. Kudermetova","A. Kontaveit","J. Grabher","M. Sherif","M. Linette","D. Parry",
			"T. Townsend","Y. Bonaventure","E. Alexandrova","I. Begu","E. Mandlik","L. Bronzetti","L. Siegemund","L. Fernandez","A. Cornet",
			"K. Sebov","C. Garcia","A. Sabalenka","T. Martincova","A. Hartono","S. Rogers","L. Davis","D. Kovinic","G. Muguruza","E. Mertens",
			"M. Trevisan","A. Schmiedlova","A. Pavlyuchenkova","C. Giorgi","C. Liu","M. Brengle","V. Tomova","B. Bencic","B. Haddad Maia",
			"N. Parrizas Diaz","S. Stephens","A. Potapova","D. Vekic","O. Selekhmeteva","J. Paolini","L. Samsonova","K. Kanepi","K. Birrell",
			"L. Fruhvirtova","J. Fourlis","A. Riske-Amritraj","M. Vondrousova","T. Zidansek","O. Jabeur"]`)),
		Name:   "2023 Australian Open - Women",
		Format: 2,
		Active: true,
	}
	DB.Create(&teams1)
	DB.Create(&teams2)
	DB.Create(&teams3)
	DB.Create(&teams4)
	fmt.Println("Seeded teams ID: ", teams1.ID, teams2.ID, teams3.ID, teams4.ID)

	return nil

}
