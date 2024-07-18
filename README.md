# Hi From Ghennadi

Thank you for taking your time to review this code. Don't want to take too much time describing the application as it is as described in the requirements.

Few things I want to note though:

  * I chose Fiber as a framework for the backend because I've heard that it's cool and I wanted to try it out.
  * I'm first and foremost a Python developer and even though I've been coding in Golang for the past ~1.5 years, I'm still learning, so please excuse the weirdnsess if encounter it.
  * This application was very time consuming and there are probably plenty of bugs which I didn't manage to fix.
  * Also there are tongs of things which can (should) be improved. Unit tests, transactons for DB operations, better session management, stronger security, design flaws. Would be happy to discuss all potential issues and my suggestions to improve them during the interview (if I will get the chance)
  * The interface is just a simple CLI application written in Python (I originally wanted to make either UI or Telegram Bot but I don't have time unfortunately)

# How to run things

  *  First and foremost - there is a `.env` file where you can tweak some basic application config.
  * Then I suggest you to just go ahead and bootstrap the docker environment by running `docker-compose up`

# Few words about interface
   
Instead of doing direct API calls, I reccoment using interface which is pretty straightforwad CLI util.

Before using the interface, you need to install the dependencies (just a few). You can do it by running `make deps` or just executing directly `pip install -r requirements.txt`


You can use interface by running:

`./interface.py --help`

or

`python interface.py --help`

Please be aware that every interface command also has a `--help` option, for example:

` ./interface.py create-post --help
`


# Authentication

The authentication is a pretty simple and straightforward method of using username and password. Password is hashed and after registration user can login with credentials to get back a JWT token. The token contains the user auth data.

# API Endpoints And Interface Methods

## /register


**CURL**
```
curl --request POST \
  --url http://127.0.0.1:3000/auth/register \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data name=ghennadi \
  --data password=mypass
```

**INTERFACE**
```
 ./interface.py register POST
```

## /login


**CURL**
```
curl --request POST \
  --url http://127.0.0.1:3000/auth/login \
  --header 'Content-Type: application/json' \
  --data '{
  "name": "ghennadi",
  "password": "mypass"
}'
```

*Note*: You will get back a token which shall be user as a bearer token for all subsequent API reusts.

**INTERFACE**
```
./interface.py login POST
```

*Note*: You will get back a token which will be auto saved to a `.token` file or overwise can be used directly in --token argument to the interface.

## /api/v1/users GET


**CURL**
```
curl --request GET \
  --url http://127.0.0.1:3000/api/v1/users \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1OTI4MDEsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.hjyMbY_7AZEH1IFOz7Edh3MO3w-B9tvYJQUERaV4tNc'
```
**INTERFACE**
```
./interface.py get-users
```

## /api/v1/user/${id}


**CURL**
```
curl --request GET \
  --url http://127.0.0.1:3000/api/v1/users/66995b2a7584f1669ac4ec22 \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1OTI4MDEsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.hjyMbY_7AZEH1IFOz7Edh3MO3w-B9tvYJQUERaV4tNc' \
  --header 'Content-Type: application/json'
```
**INTERFACE**
```
 ./interface.py get-user --id 66995b2a7584f1669ac4ec22
```

## /api/v1/me


**CURL**
```
curl --request GET \
  --url http://127.0.0.1:3000/api/v1/users/me \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1ODU1ODUsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.yqS7NmXcSa6Ow7vYvvEyDmozINVZnHybbA-tI080I0g' \
  --header 'Content-Type: application/json'
```
**INTERFACE**
```
 ./interface.py me
```

## /api/v1/users/${id} PATCH

*Note*: Form request, either name or avatar must be provided.
**CURL**

```
curl --request PATCH \
  --url http://127.0.0.1:3000/api/v1/users/66995b2a7584f1669ac4ec22 \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1ODU1ODUsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.yqS7NmXcSa6Ow7vYvvEyDmozINVZnHybbA-tI080I0g' \
  --header 'Content-Type: multipart/form-data' \
  --form avatar=@cropped_selfie.png
```
**INTERFACE**

```
./interface.py update-user --avatar ~/Downloads/IMG_3408.JPG --id 66995b2a7584f1669ac4ec22
```

## /api/v1/posts POST

*Note* - just all posts, might be useful for debugging

**CURL**

```
curl --request POST \
  --url http://127.0.0.1:3000/api/v1/posts \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1ODU1ODUsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.yqS7NmXcSa6Ow7vYvvEyDmozINVZnHybbA-tI080I0g' \
  --header 'Content-Type: application/json' \
  --data '{
  "content": "https://i.natgeofe.com/n/548467d8-c5f1-4551-9f58-6817a8d2c45e/NationalGeographic_2572187_square.jpg"
}'
```
**INTERFACE**

```
./interface.py create-post --content "https://upload.wikimedia.org/wikipedia/commons/thumb/1/15/Cat_August_2010-4.jpg/362px-Cat_August_2010-4.jpg"
```

## /api/v1/posts GET
**CURL**

```
curl --request GET \
  --url http://127.0.0.1:3000/api/v1/posts \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1ODU1ODUsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.yqS7NmXcSa6Ow7vYvvEyDmozINVZnHybbA-tI080I0g' \
  --header 'Content-Type: application/json'
  ```

**INTERFACE**
```
 ./interface.py get-posts
```

## /api/v1/posts/${id}/like POST

**CURL**

```
curl --request POST \
  --url http://127.0.0.1:3000/api/v1/posts/66997830dd7ff28683926d82/like \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1ODU1ODUsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.yqS7NmXcSa6Ow7vYvvEyDmozINVZnHybbA-tI080I0g' \
  --header 'Content-Type: application/json'
  ```

**INTERFACE**
```
 ./interface.py like-post --id 66997654dd7ff28683926d6f
```

## /api/v1/notifications GET
**CURL**

```
curl --request GET \
  --url http://127.0.0.1:3000/api/v1/notifications \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1ODU1ODUsInVzZXJfaWQiOiI2Njk5NWIyYTc1ODRmMTY2OWFjNGVjMjIifQ.yqS7NmXcSa6Ow7vYvvEyDmozINVZnHybbA-tI080I0g' \
  --header 'Content-Type: application/json'
  ```
**INTERFACE**

```
./interface.py notifications
```

*Note*: Notifications are deleted after retrivial.