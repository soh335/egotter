# egotter

mention, event, keyword notification by im.kayac.com for twitter

## how to use

### configuration

```
git clone ...
git checkout -b heroku
cp config.json.sample config.json
# edit config.json
git add config.json
git commit -m 'add config.json'
```

### deploy to heroku

```
git checkout heroku
git merge --no-ff master
git push heroku heroku:master
```
