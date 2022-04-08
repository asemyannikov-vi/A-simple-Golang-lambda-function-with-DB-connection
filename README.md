# How to launch a λ-function with Database?

To launch a __λ-function__ from the current repo, first prepare a __λ-package__ by typing
```
make prepare-lambda-package
```
and then run the __λ-function__ with __SAM__
```
make launch-lambda-with-sam
```

To gracefuly down the __λ-function__, type
```
make lambda-network-down
```

In order to make sure that everything is working correctly, you need to send an HTTP request using `curl` and get a response from the __λ-function__.

To create a database you should send a `POST` HTTP request
```
curl -X POST localhost:3000/

{
    "message": "Schema 'skeleton' has been created."
}
```

To check a health of database you should send a `GET` HTTP request
```
curl -X GET localhost:3000/

{
    "message": "Schema 'skeleton' has been created and alive."
}
```

To remove a database you should send a `DELETE` HTTP request
```
curl -X DELETE localhost:3000/

{
    "message": "Schema 'skeleton' has been removed."
}
```

Please, read more about Golang __λ-function__ in [this artical](https://teletype.in/@alexander.semyannikov/golang-lambda-function).