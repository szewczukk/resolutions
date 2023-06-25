# Resolutions - projekt zaliczeniowy na Aplikacje Backendowe

## Architektura

Aplikacja składa się z frontendu, gateway serwisu, oraz trzech serwisach funkcjonalnych: user, score i resolution service.

Frontend jest napisany w React-ie i komunikuje się za pomocą `REST API` z serwisem bramowym, który wystawia serwisy funkcjonalne na żądania HTTP. Serwis bramowy komunikuje się z serwisami funkcjonalnymi za pomocą protokołu `gRPC`. A serwisy funkcjonalne komunikują się między sobą synchronicznie za pomocą `gRPC` oraz asynchronicznie przez `RabbitMQ`. Uwierzytelnianie użytkowników następuje za pomocą tokenów `JWT`. System automatycznie tworzy oddzielną bazę `SQLite` dla każdego mikroserwisu, przez co mikroserwisy są bardziej niezależne oraz nie następuje coupling.

## Jak uruchomić projekt

Oprogramowanie wymagane do uruchomienia systemu:

1. RabbitMQ, którego opis instalacji znajduje się: https://www.rabbitmq.com/download.html
1. Język Go w wersji 1.20.5: https://go.dev/dl/
1. Menadżer pakietów NodeJS (https://nodejs.org/en), najlepiej pnpm: https://pnpm.io/

Przed uruchomieniem projektu, należy uruchomić RabbitMQ, a następnie należy każdy katalog mikroserwisu, które można znaleźć w katalogu `services`, uruchomić w oddzielnym okienku terminala, a w nich uruchomić komendę `go run .`. Zależności zainstalują się automatycznie.

Następnie, aby uruchomić aplikację internetową, należy przejść do folderu `frontend` oraz zainstalować zależności komendą `pnpm i`, a następnie uruchomić projekt przez `pnpm dev`.

## Funkcjonalności projektu

1. Rejestracja
1. Logowanie, które wykorzystuje tokeny JWT
1. Tablica wyników, która jest dostępna dla niezalogowanych użytkowników i wyświetla wszystkich użytkowników w systemie i ich zdobyte punkty
1. Dashboard użytkownika, na którym można tworzyć, usuwać oraz kończyć pewne postanowienia (za co otrzymuje się punkty)
