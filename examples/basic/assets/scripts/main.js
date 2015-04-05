console.log('Woop');
addEventListener('foo', function (player) {
  var pos = player.Pos();
  console.log("Got 'foo' event! Player X:", pos.X, "Player Y:", pos.Y);
});
