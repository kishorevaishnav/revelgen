[![baby-gopher](https://raw.github.com/drnic/babygopher-site/gh-pages/images/babygopher-logo-small.png)](http://www.babygopher.org)
revelgen
========

![version](http://img.shields.io/badge/version-pre--α-4ECDC4.svg?style=flat)

This tool is to help Revel(http://revel.github.io) developers initially to generate models, controllers, views either inidividually or by scaffold.

Here are the steps how to install:

```
go get github.com/kishorevaishnav/revelgen
sh build.sh
```

Below mentioned is the an example of one of the normal usage if you are using GORP:

```
revel new accounts
cd $GOPATH/src/accounts
revelgen scaffold ledger id:int name:string order:int status:bool
```
Note: Other ORM is in progress.


Other simple usage:

```
cd $GOPATH/src/accounts
revelgen model ledger id:int name:string order:int status:bool
revelgen controller ledger index show list
```

Please feel free to file Issues if you see any problem or wanted something to add new.
