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
