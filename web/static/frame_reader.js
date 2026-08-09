let frameSocket = new WebSocket("ws://localhost:3000/frame");
frameSocket.binaryType = "arraybuffer";
let paletteSocket = new WebSocket("ws://localhost:3000/palette");
paletteSocket.binaryType = "arraybuffer";
let patternSocket = new WebSocket("ws://localhost:3000/pattern");
patternSocket.binaryType = "arraybuffer";
let nameTableSocket = new WebSocket("ws://localhost:3000/name");
nameTableSocket.binaryType = "arraybuffer";

window.addEventListener("beforeunload", function() {
  console.log("Page reload");
  frameSocket.close(1000, "page reload");
  paletteSocket.close(1000, "page reload");
  patternSocket.close(1000, "page reload");
  nameTableSocket.close(1000, "page reload");
});

setupSocketAndCanvas(frameSocket, "frame", 256, 240, 3)
setupSocketAndCanvas(paletteSocket, "palette", 9, 5, 20)
setupSocketAndCanvas(patternSocket, "pattern", 256, 128, 2)
setupSocketAndCanvas(nameTableSocket, "nameTable", 512, 512, 1)

function setupSocketAndCanvas(socket, canvasName, width, height, scale) {
  const canvas = document.getElementById(canvasName);
  canvas.style.width = (width * scale) + "px";
  canvas.style.height = (height * scale) + "px";
  canvas.width = width;
  canvas.height = height;
  canvas.style.imageRendering = "pixelated";
  const ctx = canvas.getContext("2d", { willReadFrequently: true });
  ctx.imageSmoothingEnabled = false;

  let imgDataBuffer = new Uint8ClampedArray(width*height*4);
  let imgData = ctx.createImageData(width, height);
  
  socket.onopen = () => {
    socket.send(canvasName);
  }
  
  socket.onmessage = (event) => {
    imgDataBuffer.set(new Uint8ClampedArray(event.data));
    imgData.data.set(imgDataBuffer);
    ctx.putImageData(imgData, 0, 0);
  }

  socket.onerror = (err) => {
      console.log(err);
  }
}