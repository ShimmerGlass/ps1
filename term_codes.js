var fs = require('fs');

var stdinBuffer = fs.readFileSync(0);
var json = JSON.parse(stdinBuffer.toString());

for (var i in json) {
  console.log(
      i + ": HSL{" + json[i].hsl.h + ", " + json[i].hsl.s + ", " +
      json[i].hsl.l + "},");
}