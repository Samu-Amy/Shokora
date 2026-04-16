# Test mancanti sulle route

## Auth — handler non testati:

- `GET /auth/google` — verifica che ritorni un URL OAuth valido
- `POST /auth/google/callback` — con code valido/invalido/mancante
- `GET /auth/logout~` — con e senza sessione attiva, verifica pulizia cookie
- `POST /auth/verify-email/send~` — con email esistente, non esistente, già verificata
- `POST /auth/verify-email/otp` — OTP valido, scaduto, errato, già usato
- `POST /auth/verify-email/{token}` — magic link valido, scaduto, malformato
- `POST /auth/reset-password/send` — stessi casi
- `POST /auth/reset-password/otp e /{token}` — verifica ritorno reset session token
- `PATCH /auth/reset-password` — con reset token valido/scaduto, nuova password uguale alla vecchia
- `POST /auth/verify-2fa/send e /otp` — verifica che senza 2FA attivo non si possa usare

## Route protette:

- `GET /me` — senza auth (401), con auth valida, con token scaduto
- `PATCH /update-password` — vecchia password errata, nuova uguale alla vecchia, utente non autenticato
- `PATCH /user` — aggiornamento dati validi/invalidi, utente non autenticato
- `GET /employee/users/{userId}` — con ruolo customer (403), con ruolo employee (200), userId inesistente (404)
- `Route prodotti employee` — con/senza permessi specifici, con ruolo customer

## Middleware da testare esplicitamente:

- `authMiddleware` — token mancante, malformato, scaduto, sessione revocata
- `roleMiddleware` — accesso a /employee/\* con ruolo customer/admin/dev
- `userVerifiedMiddleware` — utente non verificato su route protette
- `permissionMiddleware` — employee senza il permesso specifico (es. EmplPermProductAdd)
- `Rate limiter` — superamento soglia (verifica 429)

# Test sui service (i più utili)

## Questi coprono logica che le route tests non possono isolare bene:

- `Auth.RegisterUser` — email duplicata, trimming campi, hashing password, invio mail verifica
- `Auth.LoginUser` — utente con 2FA attivo (verifica che non ritorni il token ma il verificationId), utente non verificato (verifica che ritorni comunque l'utente ma con flag), password errata
- `Auth.HandleAuthTokensCheck` - con/senza access/refresh token (oppure scaduto/i), reuse refresh token
- `Auth.VerifyEmailWithOTP` / WithMagicLink — OTP già usato, scaduto, appartenente ad altro utente
- `Auth.ResetPassword` — verifica che dopo il reset la vecchia password non funzioni più, che il reset token sia invalidato
- `Auth.TwoFactorAuthWithOTP` — verifica che dopo il 2FA la sessione sia creata correttamente
- `Auth.LogoutUser` — verifica che la sessione sia effettivamente eliminata dal DB e che il token non sia più valido

## Note

- Aggiungere un helper makeAuthenticatedRequest che gestisce il cookie di sessione, così i test delle route protette diventano molto più semplici da scrivere
