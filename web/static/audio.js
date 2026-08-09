const SAMPLE_RATE = 44100;

let AudioContext = window.AudioContext || webkitAudioContext.AudioContext;
let audioCtx = null;
let player = null;

// const playPause = document.getElementsByClassName("playPause")[0];
playPauseButton.addEventListener("click", function() {
    InitContext().then(function() {
            if (audioCtx.state === "suspended") {
                audioCtx.resume();
            } else if (audioCtx.state === "running") {
                audioCtx.suspend();
            }
        }
    );
})

async function InitContext() {
    if (!audioCtx) {
        audioCtx = new AudioContext({ sampleRate: SAMPLE_RATE});
        await audioCtx.audioWorklet.addModule("static/audio_worklet.js");
    }

    PlaySound();
}

function PlaySound() {
    const exampleWorkletNode = new AudioWorkletNode(audioCtx, "audio");
    exampleWorkletNode.connect(audioCtx.destination);
    player = new Player(exampleWorkletNode);
}

class Player {
    constructor(audioWorklet) {
        this.audioWorklet = audioWorklet;
        this.dataRecaiverSocket = new WebSocket("ws://localhost:3011/audio");
        this.dataRecaiverSocket.binaryType = "arraybuffer";

        this.dataRecaiverSocket.onmessage = async (event) => {
            this.audioWorklet.port.postMessage(new Float32Array(await event.data));
        };
        this.audioWorklet.port.onmessage = (event) => {
            console.log(event.data)
        }
    }
}