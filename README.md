# golang-api-crud


### Deploy heroku
    heroku login
    heroku create
    heroku config:set MONGO_USER=
    heroku config:set MONGO_HOST=
    heroku config:set MONGO_PORT=
    heroku config:set MONGO_COLLECTION=

### For git deploy
    git push --set-upstream heroku master

### For container deploy
    heroku container:login 
    heroku container:push web
    heroku container:release web
