var frameSocket = new WebSocket("ws://localhost:3000/frame");
frameSocket.binaryType = "arraybuffer";
var patternSocket = new WebSocket("ws://localhost:3000/pattern");
patternSocket.binaryType = "arraybuffer";

setupSocketAndCanvas(frameSocket, "frame", 256, 244, 3)
setupSocketAndCanvas(patternSocket, "pattern", 256, 128, 3)

function setupSocketAndCanvas(socket, canvasName, width, height, scale) {
  const canvas = document.getElementById(canvasName);
  const ctx = canvas.getContext("2d");
  ctx.scale(scale, scale);

  socket.onopen = () => {
    socket.send(canvasName)
  }

  socket.onmessage = (event) => {
      var imgData = new Uint8ClampedArray(event.data);
      var image = ctx.createImageData(width, height);
      for (var i = 0; i < image.data.length; i += 1) {
          image.data[i] = imgData[i];
      }
      var resizeWidth = width * scale;
      var resizeHeight = height * scale;

      Promise.all([
          createImageBitmap(image)
        ]).then((img) => {
          ctx.drawImage(img[0], 0, 0);
        });
  }

  socket.onerror = (err) => {
      console.log(err)
  }
}