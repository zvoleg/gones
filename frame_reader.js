var frameSocket = new WebSocket("ws://localhost:3000/frame");
frameSocket.binaryType = "arraybuffer";
var palletteSocket = new WebSocket("ws://localhost:3000/pallette");
palletteSocket.binaryType = "arraybuffer";
var patternSocket = new WebSocket("ws://localhost:3000/pattern");
patternSocket.binaryType = "arraybuffer";
var nameTableSocket = new WebSocket("ws://localhost:3000/name");
nameTableSocket.binaryType = "arraybuffer";

setupSocketAndCanvas(frameSocket, "frame", 256, 240, 3)
setupSocketAndCanvas(palletteSocket, "pallette", 9, 5, 20)
setupSocketAndCanvas(patternSocket, "pattern", 256, 128, 2)
setupSocketAndCanvas(nameTableSocket, "nameTable", 512, 480, 1)

function setupSocketAndCanvas(socket, canvasName, width, height, scale) {
  const canvas = document.getElementById(canvasName);
  canvas.style.width = (width * scale) + "px";
  canvas.style.height = (height * scale) + "px";
  canvas.width = width;
  canvas.height = height;
  canvas.style.imageRendering = "pixelated"
  const ctx = canvas.getContext("2d");
  ctx.imageSmoothingEnabled = false;

  socket.onopen = () => {
    socket.send(canvasName);
  }

  socket.onmessage = (event) => {
      var imgData = new Uint8ClampedArray(event.data);
      var image = ctx.getImageData(0, 0, width, height);
      for (var i = 0; i < image.data.length; i += 1) {
          image.data[i] = imgData[i];
      }

      ctx.putImageData(image, 0, 0);
  }

  socket.onerror = (err) => {
      console.log(err);
  }
}