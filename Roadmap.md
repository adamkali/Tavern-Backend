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

## TAVERN V1 1/8

### V1 LANDING PAGE

* Create a landing page to get infromation about the app ❌
    - [ ] MAKE A QUICK APP FOR LANDING PAGE
    - [ ] HAVE ALL THE BENEFITS
    - [ ] HAVE INFORMATION
      - [ ] FAQ
      - [ ] CAREERS (NOT TILL LATER AND NEED ALOT FROM PATREON)
      - [ ] SOCIAL LINKS (PATREON, INSTAGRAM, TIKTOK, YOUTUBE)

### V1 BACKEND

* Refactor endpoints to be more abstract ✅
    - [x] IOWRITER ABSTRACTION
    - [x] STRUCT ABSTRACTION
    - [x] REPOSITORY ABSTRACTION
    - [x] ADD CONCURRENCY

* Create `/api/groups` endpoints ❌
    - [ ] GET GROUP BY ID
    - [ ] UPDATE GROUP
    - [ ] DELETE GROUP
    - [ ] POST GROUP

* Research Tinder like ❌
    - [ ] Messaging
    - [ ] Swipe
    - [x] Photo Storeing
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

* Make a login experience ❌
    - [ ] SHOULD HAVE <> STAGES:
      - [ ] SHOULD HAVE A LOGIN PAGE
        - [ ] NORMAL LOGIN 
        - [ ] A BUTTON FOR SIGNUP
        - [ ] A BUTTON FOR FORGOT PASSWORD
      - [ ] ON SIGNUP IT SHOULD HAVE
        - username
        - password
        - email
      - [ ] THEN IF THE REQUEST WAS SUCSESSFUL
        - [ ] HAVE A SPACE TO INPUT A CODE
        - [ ] HAVE A BUTTON TO RESEND EMAIL
        - [ ] HAVE A BUTTON TO ENTER
        - [ ] THIS SHOULD SAVE THE AUTH TOKEN INTO THE LOCAL STORAGE
        - [ ] IF THE `/api/activate/{code}` IS SUCCESSFUL START MAKING A USER
      - [ ] MAKE A USER CREATION PAGE
        - [ ] THIS SHOULD HAVE DIFFERENT PAGES THAT GOES THROUGH WHAT THE USER WANTS
        - [ ] ONE PAGE FOR A BIO
        - [ ] MAKE A PAGE TO SELECT THEIR PLAYER PREFERENCE
          - Player
          - Dungeon Master
          - Both
          - Both (Prefers Player)
          - Both (Prefers Dungeon Master)
          - Just Talk
        - [ ] MAKE A PAGE TO SELECT THEIR TAGS
        - [ ] DEPENDING WHAT THEY HAVE CHOSEN IT WILL SHOW:
          - [ ] CHARACTER SHEET
          - [ ] PLOT TYPE

* Get Info from API ❌
    - [ ] Have data display
    - [ ] Optimization

* Have user preferences ❌
    - [ ] Have a way to set location
    - [ ] Have a way to set themes
    - [ ] Dice? 🎲

## TAVERN V2 2/4

* Implement authentication  ✅
    - [x] Implement auth api's and tokens
    - [x] Implement login

* Implement User settings ❌
    - [ ] store user settings in NoSQL
    - [ ] Have settings load in app at startup

* Research Server Solutions ✅
    - [x] Get https cert
    - [x] Domain for parent site (www.taverndnd.app) 
    - [x] Find good place to deploy code (Linode)
    - [x] Reaserch app deployment (Linode w/Docker Container)

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
