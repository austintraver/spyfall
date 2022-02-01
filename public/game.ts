// The duration for a game of Spyfall.
const duration = 5 * 60 * 1000;

// Redisplay the amount of time left in the game, doing so once per second,
// until the game ends.
function countdownTimer() {
    let timer: HTMLTimeElement | null = document.querySelector("time#timer");
    if (timer === null) {
        console.error("Could not find the timer element.");
        return;
    }
    const remaining = duration + new Date(timer.dateTime).getTime() - Date.now()
    if (remaining > 0) {
        let minutes = Math.floor((remaining / 1000 / 60) % 60).toString().padStart(2, '0')
        let seconds = Math.floor((remaining / 1000) % 60).toString().padStart(2, '0')
        timer.innerText = `${minutes}:${seconds}`
    }
}

document.addEventListener('DOMContentLoaded', () => {
    let timer: HTMLTimeElement | null = document.querySelector('time#timer');
    if (timer === null) {
        console.error("Could not find the timer element.");
        return;
    }
    timer.dateTime = new Date().toISOString()
    setInterval(countdownTimer, 1000)
});
