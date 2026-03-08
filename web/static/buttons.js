const playPauseButton = document.getElementsByClassName("playPause")[0];
playPauseButton.addEventListener("click", function() {
    play();
});

const restartButton = document.getElementsByClassName("reset")[0];
restartButton.addEventListener("click", function() {
    restart();
});

const play = async () => {
    const payload = {
        comand: "playPause"
    };
    await call(payload);
};

const restart = async () => {
    const payload = {
        comand: "reset"
    };
    await call(payload);
};

const call = async (payload) => {
    await fetch("/player", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    });
}