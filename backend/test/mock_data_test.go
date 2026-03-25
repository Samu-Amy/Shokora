package main

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

// TODO: crea dati di base per users e products

// TODO: fai db di test (compose come per dev)

// ----- DATA -----

var mockValidFirstNames = []string{
	"Mario", "Giulia", "Luca", "Sara", "Alessandro", "Francesca", "Matteo", "Chiara", "Andrea", "Martina", "Federico", "Elena", "Davide",
	"Valentina", "Simone", "Laura", "Nicola", "Giovanna", "Gabriele", "Alice", "Stefano", "Ilaria", "Tommaso", "Beatrice", "Riccardo",
	"Francesco", "Vanessa", "Alessia", "Emanuele", "Giada", "Antonio", "Sofia", "Daniele", "Claudia", "Paolo", "Camilla", "Enrico",
	"Michela", "Vincenzo", "Elisa", "Roberto", "Federica", "Salvatore", "Aurora", "Giorgio", "Veronica", "Fabio", "Chiara", "Simona", "Lorenzo",
}

var mockValidLastNames = []string{
	"Rossi", "Bianchi", "Verdi", "Ferrari", "Moretti", "Galli", "Rinaldi", "Romano", "Conti", "Costa", "Fontana", "Marini", "Ricci", "De Luca", "Longo", "Martini", "Barbieri",
	"Grassi", "Giordano", "Cattaneo", "Villa", "Serra", "Pellegrini", "Lombardi", "Villa", "", "Sanna", "Bruno", "Esposito", "Caputo", "Santoro", "D'Amico", "Vitale",
	"Gatti", "Sala", "Piras", "Bertoli", "Amato", "Testa", "Corsi", "Pagani", "De Santis", "Fabbri", "Monti", "Bernardi", "Ruggiero", "Negri", "Ferretti", "Barone", "",
}

var mockValidEmails = []string{
	"mario.rossi@example.com", "giulia.bianchi@example.com", "luca.verdi@example.com", "sara.ferrari@example.com", "alessandro.moretti@example.com", "francesca.galli@example.com",
	"matteo.rinaldi@example.com", "chiara.romano@example.com", "andrea.conti@example.com", "martina.costa@example.com", "federico.fontana@example.com", "elena.marini@example.com",
	"davide.ricci@example.com", "valentina.deluca@example.com", "simone.longo@example.com", "laura.martini@example.com", "nicola.barbieri@example.com", "giovanna.grassi@example.com",
	"gabriele.giordano@example.com", "alice.cattaneo@example.com", "stefano.villa@example.com", "ilaria.serra@example.com", "tommaso.pellegrini@example.com",
	"beatrice.lombardi@example.com", "riccardo.villa@example.com", "francesco.parisi@example.com", "vanessa.sanna@example.com", "alessia.bruno@example.com",
	"emanuele.esposito@example.com", "giada.caputo@example.com", "antonio.santoro@example.com", "sofia.damico@example.com", "daniele.vitale@example.com", "claudia.gatti@example.com",
	"paolo.sala@example.com", "camilla.piras@example.com", "enrico.bertoli@example.com", "michela.amato@example.com", "vincenzo.testa@example.com", "elisa.corsi@example.com",
	"roberto.pagani@example.com", "federica.desantis@example.com", "salvatore.fabbri@example.com", "aurora.monti@example.com", "giorgio.bernardi@example.com",
	"veronica.ruggiero@example.com", "fabio.negri@example.com", "chiara.ferretti@example.com", "simona.barone@example.com", "lorenzo.bellini@example.com",
}

var mockValidPasswords = []string{
	"Password123!",                                 // lettere maiuscole, minuscole, numeri, simbolo
	"justlettersabcd",                              // solo lettere minuscole
	"UPPERCASELETTERS12",                           // maiuscole + numeri
	"mixedCASEpassword",                            // maiuscole e minuscole
	"123456789012",                                 // solo numeri
	"complex!Pass#01",                              // lettere + numeri + simboli
	"shortButGood12!",                              // vicino a 12 caratteri
	"LongPasswordExample1234567890!",               // lunga, sicura
	"onlylowercaseletters",                         // solo lettere minuscole
	"MiXeD123456",                                  // mix lettere/numeri
	"Symbols!@#Only12",                             // solo simboli + numeri + lettere
	"JustLettersLONGname",                          // lettere lunghe
	"1234abcd5678efgh",                             // lettere e numeri
	"Mix3dCASE!@#",                                 // lettere, numeri e simboli
	"abcdefgHIJKLMN",                               // lettere maiuscole e minuscole
	"abcdefgh12345678",                             // minuscole + numeri
	"ONLYUPPERCASELETTERS",                         // solo maiuscole
	"!@#ComplexSymbol12",                           // simboli + lettere + numeri
	"lowerUPPER123",                                // mix semplice
	"SomeRandomPass2024",                           // lettere e numeri
	"SymbolsAndLetters!@#",                         // mix
	"PasswordWithLONGText1234567890",               // lunga
	"abcDEF123!@#",                                 // mix corto
	"justlowercaseletters2",                        // minuscole + numero
	"UPPERandlower123",                             // mix maiuscole/minuscole/numeri
	"Special$$Symbols123",                          // simboli + numeri + lettere
	"SimplepassWord12",                             // semplice
	"AnotherGoodPass!@",                            // simboli + lettere
	"lowercasewithnumber1",                         // minuscole + numero
	"UPPERCASEWITHSYMBOLS!@",                       // maiuscole + simboli
	"RandomLongPasswordExample123!",                // lunga e sicura
	"MiXeDletters123!",                             // mix corto
	"PasswordMinimal12",                            // vicino a 12 caratteri
	"ComplexLONGPasswordWith123Symbols!@#",         // lunga, sicura
	"JustLettersWithUppercase",                     // lettere + maiuscole
	"lowercaseandUPPER12",                          // mix semplice
	"MixedWith123Symbols!@",                        // mix medio
	"LowerUPPERSymbols123",                         // mix
	"AnotherPass123456",                            // lettere + numeri
	"ShortSym!@#12",                                // vicino a 12 caratteri
	"LongPasswordWithLettersNumbersAndSymbols123!", // molto lunga
	"lettersnumbers123",                            // semplice
	"UPPERlower123!",                               // corto, mix
	"Symbols123!@#",                                // simboli + numeri
	"RandomPassWith12345",                          // lettere + numeri
	"JustLetters12345",                             // lettere + numeri
	"LowerCaseOnlyabcd",                            // minuscole
	"UpperCASEOnlyABCD",                            // maiuscole
	"MixedLettersWith!@#",                          // mix simboli + lettere
	"SafePassword2026!",                            // semplice ma sicura
	"VeryLongPasswordWithLetters123!@#",            // lunga e sicura
}

var mockBirthdays = []time.Time{
	// TODO: fai date casuali come string "DD-MM"
}

// ----- FUNCTIONS -----

func makeRegisterUserReq(firstName, lastName, email, password, passwordConfirmation, birthday string) payloads.RegisterUserReq {
	return payloads.RegisterUserReq{
		UserDataReq: payloads.UserDataReq{
			FirstName: firstName,
			LastName:  lastName,
			Birthday:  birthday,
		},
		EmailFieldReq: payloads.EmailFieldReq{
			Email: email,
		},
		DoublePasswordFieldReq: payloads.DoublePasswordFieldReq{
			Password:             password,
			PasswordConfirmation: passwordConfirmation,
		},
	}
}
