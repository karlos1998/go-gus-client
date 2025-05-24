# gus-client

Mikroserwis w Go udostępniający dane firm z rejestru GUS (BIR1.1) na podstawie numeru NIP.

Do działania wymagany jest REGON_API_KEY, w przykładzie podany jako 'twój_klucz'. 
Musisz uzyskać go we własnym zakresie z REGON.

## Endpointy

- `GET /health` – status serwisu (zawsze dostępny)
- `GET /search?nip=1234567890` – wyszukiwanie danych firmy po NIP (może wymagać autoryzacji)

## Autoryzacja

Serwis może działać w dwóch trybach:

### 🔓 Tryb publiczny (domyślny)

Jeśli nie ustawisz zmiennej środowiskowej `ACCESS_TOKEN`, zapytania są dostępne publicznie.

### 🔒 Tryb z zabezpieczeniem tokenem

Jeśli ustawisz `ACCESS_TOKEN` w `.env`, każde zapytanie do `/search` wymaga autoryzacji:

#### Opcja 1: przez parametr query

```
GET /search?nip=1234567890&auth=sekretnytoken
```

#### Opcja 2: przez nagłówek HTTP

```
Authorization: Bearer sekretnytoken
```

#### Błąd

Jeśli token jest błędny lub nie został podany – zwracany jest błąd 401:

```json
{
  "error": "Nieprawidłowy klucz autoryzacyjny. Upewnij się czy w parametrzy query podałeś auth taki jak ustawiłeś w .env jako ACCESS_TOKEN. zapytanie powinno wyglądać tak: http://localhost:4300/search/?nip=1234567890&auth=sekretnytoken lub auth możesz podać jako header Authorization: Bearer sekretnytoken"
}
```

## Przykład odpowiedzi

```json
[
  {
    "Type": "P",
    "NIP": "5261040828",
    "REGON": "000331501",
    "Name": "GŁÓWNY URZĄD STATYSTYCZNY",
    "Street": "ul. Test-Krucza",
    "PropertyNumber": "208",
    "ApartmentNumber": "",
    "PostalCode": "00-925",
    "City": "Warszawa"
  }
]
```

## Uruchomienie

```bash
docker run -p 4300:4300 \
  -e REGON_API_KEY=twój_klucz \
  -e ACCESS_TOKEN=sekretnytoken \
  karlos98/gus-client:latest
```
