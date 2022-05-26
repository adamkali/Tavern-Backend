# Roadmap for Tavern

## BASIC FUNCTIONALITY 4/4

* Create `api/users/` endpoints ✅
    - [x] GET ALL
    - [x] GET BY ID
    - [x] POST ONE
    - [x] PUT ONE BY ID
    - [x] DELETE

* Create `api/characters` endpoints ✅
    - [x] GET BY ALL FOR USER
    - [x] POST ONE FOR USER
    - [x] UPDATE ONE FOR USER
    - [x] DELETE ONE FROM USER

* Create `api/plots` endpoints for a user ✅
    - [x] GET BY ID
    - [x] POST ONE
    - [x] PUT ONE BY ID
    - [x] DELETE

* Create repositorys ✅
    - [x] CONNECT TO MYSQL
    - [x] CREATE A SQL DATABASE
    - [x] CREATE TABLES
    - [x] IMPLEMENT TABLES INTO ENDPOINTS

## TAVERN V1 0/6

### V1 BACKEND

* Refactor endpoints to be more abstract ❌
    - [x] IOWRITER ABSTRACTION
    - [x] STRUCT ABSTRACTION
    - [x] REPOSITORY ABSTRACTION
    - [ ] ADD CONCURRENCY

* Create `/api/groups` endpoints ❌
    - [ ] GET GROUP BY ID
    - [ ] UPDATE GROUP
    - [ ] DELETE GROUP
    - [ ] POST GROUP

* Research Tinder like ❌
    - [ ] Messaging
    - [ ] Swipe
    - [ ] Photo Storeing
    - [ ] Geolocation

### V1 FRONTEND

* Make a mock app ❌
    - [ ] MAKE COLORSCHEME
    - [ ] MAKE BASIC MOCKUP AND LOOK AND FEEL
    - [ ] BIO PAGE
    - [ ] SWIPE PAGE
        * [ ] SUCCESS
        * [ ] NAT TWENTY
        * [ ] FALIURE
    - [ ] MESSAGE PAGE
        * [ ] HAVE CHATS FOR INDIVIDUALS
        * [ ] HAVE CHATS FOR GROUPS
        * [ ] HAVE DICE GAMES

* Get Info from API ❌
    - [ ] Have data display
    - [ ] Optimization

* Have user preferences ❌
    - [ ] Have a way to set location
    - [ ] Have a way to set themes
    - [ ] Dice? 🎲

## TAVERN V2 0/4

* Implement authentication  ❌
    - [ ] Implement auth api's and tokens
    - [ ] Implement login

* Implement User settings ❌
    - [ ] store user settings in NoSQL
    - [ ] Have settings load in app at startup

* Research Server Solutions ❌
    - [ ] Get https cert
    - [ ] Find good place to deploy code
    - [ ] Reaserch app deployment

* Frontend Refinment  ❌
    - [ ] Animations
    - [ ] Refine components
    - [ ] Implement multiple themes
    - [ ] 3D?
    - [ ] Icon

## TAVERN V3 0/3

* Publishment
    - [ ] Deploy to servers ❌
    - [ ] Make Patreon
    - [ ] Get on Stores
    - [ ] Google Adverts

* Add More Games ❌
    - [ ] Pathfinder
    - [ ] Call of Cthulu
    - [ ] Others

* Make metrics API ❌
    - [ ] Endpoint for activity
    - [ ] Endpoint for current users
    - [ ] Endpoint for messaging logs
