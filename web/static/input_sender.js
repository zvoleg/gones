var doc = document.getElementsByTagName("body")[0];
var controllerSocket = new WebSocket("ws://localhost:3000/input");
controllerSocket.binaryType = "arraybuffer"
var pressedKeys = {};

doc.addEventListener("keydown", (e) => {
    pressedKeys[e.key] = true;
});

doc.addEventListener("keyup", (e) => {
    pressedKeys[e.key] = false;
})

controllerSocket.onmessage = (event) => {
    var keys = Object.keys(pressedKeys);
    var value = 0;
    keys.forEach((key) => {
        if (pressedKeys[key]) {
            var mask = 0;
            switch (key) {
                case "ArrowRight":
                    mask = 1 << 0;
                    break;
                case "ArrowLeft":
                    mask = 1 << 1;
                    break;
                case "ArrowDown":
                    mask = 1 << 2;
                    break;
                case "ArrowUp":
                    mask = 1 << 3;
                    break;
                case "N": // start
                case "n":
                    mask = 1 << 4;
                    break;
                case "M": // select
                case "m":
                    mask = 1 << 5;
                    break;
                case "X": // B
                case "x":
                    mask = 1 << 6;
                    break;
                case "Z": // A
                case "z":
                    mask = 1 << 7;
                    break;
            }
            value |= mask;
        }
    });
    console.log(value);
    var buffer = new ArrayBuffer(1);
    var view = new Int8Array(buffer);
    view[0] = value;
    controllerSocket.send(view);
}
