let doc = document.getElementsByTagName("body")[0];
let controllerSocket = new WebSocket("ws://localhost:3001/input");
controllerSocket.binaryType = "arraybuffer"
let pressedKeys = {};

window.addEventListener("beforeunload", function() {
    controllerSocket.close(1000, "Page reload");
});

doc.addEventListener("keydown", (e) => {
    pressedKeys[e.key] = true;
});

doc.addEventListener("keyup", (e) => {
    pressedKeys[e.key] = false;
})

let value = 0

controllerSocket.onmessage = (event) => {
    data = new Uint8ClampedArray(event.data)
    if ((data[0] & 0x01) == 0x01) {
        value = 0
    }
    let keys = Object.keys(pressedKeys);
    keys.forEach((key) => {
        if (pressedKeys[key]) {
            let mask = 0;
            switch (key) {
                case "ArrowRight":
                    mask = 1 << 7;
                    break;
                case "ArrowLeft":
                    mask = 1 << 6;
                    break;
                case "ArrowDown":
                    mask = 1 << 5;
                    break;
                case "ArrowUp":
                    mask = 1 << 4;
                    break;
                case "N": // start
                case "n":
                    mask = 1 << 3;
                    break;
                case "M": // select
                case "m":
                    mask = 1 << 2;
                    break;
                case "X": // B
                case "x":
                    mask = 1 << 1;
                    break;
                case "Z": // A
                case "z":
                    mask = 1 << 0;
                    break;
            }
            value |= mask;
        }
    });
    let buffer = new ArrayBuffer(1);
    let view = new Int8Array(buffer);
    view[0] = value;
    controllerSocket.send(view);
}
