# IA04 Projet 

Pour réaliser le projet ”Dilemme du prisonnier itéré”, vous devez suivre ces étapes :

1. Cloner ce projet dans votre machine propre: 
   $ git clone https://gitlab.utc.fr/shenqiao/ia04-projet.git

2. Installer “nodejs” à partir de ce lien la version LTS https://nodejs.org/en/

2. Vérifier “nodejs” est bien installé et les versions de node et npm:
    $ node -v
    $ npm -v

3. Mise en place d'un projet “React” front-end:
    Entrer dans le répertoire “frontend” :
    $ cd frontend

    Créer une nouvelle application “ReactJS” en utilisant le paquetage “create-react-app” :
    $ npm install -g create-react-app
    $ npx create-react-app .  //dans le répertoire courant
    
    Exécuer l'application ReactJS :
    $ npm start               //pour confirmer l'installation est réussie

    Remplacer le répertoire "src" origine par celui téléchargé

4. Comme le projet ReactJS n'a pas la capacité de gérer les fichiers “.scss”,
    installer “node-sass” dans le répertoire “frontend” :
    Installer “yarn”: 
    $ npm install -g yarn

    Installer “node-sass” de la version 6.0.0: 
    $ yarn add node-sass@6.0.0

5. Démarrer back-end: 
    Entrer dans le répertoire “backend” :
    $ cd backend

    Lancer le server: 
    $ go run main.go

6. Démarrer front-end:
    Entrer dans le répertoire “frontend” :
    $ cd frontend

   Lancer le client: 
    $ npm start
