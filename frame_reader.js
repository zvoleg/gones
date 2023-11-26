const canvas = document.getElementById("frame");
const ctx = canvas.getContext("2d");
ctx.scale(2, 2);

var socket = new WebSocket("ws://localhost:3000/nes");
socket.binaryType = "arraybuffer";

socket.onmessage = function(event) {
    // var bytes = Uint8Array.from(atob(event.data), c => c.charCodeAt(0));
    var imgData = new Uint8ClampedArray(event.data);
    var image = ctx.createImageData(256,244);
    for (var i = 0; i < image.data.length; i += 1) {
        image.data[i] = imgData[i];
    }
    var resizeWidth = 512;
    var resizeHeight = 488;

    Promise.all([
        // Cut out two sprites from the sprite sheet
        createImageBitmap(image)
      ]).then((img) => {
        // Draw each sprite onto the canvas
        ctx.drawImage(img[0], 0, 0);
      });
}

socket.onerror = function(err) {
    console.log(err)
}