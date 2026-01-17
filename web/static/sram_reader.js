var sramSocket = new WebSocket("ws://localhost:3002/sram")
var sramCounter = 0
var sramEntities = null

sramSocket.onopen = () => {
    var sramBox = document.getElementsByClassName("sram")[0];
    for (i = 0; i < 64; i++) {
        var sramEntity = document.createElement("p");
        sramEntity.setAttribute("style", "margin: 0; padding: 0; font-size: 12px; font-family: monospace;");
        sramEntity.appendChild(document.createTextNode(""));
        sramBox.appendChild(sramEntity);
    }
    sramEntities = sramBox.getElementsByTagName("p")
}

sramSocket.onmessage = (event) => {
    var sramEntity = sramEntities[sramCounter]
    sramEntity.firstChild.nodeValue = String(sramCounter).padStart(2, '0') + ": " + event.data;
    if (sramEntity.firstChild.nodeValue.includes("y: 255")) {
        sramEntity.style.backgroundColor = "Salmon"
    } else {
        sramEntity.style.backgroundColor = "PaleGreen"
    }

    sramCounter = (sramCounter + 1) % 64;
}

sramSocket.onerror = (err) => {
    console.log(err);
}