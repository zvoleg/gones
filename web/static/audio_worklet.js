class ExampleWorklet extends AudioWorkletProcessor {
    constructor() {
        super();
        this.bufferSize = 4096;
        this.buffer = new Float32Array(this.bufferSize).fill(0.0);
        this.readPtr = 0;
        this.writePtr = 0;
        this.availableData = 0;

        this.port.onmessage = (event) => {
            const data = event.data;
            const freeSpace = this.bufferSize - this.availableData;
            if (freeSpace < data.length) {
                const toSkip = data.length - freeSpace;
                this.readPtr = (this.readPtr + toSkip) % this.bufferSize;
                this.availableData -= toSkip;
            }
            event.data.forEach(sample => {
                this.buffer[this.writePtr] = sample;
                this.writePtr = (this.writePtr + 1) % this.bufferSize;
            });
            this.availableData += event.data.length;
            this.port.postMessage(
                {
                    availableData: this.availableData,
                    availableFreeSpace: freeSpace
                }
            );
        };
    }

    process(inputs, outputs, parameters) {
        const output = outputs[0];
        const channel = output[0];
        if (this.availableData < channel.length) {
            return true;
        }

        for (let i = 0; i < channel.length; i++) { // channel.length == 128
            channel[i] = this.buffer[this.readPtr];
            this.buffer[this.readPtr] = 0.0;
            this.readPtr = (this.readPtr + 1) % this.bufferSize;
        }

        this.availableData -= channel.length;

        return true;
    }
}

registerProcessor("audio", ExampleWorklet);