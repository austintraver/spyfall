function handleEvents() {
    let src = new EventSource('/events')
    console.log(
        src.withCredentials,
        src.readyState,
        src.url,
    )
    src.addEventListener('open', () => {
        console.log("Connection to server opened.")
    })
    src.addEventListener('message', (e) => {
        console.log(`message: ${e.data}`)
    })

    src.addEventListener('error', () => {
        console.log("EventSource failed.")
    })

    let button = document.querySelector('button')
    button?.addEventListener('click', () => {
        src.close()
        console.log('Connection closed')
    })
}

document.addEventListener('DOMContentLoaded', handleEvents)

document.addEventListener('DOMContentLoaded', () => {
    let hostButton = document.querySelector('button#create')

    if (hostButton === null) {
        console.error('Unable to locate the host button.')
        return
    }

    let joinButton = document.querySelector('button#join')
    if (joinButton === null) {
        console.error('Unable to locate the join button.')
        return
    }

    hostButton.addEventListener('click', () => {

        // The lowest possible lobby ID number.
        const lo = 1001;
        // The highest possible lobby ID number.
        const hi = 9999;
        let min = Math.ceil(lo)
        let max = Math.floor(hi)

        // Generate a random lobby ID between min and max (inclusive)
        let lobbyID = Math.floor(Math.random() * (max - min + 1) + min);

        let name = document.querySelector('input#name').value
        let room = document.querySelector('input#room').value
        let url = `/lobby?name=${name}&room=${room}`
        window.location.href = url
    })

    joinButton.addEventListener('click', () => {
        let name = document.querySelector('input#name').value
        let room = document.querySelector('input#room').value
        let url = `/join?name=${name}&room=${room}`
        window.location.href = url
    })

})
