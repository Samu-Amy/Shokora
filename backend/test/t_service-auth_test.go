package main

// TODO: testa tokens (controlla scadenza, ma soprattutto controlla l'estensione della scadenza dei refresh token e i controlli su reuse, con eliminazione della sessione)

// TODO: testa invio mail (anche casi di errore) e controlla che funzioni in entrambi i casi all'interno dei service che la usano (register, login (2FA), verification, ecc.)

// TODO: test login (sia normale che con google, anche insieme (uno dopo l'altro) per verificare che un account creato con email e password non dia problemi a google ed un account creato con google (quindi senza password) non dia problemi e non permetta login con email e password)
