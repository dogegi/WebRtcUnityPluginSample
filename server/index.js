var app = require('express')();
var https = require('https');
var http = require('http').Server(app);
const fs = require('fs');

var io = require('socket.io')(http, {
  transports: ['polling']
});

const ngrok = require('ngrok');

const options = {
  cert: fs.readFileSync('cert.pem'),
  key: fs.readFileSync('key.pem'),
  requestCert: false,
  rejectUnauthorized: false
};

app.get('/', function(req, res){
  res.sendFile(__dirname + '/index.html');
});

// simple and dirty WebRTC signaling (without sending each ice candidates)
// see: https://qiita.com/massie_g/items/f5baf316652bbc6fcef1
var userId = 0;
io.on('connection', function(socket){
  socket.userId = userId ++;
  console.log('a user connected, user id: ' + socket.userId);
  socket.emit("welcome", { id: socket.userId });
  io.emit("join", { id: socket.userId } );

  socket.on('chat', function(msg){
    console.log('message from user#' + socket.userId + ": " + msg);
    io.emit('chat', {
      id: socket.userId,
      msg: msg
    });
  });
  socket.on('webrtc-offer', function(msg){
    console.log('webrtc-offer from user#' + socket.userId + ": " + msg);
    io.emit('webrtc-offer', {
      id: socket.userId,
      msg: msg
    });
  });
  socket.on('webrtc-answer', function(msg){
    console.log('webrtc-answer from user#' + socket.userId + ": " + msg);
    io.emit('webrtc-answer', {
      id: socket.userId,
      msg: msg
    });
  });
  socket.on('webrtc-icecandidate', function(msg){
    console.log('webrtc-icecandidate from user#' + socket.userId + ": " + msg);
    io.emit('webrtc-icecandidate', {
      id: socket.userId,
      msg: msg
    });
  });
  socket.on('disconnect', function() {
    io.emit("exit", { id: socket.userId });
  })
});

const sport = 3001;
https.createServer(options, app).listen(sport, function(){
//http.listen(port, function() {
  console.log('listening on *:' + sport);
});

const port = 3000;
http.listen(port, function() {
  console.log('listening on *:'+port);
  console.log('forwarding to global domain...')
  ngrok.connect(port, (err, url) => {
    console.log("ngrok server: "+url);
  });
});
