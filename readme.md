# TimeStream

TimeStream is a proof of concept that we can time how long http content is viewed from the server side. The use case in mind guaging engagement with emails where javascript is not allowed. 

### Usage

Simply ```go run main.go```. You can then access ```localhost:9999/someurl```. If you cancel the curl or stop the HTTP session in anyway, the logs will update with the length of time that connection was open. 

This gets more interesting if you pass a base64 encoded URI.

```localhost:9999/dWlkOjE4MCBtc2dpZDoxMDA0IHJjcHQ6Zm9vQGV4YW1wbGUuY29t.gif``` will result in a log entry with the duration of the connection (in seconds) and the base64 decoded information in the URI. 

### Proof of Concept in Action

Using [ngrok](https://ngrok.com/), you can then send an email pointing to your locally running TimeStream server and, if the user is loading images, you can see how long they viewed your email!

One drawback so far is that it seems that Gmail keeps the connection open for 20 seconds even if the user clicks back to their inbox earlier than that. 
